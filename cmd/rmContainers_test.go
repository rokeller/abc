package cmd

import (
	"errors"
	"testing"
)

func TestContainersRmCmd_Interface(t *testing.T) {
	tc := []testCase{
		{
			name: "missing arguments",
			args: []string{"containers", "rm"},
			err:  errors.New("requires at least 1 arg(s), only received 0"),
		},
		{
			name: "missing required flags",
			args: []string{"containers", "rm", "foo"},
			err:  errors.New("required flag(s) \"account\" not set"),
		},
	}

	executeTestCases(t, tc)
}

func TestContainersRmCmd_Functionality(t *testing.T) {
	manyContainers := sequentialStrings(0, 10, "test%04d")
	setupManyContainers(manyContainers)

	tc := []testCase{
		{
			name:   "remove containers",
			args:   []string{"containers", "rm", "-n=foo", "test0001", "test0004"},
			stdOut: "test0001\ntest0004",
		},
		{
			name:   "remove containers: with dupes",
			args:   []string{"containers", "rm", "-n=foo", "test0000", "test0000"},
			stdOut: "test0000",
		},
		{
			name: "remove containers: with dupes",
			args: []string{"containers", "rm", "-n=foo", "abc-"},
			err:  errors.New("the container name \"abc-\" is not permitted"),
		},
	}

	executeTestCases(t, tc)
}
