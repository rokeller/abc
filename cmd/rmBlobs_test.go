package cmd

import (
	"errors"
	"strings"
	"testing"
)

func TestBlobsRmCmd_Interface(t *testing.T) {
	tc := []testCase{
		{
			name: "missing required arguments",
			args: []string{"blobs", "rm"},
			err:  errors.New("requires at least 1 arg(s), only received 0"),
		},
		{
			name: "missing required flags",
			args: []string{"blobs", "rm", "test"},
			err:  errors.New("required flag(s) \"account\", \"container\" not set"),
		},
	}

	executeTestCases(t, tc)
}

func TestBlobsRmCmd_Functionality(t *testing.T) {
	manyBlobs := sequentialStrings(0, 10, "blob%04d.txt")
	setupManyBlobs(manyBlobs)

	tc := []testCase{
		{
			name: "container does not exist",
			args: []string{"blobs", "rm", "-n=foo", "-c=not-exist", "blob-1", "blob-2", "blob-3"},
			err:  errors.New("container \"not-exist\" does not exist"),
		},
		{
			name: "blob does not exist",
			args: []string{"blobs", "rm", "-n=foo", "-c=blobs", "blob-1", "blob-2", "blob-3"},
			err:  errors.New("blob \"blob-1\" does not exist in container \"blobs\""),
		},
		{
			name:   "blobs removed",
			args:   []string{"blobs", "rm", "-n=foo", "-c=blobs", "blob0000.txt", "blob0001.txt"},
			stdOut: strings.Join(manyBlobs[0:2], "\n"),
		},
	}

	executeTestCases(t, tc)
}
