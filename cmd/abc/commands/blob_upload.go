package commands

import (
	"errors"
	"github.com/bradfordwagner/go-azure-blob-cli/pkg/graceful"
	"github.com/bradfordwagner/go-azure-blob-cli/pkg/state"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

func newBlobUploadCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upload targetPath # with either -t, or -f",
		Short: "uploads data into ",
		Args: func(cmd *cobra.Command, args []string) (err error) {
			if len(args) == 0 {
				err = errors.New("target file path is required as first argument")
			} else if flagFilePath == "" && flagRawText == "" {
				err = errors.New("-t, or -f are required")
			} else if flagFilePath != "" && flagRawText != "" {
				err = errors.New("one of -t, or -f required cannot use both")
			}
			return
		},
		ValidArgsFunction: listBlobsValidArgsFunction,
		Run: func(cmd *cobra.Command, args []string) {
			ac := state.NewAppContext()
			blobUploadMain(ac, cmd, args)
			graceful.AwaitLogError(ac.Cancel, ac.Error)
		},
	}

	// flags
	createDirectoryContainerFlag(cmd)
	createRawTextFlag(cmd)
	createFilePathFlag(cmd)

	return cmd
}

func blobUploadMain(ac *state.AppContext, cmd *cobra.Command, args []string) {
	containerName := resolveDirectoryContainer(ac.Config)
	filePath := args[0]
	b, err := resolveUploadBytes()
	if err != nil {
		ac.Error <- err
		return
	}

	err = ac.Blob.WriteFile(ac.Context, containerName, filePath, b)
	if err != nil {
		ac.Error <- err
	} else {
		logrus.WithFields(map[string]interface{}{
			"container":   containerName,
			"upload_path": filePath,
		}).Info("successfully wrote file")
	}
	ac.Cancel()
}

func resolveUploadBytes() (b []byte, err error) {
	if flagRawText != "" {
		b = []byte(flagRawText)
	}
	if flagFilePath != "" {
		b, err = os.ReadFile(flagFilePath)
	}
	return
}
