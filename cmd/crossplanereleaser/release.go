package main

import (
	"context"

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
	return c.publishPackages(ctx, cfg)
}

func (c *releaseCmd) publishPackages(ctx context.Context, cfg *v1.Config) error {
	for _, pkgCfg := range cfg.XPackages {
		filename := getPackageOutputPath(cfg, &pkgCfg)
		ref, err := name.ParseReference(pkgCfg.NameTemplate)
		if err != nil {
			return errors.Wrap(err, "cannot parse image name")
		}
		if err := c.publisher.PublishPackage(ctx, filename, ref); err != nil {
			return errors.Wrapf(err, "cannot publish package %q", pkgCfg.ID)
		}
	}
	return nil
}
