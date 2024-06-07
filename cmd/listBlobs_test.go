package cmd

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestBlobsListCmd(t *testing.T) {
	is := is.New(t)

	manyBlobs := sequentialStrings(0, 10, "my/blob%04d.txt")
	setupManyBlobs(manyBlobs)

	tt := []testCase{
		{
			// extra argument 'foo' which is not supported
			args: []string{"blobs", "ls", "foo"},
			err:  errors.New("unknown command \"foo\" for \"abc blobs ls\""),
		},
		{
			// missing required flags
			args: []string{"blobs", "ls"},
			err:  errors.New("required flag(s) \"account\", \"container\" not set"),
		},
		{
			// container does not exist
			args: []string{"blobs", "ls", "-n=foo", "-c=does-not-exist"},
			err:  errors.New("container \"does-not-exist\" does not exist"),
		},
		{
			// include some blobs in results
			args:   []string{"blobs", "ls", "-n=foo", "-c=many-pages"},
			stdOut: strings.Join(manyBlobs, "\n"),
		},
	}

	for _, testCase := range tt {
		stdOut, stdErr, err := execute(t, testCase.args...)

		is.Equal(testCase.err, err)

		if testCase.err == nil {
			is.Equal(testCase.stdOut, stdOut)
			is.Equal(testCase.stdErr, stdErr)
		}
	}
}

func setupManyBlobs(blobNames []string) {
	containerClient := execCtx.serviceClient.NewContainerClient("many-pages")
	containerClient.Create(context.Background(), nil)

	for _, blobName := range blobNames {
		blobClient := containerClient.NewBlockBlobClient(blobName)
		blobClient.UploadBuffer(context.Background(), []byte{}, nil)
	}
}
