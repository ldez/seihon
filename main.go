package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ldez/seihon/publish"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/spf13/pflag"
)

type cmdOpts struct {
	imageName          string
	versions           []string
	baseRuntimeImage   string
	targets            []string
	dockerfileTemplate string
	dryRun             bool
}

func main() {
	log.SetFlags(log.Lshortfile)

	rootCmd := &cobra.Command{
		Use:     "seihon",
		Short:   "A simple tool to publish multi-arch images on the Docker Hub.",
		Long:    `A simple tool to publish multi-arch images on the Docker Hub.`,
		Version: version,
	}

	rootCmd.AddCommand(newPublishCmd())

	docCmd := &cobra.Command{
		Use:    "doc",
		Short:  "Generate documentation",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return doc.GenMarkdownTree(rootCmd, "./docs")
		},
	}

	rootCmd.AddCommand(docCmd)

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Display version",
		Run: func(_ *cobra.Command, _ []string) {
			displayVersion(rootCmd.Name())
		},
	}

	rootCmd.AddCommand(versionCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func newPublishCmd() *cobra.Command {
	opts := cmdOpts{}

	cmd := &cobra.Command{
		Use:   "publish",
		Short: "Build and publish multi-arch Docker image.",
		Long:  `Build and publish multi-arch Docker image.`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := validateRequiredFlags(cmd); err != nil {
				return err
			}

			if opts.dryRun {
				fmt.Println("IMPORTANT: you are using the dry-run mode. Use `--dry-run=false` to disable this mode.")
			}

			// TODO add an option?
			// _, travisTag := os.LookupEnv("TRAVIS_TAG")
			// if !travisTag {
			// 	log.Println("Skipping deploy")
			// 	os.Exit(0)
			// }

			return nil
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			return run(opts)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.imageName, "image-name", "i", "", "Image name (user/repo)")
	flags.StringSliceVarP(&opts.versions, "versions", "v", nil, "Image version.")
	flags.StringVar(&opts.dockerfileTemplate, "template", "./tmpl.Dockerfile", "Dockerfile template")
	flags.StringVarP(&opts.baseRuntimeImage, "base-runtime-image", "b", "alpine:3.10", "Base Docker image.")
	flags.StringSliceVar(&opts.targets, "targets", []string{"arm.v6", "arm.v7", "arm.v8", "amd64", "386"}, "Targeted architectures.")
	flags.BoolVar(&opts.dryRun, "dry-run", true, "Dry run mode.")

	return cmd
}

func run(opts cmdOpts) error {
	targetedArch, err := publish.GetTargetedArchitectures(opts.targets)
	if err != nil {
		return err
	}

	dockerPub, err := publish.NewDockerPub(opts.imageName, opts.versions, opts.baseRuntimeImage, targetedArch, opts.dockerfileTemplate)
	if err != nil {
		return err
	}

	if err = dockerPub.Execute(opts.dryRun); err != nil {
		return err
	}

	err = dockerPub.Clean(opts.dryRun)
	if err != nil {
		return err
	}

	for _, version := range opts.versions {
		manifestPub, err := publish.NewManifestPub(opts.imageName, version, targetedArch)
		if err != nil {
			return err
		}

		if err = manifestPub.Execute(opts.dryRun); err != nil {
			return err
		}
	}

	return nil
}

func validateRequiredFlags(cmd *cobra.Command) error {
	var missingFlagNames []string

	flags := cmd.Flags()
	flags.
		VisitAll(func(pflag *pflag.Flag) {
			switch pflag.Value.Type() {
			case "string":
				if len(pflag.Value.String()) == 0 {
					missingFlagNames = append(missingFlagNames, pflag.Name)
				}
			case "stringSlice":
				slice, _ := flags.GetStringSlice(pflag.Name)
				if len(slice) == 0 {
					missingFlagNames = append(missingFlagNames, pflag.Name)
				}
			}
		})

	if len(missingFlagNames) > 0 {
		return fmt.Errorf(`required flag(s) "%s" not set`, strings.Join(missingFlagNames, `", "`))
	}
	return nil
}
