package main

import (
	"context"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/afero"

	v1 "github.com/mistermx/crossplanereleaser/config/v1"
	"github.com/mistermx/crossplanereleaser/internal/build"
	"github.com/mistermx/crossplanereleaser/internal/git"
)

type buildCmd struct {
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
	for _, pkgCfg := range cfg.XPackages {
		buildCfg := &build.PackageBuildConfig{
			PackageDir:  pkgCfg.Dir,
			ExamplesDir: pkgCfg.Examples,
			OutputPath:  getPackageOutputPath(cfg, &pkgCfg),
		}
		err := c.builder.BuildPackage(ctx, buildCfg)
		if err != nil {
			return errors.Wrapf(err, "cannot build package %q", pkgCfg.ID)
		}
	}
	return nil
}

func getPackageOutputPath(cfg *v1.Config, pkgCfg *v1.XPackageConfig) string {
	return filepath.Join(cfg.Dist, pkgCfg.ID, pkgCfg.NameTemplate)
}
