package commands

import (
	"github.com/bradfordwagner/go-azure-blob-cli/pkg/graceful"
	"github.com/bradfordwagner/go-azure-blob-cli/pkg/state"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"sort"
)

func NewContainerCommand() (container *cobra.Command) {
	container = &cobra.Command{
		Use:   "container",
		Short: "functions to help manager azure blob containers",
	}

	// add sub commands
	container.AddCommand(newContainerListCommand())
	container.AddCommand(newContainerCreateCommand())
	container.AddCommand(newContainerDeleteCommand())
	return
}

// listContainersValidArgsFunction - can be used as a valid args function to return a list of valid containers for configured blob
var listContainersValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) (opts []string, directive cobra.ShellCompDirective) {
	directive = cobra.ShellCompDirectiveDefault

	// run list for options
	ac := state.NewAppContext()
	opts, err := ac.Blob.ListContainers(ac.Context)
	if err != nil {
		logrus.WithError(err).Panic("could not list containers for completion!")
	} else {
		ac.Cancel()
	}
	sort.Strings(opts)
	graceful.AwaitLogError(ac.Cancel, ac.Error)

	return
}
