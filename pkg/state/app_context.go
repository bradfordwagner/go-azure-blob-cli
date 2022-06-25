package state

import (
	"context"
	"github.com/bradfordwagner/go-azure-blob-cli/pkg/azureblob"
	"github.com/bradfordwagner/go-azure-blob-cli/pkg/config"
	"github.com/bradfordwagner/go-azure-blob-cli/pkg/graceful"
)

// AppContext - runtime app config
type AppContext struct {
	Config  config.Config
	Context context.Context
	Cancel  context.CancelFunc
	Error   chan error
	Blob    *azureblob.AzureBlob
}

// NewAppContext - creates a new app context
func NewAppContext() *AppContext {
	ctx, cancel, errChan := graceful.New()
	c := config.New()
	return &AppContext{
		Config:  c,
		Context: ctx,
		Cancel:  cancel,
		Error:   errChan,
		Blob:    azureblob.NewAzureBlob(c.StorageAccountName, c.StorageAccountKey),
	}
}
