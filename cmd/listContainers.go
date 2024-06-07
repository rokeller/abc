package cmd

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

func init() {
	containersCmd.AddCommand(listContainersCmd)
}

var listContainersCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all containers",
	Long:  "List all containers in a storage account",
	Args:  cobra.NoArgs,
	RunE:  runListContainers,
}

func runListContainers(cmd *cobra.Command, args []string) error {
	glog.Infof("list containers: %s", execCtx.serviceClient.URL())
	pager := execCtx.serviceClient.NewListContainersPager(nil)
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		if nil != err {
			return err
		}

		for _, container := range resp.ContainerItems {
			fmt.Fprintln(cmd.OutOrStdout(), *container.Name)
		}
	}

	return nil
}
