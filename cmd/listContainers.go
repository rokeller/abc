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
	Run:   runListContainers,
}

func runListContainers(cmd *cobra.Command, args []string) {
	glog.Infof("list containers: %s", execCtx.serviceClient.URL())

	pager := execCtx.serviceClient.NewListContainersPager(nil)
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		if nil != err {
			glog.Fatalf("failed to fetch page %v", err)
			continue
		}

		for _, container := range resp.ContainerItems {
			fmt.Println(*container.Name)
		}
	}
}
