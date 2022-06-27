package commands

import (
	"errors"
	"github.com/bradfordwagner/go-azure-blob-cli/pkg/graceful"
	"github.com/bradfordwagner/go-azure-blob-cli/pkg/state"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func newDirectoryCreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create ${directory_name} ...",
		Short: "creates a directory if does not exist",
		Args: func(cmd *cobra.Command, args []string) (err error) {
			// check if container name is supplied
			if len(args) == 0 || args[0] == "" {
				err = errors.New("directory name is required")
			}
			return
		},
		Run: func(cmd *cobra.Command, args []string) {
			ac := state.NewAppContext()
			directoryCreateMain(ac, cmd, args)
			graceful.AwaitLogError(ac.Cancel, ac.Error)
		},
	}

	// flags
	createDirectoryContainerFlag(cmd)

	return cmd
}

func directoryCreateMain(ac *state.AppContext, cmd *cobra.Command, args []string) {
	containerName := resolveDirectoryContainer(ac.Config)
	for i := range args {
		directoryName := args[i]
		err := ac.Blob.CreateDirectory(ac.Context, containerName, directoryName)
		if err != nil {
			ac.Error <- err
			return
		} else {
			logrus.WithFields(map[string]interface{}{
				"container": containerName,
				"directory": directoryName,
			}).Info("created directory")
		}
	}
	ac.Cancel()
}
