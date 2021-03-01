package publish

import (
	"fmt"
	"html/template"
	"os"
	"os/exec"

	"github.com/Masterminds/sprig/v3"
	"github.com/ldez/seihon/manifest"
	"github.com/opencontainers/go-digest"
)

// DockerPub Docker images builder and publisher.
type DockerPub struct {
	builds      []*exec.Cmd
	push        []*exec.Cmd
	dockerfiles []string
}

// NewDockerPub Creates a new DockerPub.
func NewDockerPub(imageName string, versions []string, baseRuntimeImage string, targets map[string]ArchDescriptor, dockerfileTemplate string) (*DockerPub, error) {
	manif, err := manifest.Get(baseRuntimeImage)
	if err != nil {
		return nil, err
	}

	pub := &DockerPub{}

	errB := orderlyBrowse(targets, func(target string, option ArchDescriptor) error {
		descriptor, err := manifest.FindManifestDescriptor(option.OS, option.GoARCH, option.Variant, manif)
		if err != nil {
			return err
		}

		dockerfile := fmt.Sprintf("%s-%s-%s.Dockerfile", option.OS, option.GoARCH, option.GoARM)
		pub.dockerfiles = append(pub.dockerfiles, dockerfile)

		err = createDockerfile(dockerfile, baseRuntimeImage, option, descriptor.Digest, dockerfileTemplate)
		if err != nil {
			return err
		}

		args := []string{"build"}
		for _, v := range versions {
			args = append(args, "-t", fmt.Sprintf("%s:%s-%s", imageName, v, target))
		}
		args = append(args, "-f", dockerfile, ".")

		dBuild := exec.Command("docker", args...)
		pub.builds = append(pub.builds, dBuild)

		for _, v := range versions {
			dPush := exec.Command("docker", "push", fmt.Sprintf(`%s:%s-%s`, imageName, v, target))
			pub.push = append(pub.push, dPush)
		}

		return nil
	})
	if errB != nil {
		return nil, errB
	}

	return pub, nil
}

// Execute Executes commands.
func (d DockerPub) Execute(dryRun bool) error {
	for _, cmd := range d.builds {
		if err := execCmd(cmd, dryRun); err != nil {
			return fmt.Errorf("failed to build: %v: %w", cmd, err)
		}
	}

	for _, cmd := range d.push {
		if err := execCmd(cmd, dryRun); err != nil {
			return fmt.Errorf("failed to push: %v: %w", cmd, err)
		}
	}

	return nil
}

// Clean Removes generated Dockerfile.
func (d DockerPub) Clean(dryRun bool) error {
	if dryRun {
		return nil
	}

	for _, dockerfile := range d.dockerfiles {
		if err := os.Remove(dockerfile); err != nil {
			return err
		}
	}

	return nil
}

func createDockerfile(dockerfile, baseRuntimeImage string, option ArchDescriptor, objDigest digest.Digest, dockerfileTemplate string) error {
	parse, err := template.New("tmpl.Dockerfile").
		Funcs(sprig.FuncMap()).
		ParseFiles(dockerfileTemplate)
	if err != nil {
		return err
	}

	data := map[string]interface{}{
		"GoOS":         option.OS,
		"GoARCH":       option.GoARCH,
		"GoARM":        option.GoARM,
		"RuntimeImage": fmt.Sprintf("%s@%s", baseRuntimeImage, objDigest),
	}

	file, err := os.Create(dockerfile)
	if err != nil {
		return err
	}

	return parse.Execute(file, data)
}
