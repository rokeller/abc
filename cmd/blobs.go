package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	addAccountFlag(blobsCmd)
	blobsCmd.PersistentFlags().StringP(
		"container", "c", "", "name of the blob container")
	blobsCmd.MarkPersistentFlagRequired("container")

	rootCmd.AddCommand(blobsCmd)
}

var blobsCmd = &cobra.Command{
	Use:   "blobs",
	Short: "Operate on blobs",
	Long:  "Operate on blobs in a storage account's blob container",
	Args:  cobra.NoArgs,

	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		rootCmd.PersistentPreRun(cmd, args)

		execCtx.containerName = cmd.Flag("container").Value.String()
	},
}
