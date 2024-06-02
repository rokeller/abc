package cmd

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

func init() {
	blobsCmd.AddCommand(rmBlobsCmd)
}

var rmBlobsCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove blobs",
	Long:  "Remove blobs in a storage account's blob container",
	Args:  cobra.MinimumNArgs(1),
	Run:   runRmBlobs,
}

func runRmBlobs(cmd *cobra.Command, args []string) {
	containerClient := execCtx.serviceClient.NewContainerClient(execCtx.containerName)

	for _, arg := range args {
		glog.Infof("remove blob: %s", arg)
		blobClient := containerClient.NewBlobClient(arg)

		blobClient.Delete(context.Background(), &blob.DeleteOptions{})
	}
}
