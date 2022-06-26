package commands

import (
	"errors"
	"github.com/bradfordwagner/go-azure-blob-cli/pkg/graceful"
	"github.com/bradfordwagner/go-azure-blob-cli/pkg/state"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var containerCreate = &cobra.Command{
	Use:   "create ${container_name}",
	Short: "creates a container if does not exist",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		// check if container name is supplied
		if len(args) == 0 || args[0] == "" {
			err = errors.New("container name is required")
		}
		return
	},
	Run: func(cmd *cobra.Command, args []string) {
		ac := state.NewAppContext()
		containerCreateMain(ac, cmd, args)
		graceful.AwaitLogError(ac.Cancel, ac.Error)
	},
}

func containerCreateMain(ac *state.AppContext, cmd *cobra.Command, args []string) {
	containerName := args[0]
	err := ac.Blob.CreateContainer(ac.Context, containerName)
	if err != nil {
		ac.Error <- err
	} else {
		logrus.WithField("container", containerName).Info("created container")
		ac.Cancel()
	}
}
