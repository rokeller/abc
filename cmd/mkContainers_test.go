package cmd

import (
	"errors"
	"strings"
	"testing"
)

func TestContainersMakeCmd_Interface(t *testing.T) {
	tc := []testCase{
		{
			name: "missing arguments",
			args: []string{"containers", "mk"},
			err:  errors.New("requires at least 1 arg(s), only received 0"),
		},
		{
			name: "missing required flags",
			args: []string{"containers", "mk", "foo"},
			err:  errors.New("required flag(s) \"account\" not set"),
		},
	}

	executeTestCases(t, tc)
}

func TestContainersMakeCmd_Functionality(t *testing.T) {
	tooLong := strings.Repeat("a", 64)
	validContainerNames := []string{"abc", "def", "ghi"}
	tc := []testCase{
		{
			name: "invalid container name: too short",
			args: []string{"containers", "mk", "-n=foo", "a"},
			err:  errors.New("the container name \"a\" is not permitted"),
		},
		{
			name: "invalid container name: too long",
			args: []string{"containers", "mk", "-n=foo", tooLong},
			err:  errors.New("the container name \"" + tooLong + "\" is not permitted"),
		},
		{
			name:   "make containers",
			args:   append([]string{"containers", "mk", "-n=foo"}, validContainerNames...),
			stdOut: strings.Join(validContainerNames, "\n"),
		},
		{
			name:   "make containers: with dupes",
			args:   []string{"containers", "mk", "-n=foo", "xyz", "xyz"},
			stdOut: "xyz",
		},
	}

	executeTestCases(t, tc)
}
