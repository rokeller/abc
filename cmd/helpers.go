package cmd

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/spf13/cobra"
)

func getFlagValue(c *cobra.Command, flagName string) *string {
	flag := c.Flag(flagName)
	if nil != flag {
		return to.Ptr(flag.Value.String())
	}

	return nil
}
