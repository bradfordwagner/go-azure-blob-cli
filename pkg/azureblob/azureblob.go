package azureblob

import (
	"bytes"
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/sirupsen/logrus"
)

// AzureBlob
type AzureBlob struct {
	l             *logrus.Entry
	connected     bool
	serviceClient azblob.ServiceClient
	account, key  string
}

// NewAzureBlob - creates a new instance of azure blob storage connector
func NewAzureBlob(account, key string) *AzureBlob {
	return &AzureBlob{
		l:       logrus.WithField("component", "azure_blob"),
		account: account,
		key:     key,
	}
}

func (a *AzureBlob) authenticate() (err error) {
	// check if we have already connected
	if a.connected {
		return
	}

	creds, err := azblob.NewSharedKeyCredential(a.account, a.key)
	if err != nil {
		a.l.WithError(err).Error("could not connect to azure blob storage")
		return err
	}

	// create service client
	serviceEndpoint := fmt.Sprintf("https://%s.blob.core.windows.net/", a.account)
	a.serviceClient, err = azblob.NewServiceClientWithSharedKey(serviceEndpoint, creds, nil)
	if err != nil {
		return err
	}
	a.connected = true

	return
}
func (a *AzureBlob) CreateContainer(ctx context.Context, name string) (err error) {
	err = a.authenticate()
	if err != nil {
		return
	}
	_, err = a.serviceClient.CreateContainer(ctx, name, nil)
	return
}

func (a *AzureBlob) ListContainers(ctx context.Context) (containerNames []string, err error) {
	// login
	err = a.authenticate()
	if err != nil {
		return
	}

	// get all the containers
	containers := a.serviceClient.ListContainers(nil)
	err = containers.Err()
	if err != nil {
		return
	}

	// parse response
	for containers.NextPage(ctx) {
		c := containers.PageResponse()
		for _, i := range c.ContainerItems {
			containerNames = append(containerNames, *i.Name)
		}
	}

	return
}

func (a *AzureBlob) ListFiles(ctx context.Context, container, prefix string) (res map[string]bool, err error) {
	err = a.authenticate()
	if err != nil {
		return
	}

	// setup container client
	containerClient := a.serviceClient.NewContainerClient(container)

	// pull blobs from hierarchy based on prefix
	hierarchy := containerClient.ListBlobsHierarchy("", &azblob.ContainerListBlobHierarchySegmentOptions{
		Prefix: &prefix,
	})

	// parse into result map
	res = make(map[string]bool)
	for hierarchy.NextPage(ctx) {
		b := hierarchy.PageResponse()
		for _, item := range b.Segment.BlobItems {
			name := *item.Name
			res[name] = true
		}
	}

	return
}

/*
WriteFile writes to a file in blob storage backend
the file path itself can be nested for example doa_2022_04_01/partitions/1/p1.buf
*/
func (a *AzureBlob) WriteFile(ctx context.Context, container, file string, b []byte) (err error) {
	err = a.CreateContainer(ctx, container)
	if err != nil {
		return
	}

	// write to blob storage
	containerClient := a.serviceClient.NewContainerClient(container)
	blobClient := containerClient.NewBlockBlobClient(file)
	_, err = blobClient.UploadBufferToBlockBlob(ctx, b, azblob.HighLevelUploadToBlockBlobOption{})

	return
}

/*
LoadFile - pulls a file from blob storage
*/
func (a *AzureBlob) LoadFile(ctx context.Context, container, file string) (b []byte, err error) {
	err = a.authenticate()
	if err != nil {
		return nil, err
	}

	// initialize client
	containerClient := a.serviceClient.NewContainerClient(container)
	blobClient := containerClient.NewBlobClient(file)

	// download from client to buffer
	res, err := blobClient.Download(ctx, nil)
	if err != nil {
		return
	}
	buffer, reader := new(bytes.Buffer), res.Body(nil)
	_, err = buffer.ReadFrom(reader)
	if err == nil {
		b = buffer.Bytes()
	}

	return
}

func (a *AzureBlob) Delete(ctx context.Context, container, file string) (err error) {
	// list all the files with specified file name
	// then delete all other matches
	files, err := a.ListFiles(ctx, container, file)
	for file := range files {
		// initialize client
		containerClient := a.serviceClient.NewContainerClient(container)
		blobClient := containerClient.NewBlobClient(file)

		// delete file
		_, err = blobClient.Delete(ctx, nil)
		if err != nil {
			return
		}
	}
	return
}
