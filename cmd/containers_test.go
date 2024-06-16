package cmd

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
)

func setupContainer(name string) *container.Client {
	containerClient := execCtx.serviceClient.NewContainerClient(name)
	containerClient.Create(context.Background(), nil)

	return containerClient
}
