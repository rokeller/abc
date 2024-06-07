package cmd

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/service"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
	name   string
	args   []string
	err    error
	stdOut string
	stdErr string
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

func executeTestCases(t *testing.T, testCases []testCase) {
	t.Helper()

	for _, testCase := range testCases {
		t.Run(testCase.name, testCase.executeTestCase)
	}
}

func (testCase testCase) executeTestCase(t *testing.T) {
	t.Helper()

	stdOut, stdErr, err := execute(t, testCase.args...)

	assert.Equal(t, testCase.err, err, "expected error must match")

	if testCase.err == nil {
		assert.Equal(t, testCase.stdOut, stdOut, "stdout must match")
		assert.Equal(t, testCase.stdErr, stdErr, "stderr must match")
	}
}

func execute(t *testing.T, args ...string) (string, string, error) {
	t.Helper()

	c := rootCmd
	bufOut := new(bytes.Buffer)
	bufErr := new(bytes.Buffer)

	c.SetOut(bufOut)
	c.SetErr(bufErr)

	// reset all flags to start clean
	c.Flags().VisitAll(func(f *pflag.Flag) {
		f.Value.Set(f.DefValue)
		f.Changed = false
	})
	c.SetArgs(args)

	c.PersistentPreRun = func(cmd *cobra.Command, args []string) {}

	err := c.Execute()

	return strings.TrimSpace(bufOut.String()),
		strings.TrimSpace(bufErr.String()),
		err
}

func sequentialStrings(start, end int, template string) []string {
	items := make([]string, end-start)
	for i := range end - start {
		items[i] = fmt.Sprintf(template, i+start)
	}

	return items
}
