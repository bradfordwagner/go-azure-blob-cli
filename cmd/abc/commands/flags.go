package commands

import (
	"github.com/bradfordwagner/go-azure-blob-cli/pkg/config"
	"github.com/spf13/cobra"
)

//flags
var (
	flagContainer  string
	flagRawText    string
	flagFilePath   string
	flagOutputPath string
	flagList       bool
)

func createDirectoryContainerFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&flagContainer, "container", "c", "", "sets the container to use for blob/directory commands, overrides AZURE_STORAGE_CONTAINER_NAME")
	err := cmd.RegisterFlagCompletionFunc("container", listContainersValidArgsFunction)
	if err != nil {
		panic(err)
	}
}

func resolveDirectoryContainer(c config.Config) (containerName string) {
	containerName = c.ContainerName
	if flagContainer != "" {
		containerName = flagContainer
	}
	return
}

func createRawTextFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&flagRawText, "text", "t", "", "raw text to place in file")
}
func createFilePathFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&flagFilePath, "file", "f", "", "file path")
}

func createListFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVarP(&flagList, "list", "l", false, "indicate new line per input")
}

func createOutputPathFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&flagOutputPath, "output", "o", "", "output file path")
}
