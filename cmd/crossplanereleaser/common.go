package main

import (
	"github.com/pkg/errors"
	"github.com/spf13/afero"

	v1 "github.com/mistermx/crossplanereleaser/config/v1"
	"github.com/mistermx/crossplanereleaser/internal/config"
	"github.com/mistermx/crossplanereleaser/internal/git"
)

func getConfig(fsys afero.Fs, g git.Backend) (*v1.Config, error) {
	cfgFileName, err := config.FindConfigFile(fsys)
	if err != nil {
		return nil, errors.Wrap(err, "cannot find config file")
	}
	cfg, err := config.Parse(fsys, cfgFileName)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse config file")
	}
	props, err := config.BuildProjectProperties(g, cfg)
	if err != nil {
		return nil, errors.Wrap(err, "cannot setup project properties")
	}

	if err := config.RenderConfigTemplates(cfg, props); err != nil {
		return nil, errors.Wrap(err, "cannot render config template fields")
	}
	return cfg, nil
}
