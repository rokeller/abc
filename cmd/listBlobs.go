package cmd

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/bloberror"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

func init() {
	blobsCmd.AddCommand(listBlobsCmd)

	listBlobsCmd.Flags().StringP(
		"prefix", "p", "", "prefix of blobs to list")
	listBlobsCmd.Flags().BoolP(
		"deleted", "d", false, "include deleted blobs; deleted blobs will be listed with an '@deleted' suffix")
	listBlobsCmd.Flags().BoolP(
		"snapshots", "s", false, "include blob snapshots; snapshots will be listed with an '@<snapshot-id>' suffix")
}

var listBlobsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all blobs",
	Long:  "List all blobs in a storage account's blob container",
	Args:  cobra.NoArgs,
	RunE:  runListBlobsFlat,
}

func runListBlobsFlat(cmd *cobra.Command, args []string) error {
	containerClient := execCtx.serviceClient.NewContainerClient(execCtx.containerName)

	prefix := getFlagValue(cmd, "prefix")
	glog.Infof("list blobs flat: %s, prefix: %s", containerClient.URL(), prefix)

	includeDeleted := getBoolFlagValue(cmd, "deleted")
	includeSnapshots := getBoolFlagValue(cmd, "snapshots")

	pager := containerClient.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{
		Prefix: prefix,
		Include: container.ListBlobsInclude{
			Deleted:   *includeDeleted,
			Snapshots: *includeSnapshots,
		},
	})

	for pager.More() {
		resp, err := pager.NextPage(context.TODO())
		if err != nil {
			if bloberror.HasCode(err, bloberror.ContainerNotFound) {
				err = fmt.Errorf("container \"%s\" does not exist", execCtx.containerName)
			}

			return err
		}

		for _, blob := range resp.ListBlobsFlatSegmentResponse.Segment.BlobItems {
			accessTierInferred := "set"
			if nil != blob.Properties.AccessTierInferred && *blob.Properties.AccessTierInferred {
				accessTierInferred = "inferred"
			}

			glog.Infof("found blob: %s (%s, %s)",
				*blob.Name, *blob.Properties.AccessTier, accessTierInferred)

			suffix := ""
			if nil != blob.Deleted {
				suffix += "@deleted"
			}
			if nil != blob.Snapshot {
				suffix += fmt.Sprintf("@%s", *blob.Snapshot)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "%s%s\n", *blob.Name, suffix)
		}
	}

	return nil
}
