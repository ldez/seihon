package publish

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

const envDockerExperimental = "DOCKER_CLI_EXPERIMENTAL=enabled"

var buildOptions = map[string]buildOption{
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

type buildOption struct {
	OS      string `json:"os"`
	GoARCH  string `json:"go_arch"`
	GoARM   string `json:"go_arm,omitempty"`
	Variant string `json:"variant,omitempty"`
}

func execCmd(cmd *exec.Cmd, dryRun bool) error {
	if dryRun {
		fmt.Println(cmd.Path, strings.Join(cmd.Args, " "))
		return nil
	}

	output, err := cmd.CombinedOutput()

	if len(output) != 0 {
		log.Println(string(output))
	}

	if err != nil {
		return err
	}
	return nil
}
