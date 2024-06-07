package cmd

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/service"
	"github.com/spf13/cobra"
)

type testCase struct {
	args   []string
	stdOut string
	stdErr string
	err    error
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

const AzuriteConnectionString = "DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;BlobEndpoint=http://127.0.0.1:10000/devstoreaccount1;"

func setup() {
	client, err := service.NewClientFromConnectionString(AzuriteConnectionString, nil)

	if nil != err {
		fmt.Printf("failed to create storage client: %v", err)
		os.Exit(1)
	}

	execCtx.serviceClient = client
}

func execute(t *testing.T, args ...string) (string, string, error) {
	t.Helper()

	c := rootCmd
	bufOut := new(bytes.Buffer)
	bufErr := new(bytes.Buffer)

	c.SetOut(bufOut)
	c.SetErr(bufErr)
	c.SetArgs(args)

	c.PersistentPreRun = func(cmd *cobra.Command, args []string) {}

	err := c.Execute()

	return strings.TrimSpace(bufOut.String()),
		strings.TrimSpace(bufErr.String()),
		err
}
