package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/loyyal/k8s-deployer-go/common"
	"github.com/loyyal/k8s-deployer-go/middleware"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/getter"
)

type DeployerController struct {
	logger *log.Logger
}

type NewReleaseRequest struct {
	ReleaseName string `json:"release_name" binding:"required"`
	Namespace   string `json:"namespace" binding:"required"`
	ChannelName string `json:"channel_name"`
}

func (controller *DeployerController) New(logger *log.Logger) *DeployerController {
	return &DeployerController{
		logger: logger,
	}
}

func (controller *DeployerController) CreateNewRelease(ctx *gin.Context) {
	fName := "controllers/deployerController/login"

	var input NewReleaseRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		common.PrepareCustomError(ctx, http.StatusBadRequest, fName, "error: invalid request body provided", fmt.Sprintf("got :%s ", input))

		return
	}

	response, err := installChart()
	if err != nil {
		common.PrepareCustomError(ctx, http.StatusBadRequest, fName, "error: invalid request body provided", fmt.Sprintf("got :%s ", input))
		return
	}

	common.PrepareCustomResponse(ctx, "created new release", struct {
		Body string `json:"body"`
	}{Body: response})
}

func (controller *DeployerController) DeleteNewRelease(ctx *gin.Context) {
	fName := "controllers/deployerController/login"

	var input NewReleaseRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		common.PrepareCustomError(ctx, http.StatusBadRequest, fName, "error: invalid request body provided", fmt.Sprintf("got :%s ", input))

		return
	}

	response, err := deleteChart()
	if err != nil {
		common.PrepareCustomError(ctx, http.StatusBadRequest, fName, "error: invalid request body provided", fmt.Sprintf("got :%s ", input))
		return
	}

	common.PrepareCustomResponse(ctx, "created new release", struct {
		Body string `json:"body"`
	}{Body: response})
}

func installChart() (string, error) {

	// Set up the Helm environment
	settings := cli.New()

	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), settings.Namespace(), os.Getenv("HELM_DRIVER"), func(format string, v ...interface{}) {
		fmt.Sprintf(format, v)
	}); err != nil {
		fmt.Printf("Error initializing action configuration: %v", err)
		return "", err
	}

	// Define the chart you want to install
	chartPath := "./charts/chaincode"

	// Load the chart
	chart, err := loader.Load(chartPath)
	if err != nil {
		fmt.Printf("Error loading chart: %v", err)
		return "", err
	}

	// Create an install action
	client := action.NewInstall(actionConfig)
	client.ReleaseName = "my-release"

	client.Namespace = settings.Namespace()

	// Set values if needed
	valueOpts := &values.Options{}
	p := getter.All(settings)

	vals, err := valueOpts.MergeValues(p)
	if err != nil {
		fmt.Printf("Error merging values: %v", err)
		return "", err
	}

	// Install the chart
	rel, err := client.Run(chart, vals)
	if err != nil {
		fmt.Printf("Error installing chart: %v", err)
		return "", err
	}

	return rel.Name, nil

}

func deleteChart() (string, error) {

	// Set up the Helm environment
	settings := cli.New()

	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), settings.Namespace(), os.Getenv("HELM_DRIVER"), func(format string, v ...interface{}) {
		fmt.Sprintf(format, v)
	}); err != nil {
		fmt.Printf("Error initializing action configuration: %v", err)
		return "", err
	}

	// Create an uninstall action
	client := action.NewUninstall(actionConfig)

	// Specify the release name to be deleted
	releaseName := "my-release"

	// Uninstall the release
	rel, err := client.Run(releaseName)
	if err != nil {
		fmt.Printf("Error uninstalling release: %v", err)
		return "", err
	}

	return rel.Release.Name, nil

}

func (controller *DeployerController) Routes(group *gin.RouterGroup) {
	authRoute := group.Group("/deployer")

	authRoute.Use(middleware.BasicAuthMiddleware())
	authRoute.POST("/create-new-release", controller.CreateNewRelease)
	authRoute.POST("/delete-release", controller.DeleteNewRelease)
}
