package commands

import (
	"github.com/bradfordwagner/go-azure-blob-cli/pkg/graceful"
	"github.com/bradfordwagner/go-azure-blob-cli/pkg/state"
	"github.com/spf13/cobra"
	"sort"
	"strings"
)

func newContainerListCommand() *cobra.Command {
	containerList := &cobra.Command{
		Use:   "list",
		Short: "list all containers for a blob storage account",
		Run: func(cmd *cobra.Command, args []string) {
			ac := state.NewAppContext()
			containerListMain(ac, cmd, args)
			graceful.AwaitLogError(ac.Cancel, ac.Error)
		},
	}

	// flags

	return containerList
}

func containerListMain(ac *state.AppContext, cmd *cobra.Command, args []string) {
	containers, err := ac.Blob.ListContainers(ac.Context)
	if err != nil {
		ac.Error <- err
	} else {
		sort.Strings(containers)
		if flagList {
			newLineJoin(containers)
		} else {
			spaceJoin(containers)
		}
		ac.Cancel()
	}
}

func spaceJoin(containers []string) {
	println(strings.Join(containers, " "))
}

func newLineJoin(containers []string) {
	println(strings.Join(containers, "\n"))
}
