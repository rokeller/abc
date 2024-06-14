package cmd

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/service"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

func init() {
	containersCmd.AddCommand(listContainersCmd)
	listContainersCmd.Flags().StringP(
		"prefix", "p", "", "prefix of containers to list")
}

var listContainersCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all containers",
	Long:  "List all containers in a storage account",
	Args:  cobra.NoArgs,
	RunE:  runListContainers,
}

func runListContainers(cmd *cobra.Command, args []string) error {
	prefix := getFlagValue(cmd, "prefix")
	glog.Infof("list containers: %s, prefix: %s", execCtx.serviceClient.URL(), *prefix)
	pager := execCtx.serviceClient.NewListContainersPager(&service.ListContainersOptions{
		Prefix: prefix,
	})

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
