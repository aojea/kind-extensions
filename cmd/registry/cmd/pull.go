package cmd

import (
	"os/exec"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// pullCmd represents the pull command
var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "pull a container image and store it in the local registry",
	Long:  "pull a container image and store it in the local registry",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a image argument")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return pullImage(args[0])
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)
}

func pullImage(image string) error {
	// pull image
	if err := exec.Command("docker", "pull", image).Run(); err != nil {
		return errors.Wrapf(err, "docker run error pulling image %s", image)
	}
	// rename image
	_, registryImage := splitReposSearchTerm(image)
	localImage := "localhost:5000/" + registryImage
	if err := exec.Command("docker", "tag", image, localImage).Run(); err != nil {
		return errors.Wrapf(err, "docker run error tagging image %s %s", image, registryImage)
	}
	// upload to the local registry
	if err := exec.Command("docker", "push", localImage).Run(); err != nil {
		return errors.Wrapf(err, "docker run error pushing image to %s", localImage)
	}
	return nil

}

// https://github.com/moby/moby/blob/bc6f4cc7032544553d2304a5b47ba235dbfe5b9c/registry/service.go#L149
// splitReposSearchTerm breaks a search term into an index name and remote name
func splitReposSearchTerm(reposName string) (string, string) {
	nameParts := strings.SplitN(reposName, "/", 2)
	var indexName, remoteName string
	if len(nameParts) == 1 || (!strings.Contains(nameParts[0], ".") &&
		!strings.Contains(nameParts[0], ":") && nameParts[0] != "localhost") {
		// This is a Docker Index repos (ex: samalba/hipache or ubuntu)
		// 'docker.io'
		indexName = "docker.io"
		remoteName = reposName
	} else {
		indexName = nameParts[0]
		remoteName = nameParts[1]
	}
	return indexName, remoteName
}
