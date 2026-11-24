package cmd

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestBlobsListCmd_Interface(t *testing.T) {
	tc := []testCase{
		{
			name: "missing required flags",
			args: []string{"blobs", "ls"},
			err:  errors.New("required flag(s) \"account\", \"container\" not set"),
		},
		{
			name: "unsupported argument",
			args: []string{"blobs", "ls", "foo"},
			err:  errors.New("unknown command \"foo\" for \"abc blobs ls\""),
		},
	}

	executeTestCases(t, tc)
}

func TestBlobsListCmd_Functionality(t *testing.T) {
	manyBlobs := sequentialStrings(1, 11, "normal/blob%04d.txt")
	setupManyBlobs(manyBlobs)

	tc := []testCase{
		{
			name:   "list of blobs",
			args:   []string{"blobs", "ls", "-n=foo", "-c=blobs", "-p=norm"},
			stdOut: strings.Join(manyBlobs, "\n"),
		},

		{
			name: "list of blobs: container does not exist",
			args: []string{"blobs", "ls", "-n=foo", "-c=blah"},
			err:  errors.New("container \"blah\" does not exist"),
		},

		// prefix filtering
		{
			name:   "filtered list of blobs - no match",
			args:   []string{"blobs", "ls", "-n=foo", "-c=blobs", "-p=foo"},
			stdOut: "",
		},
		{
			name:   "filtered list of blobs - all match",
			args:   []string{"blobs", "ls", "-n=foo", "-c=blobs", "-p=normal/blob"},
			stdOut: strings.Join(manyBlobs, "\n"),
		},
		{
			name:   "filtered list of blobs - some match",
			args:   []string{"blobs", "ls", "-n=foo", "-c=blobs", "-p=normal/blob000"},
			stdOut: strings.Join(manyBlobs[0:9], "\n"),
		},

		// // deleted blob
		// {
		// 	name:   "deleted blob",
		// 	args:   []string{"blobs", "ls", "-n=foo", "-c=blobs", "-p=del", "-d"},
		// 	stdOut: "deleted/deleted-blob.txt@deleted",
		// },

		// // blob snapshots
		// {
		// 	name:   "blob snapshots",
		// 	args:   []string{"blobs", "ls", "-n=foo", "-c=blobs", "-p=snap", "-s"},
		// 	stdOut: fmt.Sprintf("snapshot/blob-with-snapshot.txt@%s\nsnapshot/blob-with-snapshot.txt", *snapshot),
		// },

		// // blob versions
		// {
		// 	name:   "blob version",
		// 	args:   []string{"blobs", "ls", "-n=foo", "-c=blobs", "-p=ver", "--versions"},
		// 	stdOut: fmt.Sprintf("version/versioned-blob.txt@v%s\nversion/versioned-blob.txt", "blah"),
		// },
	}

	executeTestCases(t, tc)
}

func TestBlobsListCmd_Functionality_BlobSnapshots(t *testing.T) {
	snapshot := setupSnapshot("snapshot/blob-with-snapshot.txt")

	tc := []testCase{
		{
			name:   "blob snapshots",
			args:   []string{"blobs", "ls", "-n=foo", "-c=blobs", "-p=snap", "-s"},
			stdOut: fmt.Sprintf("snapshot/blob-with-snapshot.txt@%s\nsnapshot/blob-with-snapshot.txt", *snapshot),
		},
	}

	executeTestCases(t, tc)
}
