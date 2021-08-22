// Package manifest contains functions related the the Docker image manifest.
package manifest

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/docker/distribution/manifest/manifestlist"
)

const manifestPath = "./manifest.json"

// FindManifestDescriptor Finds manifest.
func FindManifestDescriptor(pOS, arch, variant string, list *manifestlist.ManifestList) (manifestlist.ManifestDescriptor, error) {
	for _, descriptor := range list.Manifests {
		if descriptor.Platform.OS == pOS &&
			descriptor.Platform.Architecture == arch &&
			descriptor.Platform.Variant == variant {
			return descriptor, nil
		}
	}

	return manifestlist.ManifestDescriptor{}, fmt.Errorf("architecture not found in manifest: %s %s %s", pOS, arch, variant)
}

// Get Gets the manifest of the baseImage.
func Get(baseImageName string) (*manifestlist.ManifestList, error) {
	if _, errExist := os.Stat(manifestPath); os.IsNotExist(errExist) {
		err := inspect(baseImageName, manifestPath)
		if err != nil {
			return nil, err
		}
	} else if errExist != nil {
		return nil, errExist
	}

	bytes, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, err
	}

	manifest := &manifestlist.ManifestList{}

	err = json.Unmarshal(bytes, manifest)
	if err != nil {
		return nil, err
	}

	return manifest, nil
}

func inspect(baseImageName, manifestPath string) error {
	cmd := exec.Command("docker", "manifest", "inspect", baseImageName)
	cmd.Env = append(os.Environ(), "DOCKER_CLI_EXPERIMENTAL=enabled")

	output, err := cmd.CombinedOutput()
	if err != nil {
		if len(output) != 0 {
			fmt.Println(string(output))
		}
		return err
	}

	return os.WriteFile(manifestPath, output, 0o666)
}
