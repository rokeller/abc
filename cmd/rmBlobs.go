package cmd

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/bloberror"
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
	RunE:  runRmBlobs,
}

func runRmBlobs(cmd *cobra.Command, args []string) error {
	containerClient := execCtx.serviceClient.NewContainerClient(execCtx.containerName)

	for _, blobName := range args {
		glog.Infof("remove blob: %s", blobName)
		blobClient := containerClient.NewBlobClient(blobName)

		_, err := blobClient.Delete(context.Background(), &blob.DeleteOptions{})

		if nil != err {
			if bloberror.HasCode(err, bloberror.ContainerNotFound) {
				err = fmt.Errorf("container \"%s\" does not exist", execCtx.containerName)
			} else if bloberror.HasCode(err, bloberror.BlobNotFound) {
				err = fmt.Errorf("blob \"%s\" does not exist in container \"%s\"",
					blobName, execCtx.containerName)
			}

			return err
		}

		fmt.Fprintln(cmd.OutOrStdout(), blobName)
	}

	return nil
}
