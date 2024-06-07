package cmd

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/bloberror"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/service"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

func init() {
	containersCmd.AddCommand(rmContainersCmd)
}

var rmContainersCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove containers",
	Long:  "Remove containers from a storage account",
	Args:  cobra.MinimumNArgs(1),
	RunE:  runRmContainers,
}

func runRmContainers(cmd *cobra.Command, args []string) error {
	for _, containerName := range args {
		glog.Infof("remove container: %s", containerName)
		_, err := execCtx.serviceClient.DeleteContainer(
			context.Background(), containerName, &service.DeleteContainerOptions{})

		containerDeleted := nil == err
		if nil != err {
			if bloberror.HasCode(err, bloberror.ContainerNotFound) {
				// the container no longer existing is the goal of this command,
				// so when it doesn't exist to begin with, that's not an error.
				err = nil
			} else if bloberror.HasCode(err, bloberror.InvalidResourceName) {
				err = fmt.Errorf("the container name \"%s\" is not permitted", containerName)
			}
		}

		if nil != err {
			return err
		}

		if containerDeleted {
			fmt.Fprintln(cmd.OutOrStdout(), containerName)
		}
	}

	return nil
}
