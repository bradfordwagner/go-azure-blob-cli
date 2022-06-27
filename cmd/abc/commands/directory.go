package commands

import (
	"github.com/bradfordwagner/go-azure-blob-cli/pkg/config"
	"github.com/spf13/cobra"
)

//flags
var (
	directoryContainerFlag string
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

func createDirectoryContainerFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&directoryContainerFlag, "container", "c", "", "sets the container to use for blob/directory commands, overrides AZURE_STORAGE_CONTAINER_NAME")
	err := cmd.RegisterFlagCompletionFunc("container", listContainersValidArgsFunction)
	if err != nil {
		panic(err)
	}
}

func resolveDirectoryContainer(c config.Config) (containerName string) {
	containerName = c.ContainerName
	if directoryContainerFlag != "" {
		containerName = directoryContainerFlag
	}
	return
}
