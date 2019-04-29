package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ldez/seihon/publish"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

type cmdOpts struct {
	imageName          string
	version            string
	baseImageName      string
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
		PreRun: func(cmd *cobra.Command, args []string) {
			requireString("image-name", cmd)
			requireString("version", cmd)
			requireString("base-image-name", cmd)
			requireString("template", cmd)

			// _, travisTag := os.LookupEnv("TRAVIS_TAG")
			// if !travisTag {
			// 	log.Println("Skipping deploy")
			// 	os.Exit(0)
			// }
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			return run(opts)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&opts.imageName, "image-name", "ldez/traefik-certs-dumper", "Image name (user/repo)")
	flags.StringVar(&opts.version, "version", "", "Image version.")
	flags.StringVar(&opts.dockerfileTemplate, "template", "./tmpl.Dockerfile", "Dockerfile template")
	flags.StringVar(&opts.baseImageName, "base-image-name", "alpine:3.9", "Base Docker image.")
	flags.StringSliceVar(&opts.targets, "targets", []string{"arm.v6", "arm.v7", "arm.v8", "amd64", "386"}, "Targeted architectures.")
	flags.BoolVar(&opts.dryRun, "dry-run", true, "Dry run mode.")

	return cmd
}

func requireString(fieldName string, cmd *cobra.Command) {
	if cmd.Flag(fieldName) == nil || cmd.Flag(fieldName).Value.String() == "" {
		log.Fatalf("%s is required", fieldName)
	}
}

func run(opts cmdOpts) error {
	dockerPub, err := publish.NewDockerPub(opts.imageName, version, opts.baseImageName, opts.targets, opts.dockerfileTemplate)
	if err != nil {
		return err
	}

	err = dockerPub.Execute(opts.dryRun)
	if err != nil {
		return err
	}

	manifestPub, err := publish.NewManifestPub(opts.imageName, version, opts.targets)
	if err != nil {
		return err
	}

	return manifestPub.Execute(opts.dryRun)
}
