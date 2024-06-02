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

var rootCmd = &cobra.Command{
	Use:   "abc",
	Short: "abc manages blobs in Azure storage",
	Long:  "abc (Azure Blob Commands) helps managing blobs in Azure storage",
	Args:  cobra.NoArgs,

	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// glog flags parsing
		flag.Parse()
	},
}

var execCtx executionContext

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	var (
		accountName string
	)

	cobra.OnInitialize(func() {
		initExecContext(accountName)
	})

	rootCmd.PersistentFlags().StringVarP(
		&accountName, "account", "n", "", "name of the storage account")
	rootCmd.MarkPersistentFlagRequired("account")

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
