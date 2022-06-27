package azureblob

import (
	"bytes"
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/dustin/go-humanize"
	"github.com/sirupsen/logrus"
	"sort"
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

func (a *AzureBlob) DeleteContainer(ctx context.Context, name string) (err error) {
	err = a.authenticate()
	if err != nil {
		return
	}
	_, err = a.serviceClient.DeleteContainer(ctx, name, nil)
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

type FileProperties struct {
	Name string
	Size uint64
}

func (fp FileProperties) LogrusStruct() map[string]interface{} {
	return map[string]interface{}{
		"size": humanize.Bytes(fp.Size),
	}
}

func (a *AzureBlob) ListFilesWithProperties(ctx context.Context, container, prefix string) (fileProperties []FileProperties, err error) {
	err = a.authenticate()
	if err != nil {
		return
	}

	// setup container client
	containerClient := a.serviceClient.NewContainerClient(container)

	// allow prefix override
	opts := &azblob.ContainerListBlobHierarchySegmentOptions{}
	if prefix != "" {
		opts.Prefix = &prefix
	}

	// pull blobs from hierarchy based on prefix if paramaterized
	hierarchy := containerClient.ListBlobsHierarchy("", opts)

	// parse into result map
	for hierarchy.NextPage(ctx) {
		b := hierarchy.PageResponse()
		for _, item := range b.Segment.BlobItems {
			fileProperties = append(fileProperties, FileProperties{
				Name: *item.Name,
				Size: uint64(*item.Properties.ContentLength),
			})
		}
	}

	// sort
	sort.Slice(fileProperties, func(i, j int) bool {
		return fileProperties[i].Name < fileProperties[j].Name
	})

	return
}

func (a *AzureBlob) ListFilesAsStrings(ctx context.Context, container, prefix string) (res []string, err error) {
	properties, err := a.ListFilesWithProperties(ctx, container, prefix)
	for _, property := range properties {
		res = append(res, property.Name)
	}
	return
}

/*
WriteFile writes to a file in blob storage backend
the file path itself can be nested for example doa_2022_04_01/partitions/1/p1.buf
*/
func (a *AzureBlob) WriteFile(ctx context.Context, container, file string, b []byte) (err error) {
	err = a.authenticate()

	// write to blob storage
	containerClient := a.serviceClient.NewContainerClient(container)
	blobClient := containerClient.NewBlockBlobClient(file)
	_, err = blobClient.UploadBufferToBlockBlob(ctx, b, azblob.HighLevelUploadToBlockBlobOption{})

	return
}

/*
WriteFile writes to a file in blob storage backend
the file path itself can be nested for example doa_2022_04_01/partitions/1/p1.buf
*/
func (a *AzureBlob) CreateDirectory(ctx context.Context, container, file string) (err error) {
	var b []byte
	err = a.WriteFile(ctx, container, file, b)
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
	//files, err := a.ListFiles(ctx, container, file)
	//for file := range files {
	//	// initialize client
	//	containerClient := a.serviceClient.NewContainerClient(container)
	//	blobClient := containerClient.NewBlobClient(file)

	//	// delete file
	//	_, err = blobClient.Delete(ctx, nil)
	//	if err != nil {
	//		return
	//	}
	//}
	return
}
