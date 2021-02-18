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
	dockerRegistryProxyImage = "ghcr.io/rpardini/docker-registry-proxy:0.6.2"
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
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			return err
		}
		return createRegistry(name)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().String(
		"name",
		cluster.DefaultName,
		"the cluster context name",
	)
	// TODO: add authentication
}

func createRegistry(name string) error {
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
	fmt.Println(nodeList)
	return nil
}

// https://github.com/rpardini/docker-registry-proxy#simple-no-auth-all-cache
func createContainer(name string) error {
	args := []string{
		"--rm",
		"--net", kindNetwork,
		"--name", fmt.Sprintf("docker-registry-proxy-%s", name),
		"--label", fmt.Sprintf("%s=%s", clusterLabelKey, name),
		"-e", "ENABLE_MANIFEST_CACHE=true",
		"-e", "REGISTRIES=\"k8s.gcr.io gcr.io quay.io your.own.registry another.public.registry\"",
		"--restart=on-failure:1",
		dockerRegistryProxyImage,
	}

	if err := exec.Command("docker", args...).Run(); err != nil {
		return errors.Wrap(err, "docker run error")
	}
	return nil
}

func configureCluster(nodes []kindnodes.Node) error {
	return nil
}
