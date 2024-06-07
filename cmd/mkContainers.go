package cmd

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/bloberror"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

func init() {
	containersCmd.AddCommand(mkContainersCmd)
}

var mkContainersCmd = &cobra.Command{
	Use:   "mk",
	Short: "Create one or multiple containers",
	Long:  "Create one or multiple containers in a storage account",
	Args:  cobra.MinimumNArgs(1),
	RunE:  runMakeContainers,
}

func runMakeContainers(cmd *cobra.Command, args []string) error {
	for _, containerName := range args {
		glog.Infof("make container: %s on %s", containerName, execCtx.serviceClient.URL())

		_, err := execCtx.serviceClient.CreateContainer(context.Background(), containerName, nil)
		// we won't treat containers that already exist as errors, but we won't
		// list them twice in output
		containerCreated := nil == err

		if nil != err {
			if bloberror.HasCode(err, bloberror.ContainerAlreadyExists) {
				err = nil
			} else if bloberror.HasCode(err, bloberror.OutOfRangeInput) {
				err = fmt.Errorf("the container name \"%s\" is not permitted", containerName)
			}
		}

		if nil != err {
			return err
		}

		if containerCreated {
			fmt.Fprintln(cmd.OutOrStdout(), containerName)
		}
	}

	return nil
}
