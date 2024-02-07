package config

import (
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"sigs.k8s.io/yaml"

	v1 "github.com/mistermx/crossplanereleaser/config/v1"
)

const (
	defaultConfigFileName = ".xpreleaser.yaml"
)

func FindConfigFile(_ afero.Fs) (string, error) {
	// TODO: Walk up the file tree and look in parent directories?
	return defaultConfigFileName, nil
}

func Parse(fsys afero.Fs, filename string) (*v1.Config, error) {
	if !filepath.IsAbs(filename) {
		var err error
		filename, err = filepath.Abs(filename)
		if err != nil {
			return nil, err
		}
	}

	raw, err := afero.ReadFile(fsys, filename)
	if err != nil {
		return nil, err
	}
	cfg := &v1.Config{}
	if err := yaml.Unmarshal(raw, cfg); err != nil {
		return nil, err
	}
	err = fillDefaults(filename, cfg)
	return cfg, err
}

func fillDefaults(filename string, cfg *v1.Config) error {
	cfg.ProjectName = valueOrFallback(cfg.ProjectName, filepath.Base(filepath.Dir(filename)))
	cfg.Dist = valueOrFallback(cfg.Dist, "dist")

	for i := range cfg.XPackages {
		cfg.XPackages[i].Examples = valueOrFallback(cfg.XPackages[i].Examples, "examples")

		if cfg.XPackages[i].ID == "" {
			if len(cfg.XPackages) > 1 {
				return errors.New("package ID is required if there is more than one package")
			}
			// If there is only one package use the project name as ID
			cfg.XPackages[i].ID = cfg.ProjectName
		}
	}
	return nil
}

func valueOrFallback[T comparable](val, fallback T) T {
	var zero T
	if val == zero {
		return fallback
	}
	return val
}
