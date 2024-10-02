package main

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/pkg/errors"
	"github.com/spf13/afero"

	v1 "github.com/mistermx/crossplanereleaser/config/v1"
	"github.com/mistermx/crossplanereleaser/internal/build"
	"github.com/mistermx/crossplanereleaser/internal/git"
	"github.com/mistermx/crossplanereleaser/internal/publish"
)

type releaseCmd struct {
	git       git.Client
	builder   build.BuilderBackend
	publisher publish.PackagePublisher

	Skip string `help:"Specify steps to skip"`
}

func (c *releaseCmd) BeforeApply() error {
	c.git = git.NewGitCLIBackend()
	c.builder = build.NewCrankCLIBackend()
	c.publisher = publish.NewCraneCLIPublisher()
	return nil
}

func (c *releaseCmd) Run(fsys afero.Fs) error {
	ctx := context.Background()
	cfg, err := getConfig(fsys, c.git)
	if err != nil {
		return err
	}

	// NOTE: This implementation does not scale very well for many build
	//       pipeline steps.
	//       If this gets extended in the future we could think about
	//       implementing a build pipeline like Goreleaser:
	//       https://github.com/goreleaser/goreleaser/blob/bba4ee2be7fa0f16b682aceef3500f608f5bf18e/internal/pipeline/pipeline.go

	if c.Skip != "build" {
		cmd := buildCmd{
			git:     c.git,
			builder: c.builder,
		}
		if err := cmd.Run(fsys); err != nil {
			return errors.Wrap(err, "build failed")
		}
	}
	return errors.Wrap(c.publishPackages(ctx, cfg), "cannot push images")
}

func (c *releaseCmd) publishPackages(ctx context.Context, cfg *v1.Config) error {
	someBuildsFailed := false
	for i, push := range cfg.Pushes {
		if err := c.publishPackage(ctx, cfg, push); err != nil {
			fmt.Fprintf(os.Stderr, "publish build %d %q failed: %s", i, push.Build, err.Error())
			someBuildsFailed = true
		}
	}
	if someBuildsFailed {
		return errors.New("not all builds and pushes were successful")
	}
	return nil
}

func (c *releaseCmd) publishPackage(ctx context.Context, cfg *v1.Config, push v1.PushConfig) error {
	build := getBuildConfigByID(cfg, push.Build)
	if build == nil {
		return errors.Errorf("no build with ID %q", push.Build)
	}
	filename := getPackageOutputPath(cfg, build)

	somePublishFailed := false
	for _, img := range push.ImageTemplates {
		ref, err := name.ParseReference(img)
		if err != nil {
			fmt.Fprintln(os.Stderr, errors.Wrap(err, "cannot parse image name").Error())
			somePublishFailed = true
			continue
		}
		if err := c.publisher.PublishPackage(ctx, filename, ref); err != nil {
			fmt.Fprintln(os.Stderr, errors.Wrapf(err, "cannot publish package %q", build.ID).Error())
			somePublishFailed = true
			continue
		}
	}
	if somePublishFailed {
		return errors.New("unable to publish all images")
	}
	return nil
}

func getBuildConfigByID(cfg *v1.Config, buildID string) *v1.BuildConfig {
	for _, b := range cfg.Builds {
		if b.ID == buildID {
			return &b
		}
	}
	return nil
}
