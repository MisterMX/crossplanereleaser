package main

import (
	"context"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/afero"

	v1 "github.com/mistermx/crossplanereleaser/config/v1"
	"github.com/mistermx/crossplanereleaser/internal/git"
	"github.com/mistermx/crossplanereleaser/internal/xpkg/build"
	"github.com/mistermx/crossplanereleaser/internal/xpkg/parse"
)

type buildCmd struct{}

func (c *buildCmd) Run(fsys afero.Fs) error {
	g := git.NewGitCLIBackend()
	cfg, err := getConfig(fsys, g)
	if err != nil {
		return err
	}
	return buildPackages(fsys, cfg)
}

func buildPackages(fsys afero.Fs, cfg *v1.Config) error {
	for _, pkgCfg := range cfg.XPackages {
		err := buildPackage(fsys, cfg, &pkgCfg)
		if err != nil {
			return errors.Wrapf(err, "cannot build package %q", pkgCfg.ID)
		}
	}
	return nil
}

func buildPackage(fsys afero.Fs, cfg *v1.Config, pkgCfg *v1.XPackageConfig) error {
	ctx := context.TODO()
	parseBackend := parse.NewFSDirBackend(fsys, pkgCfg.Dir, pkgCfg.Examples)
	pkg, err := parse.Parse(ctx, parseBackend)
	if err != nil {
		return errors.Wrap(err, "cannot parse package")
	}
	buildBackend, err := build.NewImageBackend()
	if err != nil {
		return errors.Wrap(err, "cannot setup builder backend")
	}
	if err := build.BuildImage(ctx, buildBackend, pkg); err != nil {
		return errors.Wrap(err, "cannot build image")
	}
	outputPath := build.GetXPackageOutputPath(cfg, pkgCfg)
	if err != nil {
		return err
	}
	if err := fsys.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return err
	}
	if err := buildBackend.WriteTarball(fsys, outputPath, nil); err != nil {
		return errors.Wrap(err, "cannot write image tarball")
	}
	return err
}
