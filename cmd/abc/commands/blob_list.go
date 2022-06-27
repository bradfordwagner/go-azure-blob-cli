package commands

import (
	"errors"
	"github.com/bradfordwagner/go-azure-blob-cli/pkg/graceful"
	"github.com/bradfordwagner/go-azure-blob-cli/pkg/state"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"strings"
)

func newBlobListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list ${filePath} # empty is root",
		Short: "list files in path",
		Run: func(cmd *cobra.Command, args []string) {
			ac := state.NewAppContext()
			blobListMain(ac, cmd, args)
			graceful.AwaitLogError(ac.Cancel, ac.Error)
		},
	}

	// flags
	createDirectoryContainerFlag(cmd)
	createListFlag(cmd)

	return cmd
}

func blobListMain(ac *state.AppContext, cmd *cobra.Command, args []string) {
	containerName := resolveDirectoryContainer(ac.Config)
	if containerName == "" {
		ac.Error <- errors.New("container name must be provided with -c or AZURE_STORAGE_CONTAINER_NAME")
		return
	}

	path := ""
	if len(args) > 0 {
		path = args[0]
	}

	if flagList {
		// detailed info
		files, err := ac.Blob.ListFilesWithProperties(ac.Context, containerName, path)
		if err != nil {
			ac.Error <- err
			return
		}
		for _, file := range files {
			logrus.WithFields(file.LogrusStruct()).Info(file.Name)
		}
	} else {
		// default standard list
		files, err := ac.Blob.ListFilesAsStrings(ac.Context, containerName, path)
		if err != nil {
			ac.Error <- err
			return
		}
		println(strings.Join(files, " "))
	}

	ac.Cancel()
}
