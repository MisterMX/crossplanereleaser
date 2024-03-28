package main

import (
	"context"
	"path/filepath"
	"slices"

	"github.com/pkg/errors"
	"github.com/spf13/afero"

	v1 "github.com/mistermx/crossplanereleaser/config/v1"
	"github.com/mistermx/crossplanereleaser/internal/build"
	"github.com/mistermx/crossplanereleaser/internal/git"
)

type buildCmd struct {
	ID []string `help:"Build only the specified build ids"`

	git     git.Client
	builder build.BuilderBackend
}

func (c *buildCmd) BeforeApply() error {
	c.git = git.NewGitCLIBackend()
	c.builder = build.NewCrankCLIBackend()
	return nil
}

func (c *buildCmd) Run(fsys afero.Fs) error {
	ctx := context.Background()

	cfg, err := getConfig(fsys, c.git)
	if err != nil {
		return err
	}
	return c.buildPackages(ctx, fsys, cfg)
}

func (c *buildCmd) buildPackages(ctx context.Context, fsys afero.Fs, cfg *v1.Config) error {
	for _, b := range cfg.Builds {
		if !c.shouldBeBuilt(b) {
			continue
		}

		buildCfg := &build.PackageBuildConfig{
			PackageDir:      b.Dir,
			ExamplesDir:     b.Examples,
			OutputPath:      getPackageOutputPath(cfg, &b),
			RuntimeImageTar: b.RuntimeImageTar,
		}
		// We need to call mkdir here because crank doesn't do it for us
		if err := fsys.MkdirAll(filepath.Dir(buildCfg.OutputPath), 0755); err != nil {
			return err
		}
		err := c.builder.BuildPackage(ctx, buildCfg)
		if err != nil {
			return errors.Wrapf(err, "cannot build package %q", b.ID)
		}
	}
	return nil
}

func (c *buildCmd) shouldBeBuilt(buildCfg v1.BuildConfig) bool {
	if len(c.ID) == 0 {
		return true
	}
	return slices.Contains(c.ID, buildCfg.ID)
}

func getPackageOutputPath(cfg *v1.Config, build *v1.BuildConfig) string {
	return filepath.Join(cfg.Dist, build.ID, build.NameTemplate)
}
