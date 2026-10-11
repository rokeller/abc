package cmd

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/bloberror"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

func init() {
	blobsCmd.AddCommand(uploadBlobsCmd)

	uploadBlobsCmd.Flags().StringP(
		"prefix", "p", "", "prefix of blobs to upload")
	uploadBlobsCmd.Flags().StringP(
		"content-type", "m", "",
		"content type of blobs to upload; only a single content type can be specified - if specified, applies to all uploaded blobs")
	uploadBlobsCmd.Flags().BoolP(
		"overwrite", "f", true, "overwrite the blob if it already exists")

	uploadBlobsCmd.Flags().StringToString(
		"metadata", map[string]string{},
		"metadata for blobs to upload, formatted as a comma-separated list of entries with 'key=val'")
	uploadBlobsCmd.Flags().StringToString(
		"tags", map[string]string{},
		"tags for blobs to upload, formatted as a comma-separated list of entries with 'key=val'")
}

var uploadBlobsCmd = &cobra.Command{
	Use:     "up",
	Aliases: []string{"upload"},
	Short:   "Upload blobs",
	Long:    "Uploads local files to block blobs in a storage account's blob container",
	Args:    cobra.MinimumNArgs(1),
	RunE:    runUploadBlobs,
}

func runUploadBlobs(cmd *cobra.Command, args []string) error {
	prefixFlag := getFlagValue(cmd, "prefix")
	contentTypeFlag := getFlagValue(cmd, "content-type")
	overwriteFlag := getBoolFlagValue(cmd, "overwrite")
	metadata, err := cmd.Flags().GetStringToString("metadata")
	if nil != err {
		return err
	}
	tags, err := cmd.Flags().GetStringToString("tags")
	if nil != err {
		return err
	}

	ctx := fileUploadContext{
		ctx:             cmd.Context(),
		blobPrefix:      "",
		overwrite:       true,
		contentType:     "",
		containerClient: execCtx.serviceClient.NewContainerClient(execCtx.containerName),
	}
	if nil != prefixFlag {
		ctx.blobPrefix = *prefixFlag
	}
	if nil != overwriteFlag {
		ctx.overwrite = *overwriteFlag
	}
	if nil != contentTypeFlag {
		ctx.contentType = *contentTypeFlag
	}
	if nil != metadata {
		ctx.metadata = metadata
	}
	if nil != tags {
		ctx.tags = tags
	}

	for _, arg := range args {
		if err := ctx.uploadFile(arg); nil != err {
			return err
		}
	}
	return nil
}

type fileUploadContext struct {
	ctx             context.Context
	blobPrefix      string
	overwrite       bool
	contentType     string
	containerClient *container.Client

	metadata map[string]string
	tags     map[string]string
}

func (c *fileUploadContext) uploadFile(filePath string) error {
	fileName := path.Base(filePath)
	blobName := fmt.Sprintf("%s%s", c.blobPrefix, fileName)
	glog.Infof("upload file %q to %q", fileName, blobName)
	blobClient := c.containerClient.NewBlockBlobClient(blobName)

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %q: %w", filePath, err)
	}
	defer file.Close()

	accessConditions := &blob.AccessConditions{
		ModifiedAccessConditions: &blob.ModifiedAccessConditions{},
	}
	if !c.overwrite {
		accessConditions.ModifiedAccessConditions.IfNoneMatch = to.Ptr(azcore.ETagAny)
	}

	httpHeaders := &blob.HTTPHeaders{}
	if "" != c.contentType {
		httpHeaders.BlobContentType = &c.contentType
	}

	_, err = blobClient.UploadFile(c.ctx, file, &blockblob.UploadFileOptions{
		AccessConditions: accessConditions,
		HTTPHeaders:      httpHeaders,
		Metadata:         c.convertMetadata(),
		Tags:             c.tags,
	})
	if nil != err && bloberror.HasCode(err, bloberror.BlobAlreadyExists) {
		return fmt.Errorf("failed to overwrite %q with file %q: blob already exists",
			blobName, fileName)
	}
	return err
}

func (c *fileUploadContext) convertMetadata() map[string]*string {
	if nil == c.metadata {
		return nil
	}
	res := make(map[string]*string, len(c.metadata))
	for k, v := range c.metadata {
		res[k] = &v
	}
	return res
}
