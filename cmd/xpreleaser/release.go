package main

import (
	"github.com/pkg/errors"
	"github.com/spf13/afero"

	v1 "github.com/mistermx/xpreleaser/config/v1"
	"github.com/mistermx/xpreleaser/internal/config"
	"github.com/mistermx/xpreleaser/internal/git"
	"github.com/mistermx/xpreleaser/internal/xpkg/release"
)

type releaseCmd struct{}

func (c *releaseCmd) Run(fsys afero.Fs) error {
	cfgFileName, err := config.FindConfigFile(fsys)
	if err != nil {
		return errors.Wrap(err, "cannot find config file")
	}
	cfg, err := config.Parse(fsys, cfgFileName)
	if err != nil {
		return errors.Wrap(err, "cannot parse config file")
	}
	if err := buildPackages(fsys, cfg); err != nil {
		return err
	}
	return c.releasePackages(fsys, cfg)
}

func (c *releaseCmd) releasePackages(fsys afero.Fs, cfg *v1.Config) error {
	g := git.NewGitCLIBackend()
	err := release.PublishPackages(cfg, g)
	return errors.Wrap(err, "cannot publish packages")
}
