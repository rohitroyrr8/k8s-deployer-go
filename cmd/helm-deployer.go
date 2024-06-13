package main

import (
	"fmt"
	"os"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/getter"
)

/**
* /v1/create-new-release [deployment, ]
*
**/

func mains() {
	// Set up the Helm environment
	settings := cli.New()

	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), settings.Namespace(), os.Getenv("HELM_DRIVER"), func(format string, v ...interface{}) {
		fmt.Sprintf(format, v)
	}); err != nil {
		fmt.Printf("Error initializing action configuration: %v", err)
		return
	}

	// Define the chart you want to install
	chartPath := "./charts/chaincode"

	// Load the chart
	chart, err := loader.Load(chartPath)
	if err != nil {
		fmt.Printf("Error loading chart: %v", err)
		return
	}

	// Create an install action
	client := action.NewInstall(actionConfig)
	client.ReleaseName = "my-release"

	client.Namespace = settings.Namespace()

	// // Set custom values
	// valueOpts := &values.Options{
	//     ValueFiles: []string{"values.yaml"}, // Path to your values file
	// }
	// p := getter.All(settings)

	// vals, err := valueOpts.MergeValues(p)
	// if err != nil {
	//     fmt.Printf("Error merging values: %v", err)
	//     return
	// }

	// Set values if needed
	valueOpts := &values.Options{}
	p := getter.All(settings)

	vals, err := valueOpts.MergeValues(p)
	if err != nil {
		fmt.Printf("Error merging values: %v", err)
		return
	}

	// Install the chart
	rel, err := client.Run(chart, vals)
	if err != nil {
		fmt.Printf("Error installing chart: %v", err)
		return
	}

	fmt.Printf("Successfully installed chart: %s\n", rel.Name)
}
