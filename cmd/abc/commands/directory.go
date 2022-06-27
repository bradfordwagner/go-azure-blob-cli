package commands

import (
	"github.com/spf13/cobra"
)

func NewDirectoryCommand() (container *cobra.Command) {
	container = &cobra.Command{
		Use:   "directory",
		Short: "functions to help manager azure blob directories",
	}

	// add sub commands
	container.AddCommand(newDirectoryCreateCommand())
	return
}
