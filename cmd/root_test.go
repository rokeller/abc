package cmd

import (
	"testing"

	"github.com/rokeller/abc/test_helpers"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestExecute_BareCommand(t *testing.T) {
	test_helpers.ExecuteWithExit(t, "TestExecute_BareCommand", func(t *testing.T) {
		bufOut, bufErr := captureStdOutAndErr(rootCmd)
		Execute()

		assert.Equal(t, "", bufErr.String(), "stderr matches")
		assert.Regexpf(
			t,
			"^abc \\(Azure Blob Commands\\) helps managing blobs in Azure storage\n",
			bufOut.String(),
			"root command output first line matches")
	}, 0)
}

func TestExecute_UnsupportedFlag(t *testing.T) {
	test_helpers.ExecuteWithExit(t, "TestExecute_UnsupportedFlag", func(t *testing.T) {
		rootCmd.SetArgs([]string{"--foo"})
		Execute()
	}, 1)
}

func TestClientFactory_WithoutAccount(t *testing.T) {
	c := &cobra.Command{}
	client, err := clientFactoryForRelease(c)

	assert.Nil(t, client, "the client must be nil")
	assert.Nil(t, err, "the error must be nil")
}

func TestClientFactory_WithoutCredentials(t *testing.T) {
	c := &cobra.Command{}
	addAccountFlag(c)
	c.SetArgs([]string{"--account=foo"})
	client, err := clientFactoryForRelease(c)

	assert.NotNil(t, client, "the client must not be nil")
	assert.Nil(t, err, "the error must be nil")
}
