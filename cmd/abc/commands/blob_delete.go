package commands

import (
	"errors"
	"github.com/bradfordwagner/go-azure-blob-cli/pkg/graceful"
	"github.com/bradfordwagner/go-azure-blob-cli/pkg/state"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func newBlobDeleteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "delete ${path}",
		Short:             "deletes all data from a given path including subpaths, * for all",
		ValidArgsFunction: listBlobsValidArgsFunction,
		Run: func(cmd *cobra.Command, args []string) {
			ac := state.NewAppContext()
			blobDeleteMain(ac, cmd, args)
			graceful.AwaitLogError(ac.Cancel, ac.Error)
		},
	}

	// flags
	createDirectoryContainerFlag(cmd)

	return cmd
}

func blobDeleteMain(ac *state.AppContext, cmd *cobra.Command, args []string) {
	containerName := resolveDirectoryContainer(ac.Config)
	var (
		blobPath string
		err      error
	)
	if len(args) == 0 {
		err = errors.New("no path for deletion provided")
		return
	} else {
		blobPath = args[0]
		if args[0] == "*" {
			blobPath = "" // delete everything AHHHHHHHH
		}
		err = ac.Blob.Delete(ac.Context, containerName, blobPath)
	}

	if err != nil {
		ac.Error <- err
	} else {
		logrus.Info("deleted: ", blobPath)
		ac.Cancel()
	}
}
