package cmd

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func TestContainersListCmd_Interface(t *testing.T) {
	tc := []testCase{
		{
			name: "unsupported argument",
			args: []string{"containers", "ls", "foo"},
			err:  errors.New("unknown command \"foo\" for \"abc containers ls\""),
		},
		{
			name: "missing required flags",
			args: []string{"containers", "ls"},
			err:  errors.New("required flag(s) \"account\" not set"),
		},
	}

	executeTestCases(t, tc)
}

func TestContainersListCmd_Functionality(t *testing.T) {
	manyContainers := sequentialStrings(1, 11, "test%04d")
	setupManyContainers(manyContainers)

	tc := []testCase{
		{
			name:   "list of containers",
			args:   []string{"containers", "ls", "-n=foo"},
			stdOut: strings.Join(manyContainers, "\n"),
		},
		{
			name:   "filtered list of containers - no match",
			args:   []string{"containers", "ls", "-n=foo", "-p=foo"},
			stdOut: "",
		},
		{
			name:   "filtered list of containers - all match",
			args:   []string{"containers", "ls", "-n=foo", "-p=test"},
			stdOut: strings.Join(manyContainers, "\n"),
		},
		{
			name:   "filtered list of containers - some match",
			args:   []string{"containers", "ls", "-n=foo", "-p=test000"},
			stdOut: strings.Join(manyContainers[0:9], "\n"),
		},
	}

	executeTestCases(t, tc)
}

func setupManyContainers(containerNames []string) {
	pager := execCtx.serviceClient.NewListContainersPager(nil)
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		if nil != err {
			panic(err)
		}

		for _, container := range resp.ContainerItems {
			execCtx.serviceClient.DeleteContainer(context.Background(), *container.Name, nil)
		}
	}

	for _, containerName := range containerNames {
		execCtx.serviceClient.CreateContainer(context.Background(), containerName, nil)
	}
}
