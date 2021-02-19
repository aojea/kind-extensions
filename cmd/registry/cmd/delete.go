package cmd

import (
	"os/exec"

	"github.com/spf13/cobra"
	"sigs.k8s.io/kind/pkg/cluster"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Create a container registry for the specified KIND cluster",
	Long:  "Create a container registry for the specified KIND cluster",
	RunE: func(cmd *cobra.Command, args []string) error {
		return exec.Command("docker", "rm", "-f", registry).Run()
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().String(
		"name",
		cluster.DefaultName,
		"the cluster context name",
	)
}
