package cmd

import (
	"errors"
	"strings"
	"testing"
)

func TestBlobsListCmd_Interface(t *testing.T) {
	tc := []testCase{
		{
			name: "unsupported argument",
			args: []string{"blobs", "ls", "foo"},
			err:  errors.New("unknown command \"foo\" for \"abc blobs ls\""),
		},
		{
			name: "missing required flags",
			args: []string{"blobs", "ls"},
			err:  errors.New("required flag(s) \"account\", \"container\" not set"),
		},
	}

	executeTestCases(t, tc)
}

func TestBlobsListCmd_Functionality(t *testing.T) {
	manyBlobs := sequentialStrings(1, 11, "blob%04d.txt")
	setupManyBlobs(manyBlobs)

	tc := []testCase{
		{
			name:   "list of blobs",
			args:   []string{"blobs", "ls", "-n=foo", "-c=blobs"},
			stdOut: strings.Join(manyBlobs, "\n"),
		},
		{
			name: "list of blobs: container does not exist",
			args: []string{"blobs", "ls", "-n=foo", "-c=blah"},
			err:  errors.New("container \"blah\" does not exist"),
		},
		{
			name:   "filtered list of blobs - no match",
			args:   []string{"blobs", "ls", "-n=foo", "-c=blobs", "-p=foo"},
			stdOut: "",
		},
		{
			name:   "filtered list of blobs - all match",
			args:   []string{"blobs", "ls", "-n=foo", "-c=blobs", "-p=blob"},
			stdOut: strings.Join(manyBlobs, "\n"),
		},
		{
			name:   "filtered list of blobs - some match",
			args:   []string{"blobs", "ls", "-n=foo", "-c=blobs", "-p=blob000"},
			stdOut: strings.Join(manyBlobs[0:9], "\n"),
		},
	}

	executeTestCases(t, tc)
}
