package config

import (
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"sigs.k8s.io/yaml"

	v1 "github.com/mistermx/crossplanereleaser/config/v1"
)

const (
	defaultConfigFileName = ".crossplanereleaser.yaml"
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
	if err := yaml.UnmarshalStrict(raw, cfg, yaml.DisallowUnknownFields); err != nil {
		return nil, err
	}
	err = fillDefaults(filename, cfg)
	return cfg, err
}

func fillDefaults(filename string, cfg *v1.Config) error {
	cfg.ProjectName = valueOrFallback(cfg.ProjectName, filepath.Base(filepath.Dir(filename)))
	cfg.Dist = valueOrFallback(cfg.Dist, "dist")

	for i := range cfg.Builds {
		cfg.Builds[i].Examples = valueOrFallback(cfg.Builds[i].Examples, "examples")

		if cfg.Builds[i].ID == "" {
			if len(cfg.Builds) > 1 {
				return errors.New("build ID is required if there is more than one build")
			}
			// If there is only one build use the project name as ID
			cfg.Builds[i].ID = cfg.ProjectName
		}

		cfg.Builds[i].NameTemplate = valueOrFallback(cfg.Builds[i].NameTemplate, cfg.Builds[i].ID+".xpkg")
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
