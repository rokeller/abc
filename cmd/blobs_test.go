package cmd

import "context"

func setupManyBlobs(blobNames []string) {
	containerClient := execCtx.serviceClient.NewContainerClient("blobs")
	containerClient.Create(context.Background(), nil)

	for _, blobName := range blobNames {
		blobClient := containerClient.NewBlockBlobClient(blobName)
		blobClient.UploadBuffer(context.Background(), []byte{}, nil)
	}
}
