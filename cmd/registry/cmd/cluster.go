package cmd

import (
	"bytes"
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"sigs.k8s.io/kind/pkg/cluster"
	kindcmd "sigs.k8s.io/kind/pkg/cmd"
)

// clusterCmd represents the create command
var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Configure a KIND cluster to use the local registry",
	Long:  "Configure a KIND cluster to use the local registry",
	RunE: func(cmd *cobra.Command, args []string) error {
		return configureCluster(cmd)
	},
}

func init() {
	rootCmd.AddCommand(clusterCmd)
	clusterCmd.Flags().String(
		"name",
		cluster.DefaultName,
		"the cluster context name",
	)

}

func configureCluster(cmd *cobra.Command) error {
	// containers reach the registry directly
	registryURL := fmt.Sprintf("http://%s:5000", registry)

	name, err := cmd.Flags().GetString("name")
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

	// configure local registry as a mirror in all nodes
	mirrorRegistry := `
[plugins."io.containerd.grpc.v1.cri".registry]
	[plugins."io.containerd.grpc.v1.cri".registry.mirrors]
	   [plugins."io.containerd.grpc.v1.cri".registry.mirrors."*"]
	   endpoint = ["` + registryURL + `"]
   [plugins."io.containerd.grpc.v1.cri".registry.configs]
	   [plugins."io.containerd.grpc.v1.cri".registry.configs."*"]
			[plugins."io.containerd.grpc.v1.cri".registry.configs."*".tls]
			  insecure_skip_verify = true 
`

	for _, n := range nodeList {
		// Get current containerd configuration
		var buff bytes.Buffer
		if err := n.Command("cat", "/etc/containerd/config.toml").SetStdout(&buff).Run(); err != nil {
			return errors.Wrapf(err, "failed to read containerd config from node %s", n.String())
		}
		// Append the mirror registry
		// TODO this may use patching or detect if exists already
		buff.WriteString(mirrorRegistry)
		// Write the new configuration
		if err := n.Command("cp", "/dev/stdin", "/etc/containerd/config.toml").SetStdin(&buff).Run(); err != nil {
			return errors.Wrapf(err, "failed to write containerd config to node %s", n.String())
		}

		if err := n.Command("systemctl", "restart", "containerd").Run(); err != nil {
			return errors.Wrap(err, "failed to restart containerd")
		}
	}

	return nil
}
