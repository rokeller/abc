package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

func init() {
	blobsCmd.AddCommand(listBlobsCmd)
}

var listBlobsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all blobs",
	Long:  "List all blobs in a storage account's blob container",
	Args:  cobra.NoArgs,
	Run:   runListBlobsFlat,
}

func runListBlobsFlat(cmd *cobra.Command, args []string) {
	containerClient := execCtx.serviceClient.NewContainerClient(execCtx.containerName)
	glog.Infof("list blobs flat: %s", containerClient.URL())

	pager := containerClient.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{
		Include: container.ListBlobsInclude{Snapshots: false, Versions: false},
	})

	for pager.More() {
		resp, err := pager.NextPage(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
		for _, blob := range resp.ListBlobsFlatSegmentResponse.Segment.BlobItems {
			accessTierInferred := "set"
			if nil != blob.Properties.AccessTierInferred && *blob.Properties.AccessTierInferred {
				accessTierInferred = "inferred"
			}

			glog.Infof("found blob: %s (%s, %s)",
				*blob.Name, *blob.Properties.AccessTier, accessTierInferred)
			fmt.Println(*blob.Name)
		}
	}
}
