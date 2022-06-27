package commands

import (
	"github.com/bradfordwagner/go-azure-blob-cli/pkg/graceful"
	"github.com/bradfordwagner/go-azure-blob-cli/pkg/state"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"sort"
)

func NewBlobCommand() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "blob",
		Short: "functions to help manager azure blobs",
	}

	// global flags - for all blob subcommands
	createDirectoryContainerFlag(cmd)

	// add sub commands
	cmd.AddCommand(newBlobUploadCommand())
	cmd.AddCommand(newBlobListCommand())
	cmd.AddCommand(newBlobDownloadCommand())
	cmd.AddCommand(newBlobDeleteCommand())

	return
}

// listBlobsValidArgsFunction - can be used as a valid args function to return a list of valid blobs for a configured container
var listBlobsValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) (opts []string, directive cobra.ShellCompDirective) {
	directive = cobra.ShellCompDirectiveDefault

	// resolve path prefix
	pathPrefix := ""
	if len(args) > 0 {
		pathPrefix = args[0]
	}

	// run list for options
	ac := state.NewAppContext()
	containerName := resolveDirectoryContainer(ac.Config)
	opts, err := ac.Blob.ListFilesAsStrings(ac.Context, containerName, pathPrefix)
	if err != nil {
		logrus.WithError(err).Panic("could not list containers for completion!")
	} else {
		ac.Cancel()
	}
	sort.Strings(opts)
	graceful.AwaitLogError(ac.Cancel, ac.Error)

	return
}
