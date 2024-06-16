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

func getBoolFlagValue(c *cobra.Command, flagName string) *bool {
	val, err := c.Flags().GetBool(flagName)
	if nil != err {
		return nil
	}

	return &val
}
