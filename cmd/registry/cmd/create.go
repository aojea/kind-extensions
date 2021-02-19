package cmd

import (
	"fmt"
	"os/exec"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	dockerRegistryImage = "registry:2"
	// clusterLabelKey is applied to each "node" docker container for identification
	clusterLabelKey = "io.x-k8s.kind.extension.cluster"
	kindNetwork     = "kind"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a container registry for the specified KIND cluster",
	Long:  "Create a container registry for the specified KIND cluster",
	RunE: func(cmd *cobra.Command, args []string) error {
		return createRegistry(cmd)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().Bool(
		"retain",
		false,
		"don't clean the registry container",
	)
}

func createRegistry(cmd *cobra.Command) error {
	retain, err := cmd.Flags().GetBool("retain")
	if err != nil {
		return err
	}

	args := []string{"run",
		"-d",                 // run in the background
		"--net", kindNetwork, // attach to the KIND network
		"--name", registry, // well known name
		"-p", "5000:5000", // listen on port 5000
		"--label", fmt.Sprintf("%s=%s", clusterLabelKey, registry), // label as a KIND cluster node
		"--restart=always",
		dockerRegistryImage,
	}

	if err := exec.Command("docker", args...).Run(); err != nil {
		// try to clean as much as possible if retain not enabled
		if !retain {
			exec.Command("docker", "rm", "-f", registry).Run()
		}
		return errors.Wrap(err, "docker run error")
	}
	return nil

}
