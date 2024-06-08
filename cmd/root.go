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

	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// glog flags parsing
		flag.Parse()

		if client, err := clientFactory(cmd); nil != err {
			return err
		} else {
			execCtx.serviceClient = client
			return nil
		}
	},
}

var clientFactory func(c *cobra.Command) (*service.Client, error) = clientFactoryForRelease
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
	cobra.EnableTraverseRunHooks = true
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
}

func clientFactoryForRelease(c *cobra.Command) (*service.Client, error) {
	accountFlag := c.Flag("account")
	var accountName string

	if nil != accountFlag {
		accountName = accountFlag.Value.String()
	} else {
		// if the flag is not present, it's because it's not needed -- otherwise
		// cobra would have already enforced it / reported an error.
		// in fact, the flag will not be present for some automatic commands like
		// for instance 'completion' and 'help' commands.
		return nil, nil
	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if nil != err {
		return nil, fmt.Errorf("failed to get Azure credentials: %v", err)
	}

	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)
	serviceClient, err := service.NewClient(serviceURL, cred, nil)
	if nil != err {
		return nil, fmt.Errorf("failed to get blob service client: %v", err)
	}

	return serviceClient, nil
}
