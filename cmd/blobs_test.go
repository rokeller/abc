package cmd

import (
	"context"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
)

func setupManyBlobs(blobNames []string) {
	containerClient := setupContainer("blobs")

	for _, blobName := range blobNames {
		blobClient := containerClient.NewBlockBlobClient(blobName)
		blobClient.UploadBuffer(context.Background(), []byte{}, nil)
	}
}

func setupSnapshot(blobName string) *string {
	containerClient := setupContainer("blobs")

	blobClient := containerClient.NewBlockBlobClient(blobName)
	blobClient.UploadBuffer(context.Background(), []byte{}, nil)
	resp, _ := blobClient.CreateSnapshot(context.Background(), &blob.CreateSnapshotOptions{})

	return resp.Snapshot
}

func setupBlob(containerName, blobName string) {
	containerClient := setupContainer(containerName)
	blobClient := containerClient.NewBlockBlobClient(blobName)
	blobClient.UploadBuffer(context.Background(), []byte{}, nil)
}

func setupTempFile(filePath, content string) {
	os.WriteFile(filePath, []byte(content), 0644)
}
