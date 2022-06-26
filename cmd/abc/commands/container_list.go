package commands

import (
	"github.com/bradfordwagner/go-azure-blob-cli/pkg/graceful"
	"github.com/bradfordwagner/go-azure-blob-cli/pkg/state"
	"github.com/spf13/cobra"
	"strings"
)

var containerList = &cobra.Command{
	Use:   "list",
	Short: "list all containers for a blob storage account",
	Run: func(cmd *cobra.Command, args []string) {
		ac := state.NewAppContext()
		blobListMain(ac, cmd, args)
		graceful.AwaitLogError(ac.Error)
	},
}

func blobListMain(ac *state.AppContext, cmd *cobra.Command, args []string) {
	containers, err := ac.Blob.ListContainers(ac.Context)
	if err != nil {
		ac.Error <- err
	} else {
		println(strings.Join(containers, " "))
		ac.Cancel()
	}
}
