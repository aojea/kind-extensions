package cmd

import (
	"fmt"
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
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			return err
		}
		containerName := fmt.Sprintf("docker-registry-proxy-%s", name)
		return exec.Command("docker", "rm", "-f", containerName).Run()
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

func deleteRegistry(cmd *cobra.Command) error {

	return nil
}
