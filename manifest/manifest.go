package manifest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/docker/distribution/manifest/manifestlist"
)

// FindManifestDescriptor Finds manifest.
func FindManifestDescriptor(os, arch, variant string, list *manifestlist.ManifestList) (manifestlist.ManifestDescriptor, error) {
	for _, descriptor := range list.Manifests {
		if descriptor.Platform.OS == os &&
			descriptor.Platform.Architecture == arch &&
			descriptor.Platform.Variant == variant {
			return descriptor, nil
		}
	}

	return manifestlist.ManifestDescriptor{}, fmt.Errorf("not supported: %s %s %s", os, arch, variant)
}

// Get Gets the manifest of the baseImage.
func Get(baseImageName string) (*manifestlist.ManifestList, error) {
	manifestPath := "./manifest.json"

	if _, errExist := os.Stat(manifestPath); os.IsNotExist(errExist) {
		err := inspect(baseImageName, manifestPath)
		if err != nil {
			return nil, err
		}
	} else if errExist != nil {
		return nil, errExist
	}

	bytes, err := ioutil.ReadFile(manifestPath)
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
	cmd.Env = append(cmd.Env, "DOCKER_CLI_EXPERIMENTAL=enabled")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	return ioutil.WriteFile(manifestPath, output, 0666)
}