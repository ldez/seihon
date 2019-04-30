package publish

import (
	"fmt"
	"os/exec"
)

// ManifestPub Manifest builder and publisher.
type ManifestPub struct {
	manifestAnnotate []*exec.Cmd
	manifestCreate   *exec.Cmd
	manifestPush     *exec.Cmd
}

// NewManifestPub Creates a new ManifestPub.
func NewManifestPub(imageName, version string, targets map[string]ArchDescriptor) (*ManifestPub, error) {
	pub := &ManifestPub{}

	_ = orderlyBrowse(targets, func(target string, option ArchDescriptor) error {
		ma := []string{
			"manifest", "annotate",
			fmt.Sprintf("%s:%s", imageName, version),
			fmt.Sprintf("%s:%s-%s", imageName, version, target),
			fmt.Sprintf("--os=%s", option.OS),
			fmt.Sprintf("--arch=%s", option.GoARCH),
		}
		if option.Variant != "" {
			ma = append(ma, fmt.Sprintf("--variant=%s", option.Variant))
		}

		cmdMA := exec.Command("docker", ma...)
		cmdMA.Env = append(cmdMA.Env, envDockerExperimental)
		pub.manifestAnnotate = append(pub.manifestAnnotate, cmdMA)

		return nil
	})

	mc := []string{"manifest", "create", "--amend", fmt.Sprintf("%s:%s", imageName, version)}
	_ = orderlyBrowse(targets, func(target string, option ArchDescriptor) error {
		mc = append(mc, fmt.Sprintf("%s:%s-%s", imageName, version, target))
		return nil
	})

	cmdMC := exec.Command("docker", mc...)
	cmdMC.Env = append(cmdMC.Env, envDockerExperimental)
	pub.manifestCreate = cmdMC

	cmdMP := exec.Command("docker", "manifest", "push", fmt.Sprintf("%s:%s", imageName, version))
	cmdMP.Env = append(cmdMP.Env, envDockerExperimental)
	pub.manifestPush = cmdMP

	return pub, nil
}

// Execute Executes commands.
func (m ManifestPub) Execute(dryRun bool) error {
	if err := execCmd(m.manifestCreate, dryRun); err != nil {
		return fmt.Errorf("failed to create manifest: %v: %v", m.manifestCreate, err)
	}

	for _, cmd := range m.manifestAnnotate {
		if err := execCmd(cmd, dryRun); err != nil {
			return fmt.Errorf("failed to annotate manifest: %v: %v", cmd, err)
		}
	}

	if err := execCmd(m.manifestPush, dryRun); err != nil {
		return fmt.Errorf("failed to push manifest: %v: %v", m.manifestPush, err)
	}

	return nil
}
