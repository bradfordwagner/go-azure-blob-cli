package commands

import (
	"github.com/bradfordwagner/go-azure-blob-cli/pkg/graceful"
	"github.com/bradfordwagner/go-azure-blob-cli/pkg/state"
	"github.com/spf13/cobra"
)

func newBlobDownloadCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "download ${path}",
		Short:             "downloads data from the blob to local dir",
		ValidArgsFunction: listBlobsValidArgsFunction,
		Run: func(cmd *cobra.Command, args []string) {
			ac := state.NewAppContext()
			blobDownloadMain(ac, cmd, args)
			graceful.AwaitLogError(ac.Cancel, ac.Error)
		},
	}

	// flags
	createDirectoryContainerFlag(cmd)
	createOutputPathFlag(cmd)

	return cmd
}

func blobDownloadMain(ac *state.AppContext, cmd *cobra.Command, args []string) {
	containerName := resolveDirectoryContainer(ac.Config)
	var err error
	if len(args) == 0 {
		err = ac.Blob.Download(ac.Context, containerName, "", flagOutputPath)
	} else {
		for _, filePath := range args {
			err = ac.Blob.Download(ac.Context, containerName, filePath, flagOutputPath)
		}
	}

	if err != nil {
		ac.Error <- err
		return
	}

	ac.Cancel()
}
