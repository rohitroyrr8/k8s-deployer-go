package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/loyyal/k8s-deployer-go/controllers"
)

var (
	service            = "k8s-deployer-go"
	version            = "1.0"
	ctx                context.Context
	logger             *log.Logger
	server             *gin.Engine
	deployerController *controllers.DeployerController
)

func init() {
	ctx = context.TODO()
	deployerController.New(logger)

	server = gin.Default()
}

func main() {
	logger = log.New(os.Stderr, fmt.Sprintf(service+"[%s]: ", version), log.Llongfile|log.Lmicroseconds|log.LstdFlags)
	defer func() {
		logger.Println("exiting the services...")
	}()

	server.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	basepath := server.Group("/v1")
	deployerController.Routes(basepath)

	server.Run()

}
