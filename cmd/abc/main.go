package main

import (
	"github.com/bradfordwagner/go-azure-blob-cli/cmd/abc/commands"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "abc",
}

func init() {
	rootCmd.AddCommand(commands.NewContainerCommand())
}

func main() {
	// setup logrus
	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetLevel(logrus.InfoLevel)

	// cobra
	if err := rootCmd.Execute(); err != nil {
		logrus.WithError(err).Error("execution failed")
		os.Exit(1)
	}
}
