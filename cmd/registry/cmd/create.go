package cmd

import (
	"fmt"
	"os/exec"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"sigs.k8s.io/kind/pkg/cluster"
	kindnodes "sigs.k8s.io/kind/pkg/cluster/nodes"
	kindcmd "sigs.k8s.io/kind/pkg/cmd"
)

const (
	dockerRegistryProxyImage = "rpardini/docker-registry-proxy:0.6.3"
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
	createCmd.Flags().String(
		"name",
		cluster.DefaultName,
		"the cluster context name",
	)
	createCmd.Flags().Bool(
		"retain",
		false,
		"don't clean the registry container",
	)
	// TODO: add authentication
}

func createRegistry(cmd *cobra.Command) error {
	name, err := cmd.Flags().GetString("name")
	if err != nil {
		return err
	}
	retain, err := cmd.Flags().GetBool("retain")
	if err != nil {
		return err
	}

	logger := kindcmd.NewLogger()
	provider := cluster.NewProvider(
		cluster.ProviderWithLogger(logger),
	)

	// Check if the cluster nodes exist
	nodeList, err := provider.ListInternalNodes(name)
	if err != nil {
		return err
	}
	if len(nodeList) == 0 {
		return fmt.Errorf("no nodes found for cluster %q", name)
	}
	err = createContainer(name, retain)
	if err != nil {
		return err
	}
	return nil
}

// https://github.com/rpardini/docker-registry-proxy#simple-no-auth-all-cache
func createContainer(name string, retain bool) error {
	containerName := fmt.Sprintf("docker-registry-proxy-%s", name)
	args := []string{"run",
		"-d",                 // run in the background
		"--net", kindNetwork, // attach to the KIND network
		"--name", containerName, // well known name
		"--label", fmt.Sprintf("%s=%s", clusterLabelKey, name), // label as a KIND cluster node
		"-e", "ENABLE_MANIFEST_CACHE=true",
		"-e", "REGISTRIES=k8s.gcr.io gcr.io quay.io", // TODO: pass environment variables directly
		"--restart=on-failure:1", // same as KIND
		dockerRegistryProxyImage,
	}

	if err := exec.Command("docker", args...).Run(); err != nil {
		// try to clean as much as possible if retain not enabled
		if !retain {
			exec.Command("docker", "rm", "-f", containerName).Run()
		}
		return errors.Wrap(err, "docker run error")
	}
	return nil
}

func configureCluster(nodes []kindnodes.Node) error {
	return nil
}
