package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"sigs.k8s.io/kind/pkg/cmd"
	"sigs.k8s.io/kind/pkg/cmd/kind/create"
)

var rootCmd = &cobra.Command{
	Use:   "zones",
	Short: "Zones create a KIND cluster emulating AZ zones",
	Long:  "Zones create a KIND cluster emulating AZ zones",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func Execute() {
	logger := cmd.NewLogger()
	streams := cmd.StandardIOStreams()
	rootCmd.AddCommand(create.NewCommand(logger, streams))
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}
