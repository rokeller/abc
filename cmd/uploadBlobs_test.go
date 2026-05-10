package cmd

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"syscall"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
)

func Test_runUploadBlobs(t *testing.T) {
	container := setupContainer("uploads")
	tempBase := os.TempDir()
	setupTempFile(path.Join(tempBase, "upload-blob-1.txt"), "0001")
	setupTempFile(path.Join(tempBase, "upload-blob-2.txt"), "0002")
	setupBlob("uploads", "upload-blob-2.txt")

	tc := []testCase{
		{
			name: "NoArgs",
			args: []string{"blobs", "up"},
			err:  errors.New("requires at least 1 arg(s), only received 0"),
		},
		{
			name: "FileDoesNotExist",
			args: []string{
				"blobs", "up", "-n=foo", "-c=uploads",
				path.Join(tempBase, "file-does-not-exist.txt"),
			},
			err: fmt.Errorf("failed to open file \"/tmp/file-does-not-exist.txt\": %w",
				&fs.PathError{Op: "open", Path: "/tmp/file-does-not-exist.txt", Err: syscall.ENOENT}),
		},
		{
			name: "WithoutFlagsUploadFlags",
			args: []string{
				"blobs", "up", "-n=foo", "-c=uploads",
				path.Join(tempBase, "upload-blob-1.txt"),
			},
			verify: func(t *testing.T) {
				content, _ := downloadBlobContent(t, container, "upload-blob-1.txt")
				if string(content) != "0001" {
					t.Errorf("got content %q, want %q", string(content), "0001")
				}
			},
		},
		{
			name: "WithNoOverwrite",
			args: []string{
				"blobs", "up", "-n=foo", "-c=uploads", "--overwrite=false",
				path.Join(tempBase, "upload-blob-2.txt"),
			},
			err: errors.New("failed to overwrite \"upload-blob-2.txt\" with file \"upload-blob-2.txt\": blob already exists"),
		},
		{
			name: "WithOverwriteAndContentType",
			args: []string{
				"blobs", "up", "-n=foo", "-c=uploads",
				"--overwrite", "--content-type=text/plain",
				path.Join(tempBase, "upload-blob-2.txt"),
			},
			verify: func(t *testing.T) {
				content, contentType := downloadBlobContent(t, container, "upload-blob-2.txt")
				if string(content) != "0002" {
					t.Errorf("got content %q, want %q", string(content), "0002")
				}
				if contentType != "text/plain" {
					t.Errorf("got content type %q, want %q", contentType, "text/plain")
				}
			},
		},
		{
			name: "WithPrefixAndNoOverwrite",
			args: []string{
				"blobs", "up", "-n=foo", "-c=uploads",
				"--overwrite=false", "--content-type=text/plain", "--prefix=test/",
				path.Join(tempBase, "upload-blob-2.txt"),
			},
			verify: func(t *testing.T) {
				content, contentType := downloadBlobContent(t, container, "test/upload-blob-2.txt")
				if string(content) != "0002" {
					t.Errorf("got content %q, want %q", string(content), "0002")
				}
				if contentType != "text/plain" {
					t.Errorf("got content type %q, want %q", contentType, "text/plain")
				}
			},
		},
		{
			name: "MultipleFiles",
			args: []string{
				"blobs", "up", "-n=foo", "-c=uploads", "--prefix=multi/",
				path.Join(tempBase, "upload-blob-1.txt"),
				path.Join(tempBase, "upload-blob-2.txt"),
			},
			verify: func(t *testing.T) {
				content, _ := downloadBlobContent(t, container, "multi/upload-blob-1.txt")
				if string(content) != "0001" {
					t.Errorf("got content %q, want %q", string(content), "0001")
				}
				content, _ = downloadBlobContent(t, container, "multi/upload-blob-2.txt")
				if string(content) != "0002" {
					t.Errorf("got content %q, want %q", string(content), "0002")
				}
			},
		},
	}

	executeTestCases(t, tc)
}

func downloadBlobContent(t *testing.T, container *container.Client, blobName string) ([]byte, string) {
	client := container.NewBlobClient(blobName)
	resp, err := client.DownloadStream(t.Context(), &blob.DownloadStreamOptions{})
	if nil != err {
		t.Errorf("got error: %v", err)
	}
	defer resp.Body.Close()
	content, err := io.ReadAll(resp.Body)
	if nil != err {
		t.Errorf("got error: %v", err)
	}
	return content, *resp.ContentType
}
