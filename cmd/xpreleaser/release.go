package main

import (
	"github.com/pkg/errors"
	"github.com/spf13/afero"

	v1 "github.com/mistermx/xpreleaser/config/v1"
	"github.com/mistermx/xpreleaser/internal/git"
	"github.com/mistermx/xpreleaser/internal/xpkg/release"
)

type releaseCmd struct {
	SkipBuild bool `help:"Don't execute build step before release."`
}

func (c *releaseCmd) Run(fsys afero.Fs) error {
	g := git.NewGitCLIBackend()
	cfg, err := getConfig(fsys, g)
	if err != nil {
		return err
	}

	if !c.SkipBuild {
		if err := buildPackages(fsys, cfg); err != nil {
			return err
		}
	}

	return c.releasePackages(fsys, cfg)
}

func (c *releaseCmd) releasePackages(fsys afero.Fs, cfg *v1.Config) error {
	err := release.PublishPackages(cfg)
	return errors.Wrap(err, "cannot publish packages")
}
