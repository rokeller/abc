package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	addAccountFlag(containersCmd)
	rootCmd.AddCommand(containersCmd)
}

var containersCmd = &cobra.Command{
	Use:   "containers",
	Short: "Operate on containers",
	Long:  "Operate on containers in a storage account",
	Args:  cobra.NoArgs,
}
