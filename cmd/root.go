package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/service"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)
var version string

var rootCmd = &cobra.Command{
	Use:   "abc",
	Short: "abc manages blobs in Azure storage",
	Long:  "abc (Azure Blob Commands) helps managing blobs in Azure storage",
	Args:  cobra.NoArgs,

	Version: version,

	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// glog flags parsing
		flag.Parse()

		accountFlag := cmd.Flag("account")
		if nil != accountFlag {
			accountName := accountFlag.Value.String()
			initExecContext(accountName)
		}
	},
}

var execCtx executionContext

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(rootCmd.OutOrStdout(), err)
		os.Exit(1)
	}
}

func addAccountFlag(c *cobra.Command) {
	c.PersistentFlags().StringP("account", "n", "", "name of the storage account")
	c.MarkPersistentFlagRequired("account")
}

func init() {
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
}

func initExecContext(accountName string) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if nil != err {
		glog.Exitf("failed to get Azure credentials: %v", err)
	}

	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)
	serviceClient, err := service.NewClient(serviceURL, cred, nil)
	if nil != err {
		glog.Exitf("failed to get blob service client: %v", err)
	}

	glog.Infof("setting execution context: %s", serviceURL)
	execCtx = executionContext{
		serviceURL:    serviceURL,
		serviceClient: serviceClient,
	}
}
