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

func TestGetFlagValue(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().StringP("test-flag", "f", "default-value", "")

	assert.Nil(t, getFlagValue(cmd, "does-not-exist"), "non-existing flags' values must be nil")
	assert.NotNil(t, getFlagValue(cmd, "test-flag"), "existing flags' values must not be nil")
	assert.Equal(t, "default-value", *getFlagValue(cmd, "test-flag"), "unset existing flags' values must equal the default")

	cmd.Flags().Set("test-flag", "test-value")
	assert.NotNil(t, getFlagValue(cmd, "test-flag"), "existing flags' values must not be nil")
	assert.Equal(t, "test-value", *getFlagValue(cmd, "test-flag"), "existing flags' values must be correct")
}

func TestGetBoolFlagValue(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().BoolP("false-flag", "f", false, "")
	cmd.Flags().BoolP("true-flag", "t", true, "")

	assert.Nil(t, getBoolFlagValue(cmd, "does-not-exist"), "non-existing flags' values must be nil")
	assert.NotNil(t, getBoolFlagValue(cmd, "false-flag"), "existing flags' values must not be nil")
	assert.False(t, *getBoolFlagValue(cmd, "false-flag"), "unset existing flags' values must equal the default")
	assert.NotNil(t, getBoolFlagValue(cmd, "true-flag"), "existing flags' values must not be nil")
	assert.True(t, *getBoolFlagValue(cmd, "true-flag"), "unset existing flags' values must equal the default")

	cmd.Flags().Set("false-flag", "true")
	cmd.Flags().Set("true-flag", "false")
	assert.NotNil(t, getBoolFlagValue(cmd, "false-flag"), "existing flags' values must not be nil")
	assert.True(t, *getBoolFlagValue(cmd, "false-flag"), "existing flags' values must be correct")
	assert.NotNil(t, getBoolFlagValue(cmd, "true-flag"), "existing flags' values must not be nil")
	assert.False(t, *getBoolFlagValue(cmd, "true-flag"), "existing flags' values must be correct")
}

func setup() {
	clientFactory = clientFactoryForTest
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

	// reset all flags to start clean
	cmdFn := func(c *cobra.Command) {
		c.Flags().VisitAll(func(f *pflag.Flag) {
			f.Value.Set(f.DefValue)
			f.Changed = false
		})
	}
	visitCommands(c.Commands(), cmdFn)

	bufOut, bufErr := captureStdOutAndErr(c)
	c.SetArgs(args)
	err := c.Execute()

	return strings.TrimSpace(bufOut.String()),
		strings.TrimSpace(bufErr.String()),
		err
}

func captureStdOutAndErr(c *cobra.Command) (bufOut, bufErr *bytes.Buffer) {
	bufOut = new(bytes.Buffer)
	bufErr = new(bytes.Buffer)

	c.SetOut(bufOut)
	c.SetErr(bufErr)

	return
}

func sequentialStrings(start, end int, template string) []string {
	items := make([]string, end-start)
	for i := range end - start {
		items[i] = fmt.Sprintf(template, i+start)
	}

	return items
}

func visitCommands(cs []*cobra.Command, cmdFn func(c *cobra.Command)) {
	for _, c := range cs {
		cmdFn(c)
		visitCommands(c.Commands(), cmdFn)
	}
}

const AzuriteConnectionString = "DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;BlobEndpoint=http://127.0.0.1:10000/devstoreaccount1;"

func clientFactoryForTest(c *cobra.Command) (*service.Client, error) {
	client, err := service.NewClientFromConnectionString(AzuriteConnectionString, nil)

	if nil != err {
		return nil, err
	}

	return client, nil
}
