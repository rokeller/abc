package cmd

import "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/service"

type executionContext struct {
	serviceURL    string
	serviceClient *service.Client

	containerName string
}
