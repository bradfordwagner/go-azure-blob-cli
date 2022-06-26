package commands

import (
	"github.com/spf13/cobra"
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
