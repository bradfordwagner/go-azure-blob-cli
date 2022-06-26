package commands

import (
	"errors"
	"github.com/bradfordwagner/go-azure-blob-cli/pkg/graceful"
	"github.com/bradfordwagner/go-azure-blob-cli/pkg/state"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func newContainerDeleteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete ${container_name} ...",
		Short: "deletes a container",
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) (opts []string, directive cobra.ShellCompDirective) {
			directive = cobra.ShellCompDirectiveDefault

			// run list for options
			ac := state.NewAppContext()
			opts, err := ac.Blob.ListContainers(ac.Context)
			if err != nil {
				logrus.WithError(err).Panic("could not list containers for completion!")
			} else {
				ac.Cancel()
			}
			graceful.AwaitLogError(ac.Cancel, ac.Error)

			return
		},
		Args: func(cmd *cobra.Command, args []string) (err error) {
			// check if container name is supplied
			if len(args) == 0 || args[0] == "" {
				err = errors.New("container name is required")
			}
			return
		},
		Run: func(cmd *cobra.Command, args []string) {
			ac := state.NewAppContext()
			containerDeleteMain(ac, cmd, args)
			graceful.AwaitLogError(ac.Cancel, ac.Error)
		},
	}
	return cmd
}

func containerDeleteMain(ac *state.AppContext, cmd *cobra.Command, args []string) {
	for i := range args {
		containerName := args[i]
		err := ac.Blob.DeleteContainer(ac.Context, containerName)
		if err != nil {
			ac.Error <- err
			return
		} else {
			logrus.WithField("container", containerName).Info("deleted container")
		}
	}
	ac.Cancel()
}
