package publish

import (
	"fmt"
	"html/template"
	"os"
	"os/exec"

	"github.com/ldez/seihon/manifest"
	"github.com/opencontainers/go-digest"
)

// DockerPub Docker images builder and publisher.
type DockerPub struct {
	builds []*exec.Cmd
	push   []*exec.Cmd
}

// NewDockerPub Creates a new DockerPub.
func NewDockerPub(imageName, version, baseImageName string, targets map[string]ArchDescriptor, dockerfileTemplate string) (*DockerPub, error) {
	manif, err := manifest.Get(baseImageName)
	if err != nil {
		return nil, err
	}

	pub := &DockerPub{}

	for target, option := range targets {
		descriptor, err := manifest.FindManifestDescriptor(option.OS, option.GoARCH, option.Variant, manif)
		if err != nil {
			return nil, err
		}

		dockerfile := fmt.Sprintf("%s-%s-%s.Dockerfile", option.OS, option.GoARCH, option.GoARM)

		err = createDockerfile(dockerfile, baseImageName, option, descriptor.Digest, dockerfileTemplate)
		if err != nil {
			return nil, err
		}

		dBuild := exec.Command("docker", "build",
			"-t", fmt.Sprintf("%s:%s-%s", imageName, version, target),
			"-f", dockerfile,
			".")
		pub.builds = append(pub.builds, dBuild)

		dPush := exec.Command("docker", "push", fmt.Sprintf(`%s:%s-%s`, imageName, version, target))
		pub.push = append(pub.push, dPush)
	}

	return pub, nil
}

// Execute Executes commands.
func (d DockerPub) Execute(dryRun bool) error {
	for _, cmd := range d.builds {
		if err := execCmd(cmd, dryRun); err != nil {
			return fmt.Errorf("failed to build: %v: %v", cmd, err)
		}
	}

	for _, cmd := range d.push {
		if err := execCmd(cmd, dryRun); err != nil {
			return fmt.Errorf("failed to push: %v: %v", cmd, err)
		}
	}

	return nil
}

func createDockerfile(dockerfile string, baseImageName string, option ArchDescriptor, digest digest.Digest, dockerfileTemplate string) error {
	parse, err := template.New("tmpl.Dockerfile").ParseFiles(dockerfileTemplate)
	if err != nil {
		return err
	}

	data := map[string]interface{}{
		"GoOS":         option.OS,
		"GoARCH":       option.GoARCH,
		"GoARM":        option.GoARM,
		"RuntimeImage": fmt.Sprintf("%s@%s", baseImageName, digest),
	}

	file, err := os.Create(dockerfile)
	if err != nil {
		return err
	}

	return parse.Execute(file, data)
}
