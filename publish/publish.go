package publish

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
)

const envDockerExperimental = "DOCKER_CLI_EXPERIMENTAL=enabled"

var availableArchitectures = map[string]ArchDescriptor{
	"386": {
		OS:     "linux",
		GoARCH: "386",
	},
	"amd64": {
		OS:     "linux",
		GoARCH: "amd64",
	},
	"arm.v5": {
		OS:      "linux",
		GoARCH:  "arm",
		GoARM:   "5",
		Variant: "v5",
	},
	"arm.v6": {
		OS:      "linux",
		GoARCH:  "arm",
		GoARM:   "6",
		Variant: "v6",
	},
	"arm.v7": {
		OS:      "linux",
		GoARCH:  "arm",
		GoARM:   "7",
		Variant: "v7",
	},
	"arm.v8": {
		OS:      "linux",
		GoARCH:  "arm64",
		Variant: "v8",
	},
}

// ArchDescriptor An architecture descriptor for an architecture.
type ArchDescriptor struct {
	OS      string `json:"os"`
	GoARCH  string `json:"go_arch"`
	GoARM   string `json:"go_arm,omitempty"`
	Variant string `json:"variant,omitempty"`
}

// GetTargetedArchitectures Gets architecture descriptors.
func GetTargetedArchitectures(targets []string) (map[string]ArchDescriptor, error) {
	targetedArch := make(map[string]ArchDescriptor)

	for _, target := range targets {
		option, ok := availableArchitectures[target]
		if !ok {
			return nil, fmt.Errorf("unsupported platform: %s", target)
		}
		targetedArch[target] = option
	}

	return targetedArch, nil
}

func orderlyBrowse(targets map[string]ArchDescriptor, apply func(string, ArchDescriptor) error) error {
	var keys []string

	for target := range targets {
		keys = append(keys, target)
	}

	sort.Strings(keys)

	for _, key := range keys {
		err := apply(key, targets[key])
		if err != nil {
			return err
		}
	}

	return nil
}

func execCmd(cmd *exec.Cmd, dryRun bool) error {
	cmd.Env = append(os.Environ(), envDockerExperimental)
	if dryRun {
		fmt.Println(strings.Join(cmd.Args, " "))
		return nil
	}

	output, err := cmd.CombinedOutput()

	if len(output) != 0 {
		fmt.Println(string(output))
	}

	if err != nil {
		return err
	}
	return nil
}
