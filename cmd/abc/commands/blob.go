package commands

import (
	"github.com/spf13/cobra"
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

	return
}
