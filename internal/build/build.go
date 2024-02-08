package build

import (
	"context"
	"path/filepath"

	v1 "github.com/mistermx/crossplanereleaser/config/v1"
)

type PackageBuildConfig struct {
	PackageDir  string
	ExamplesDir string
	OutputPath  string
}

type BuilderBackend interface {
	BuildPackage(ctx context.Context, cfg *PackageBuildConfig) error
}

type Builder struct {
	backend BuilderBackend
}

func NewBuilder(backend BuilderBackend) *Builder {
	return &Builder{
		backend: backend,
	}
}

func (b *Builder) BuildPackagesForConfig(ctx context.Context, cfg *v1.Config) error {
	for _, pkgCfg := range cfg.XPackages {
		buildCfg := &PackageBuildConfig{
			PackageDir:  pkgCfg.Dir,
			ExamplesDir: pkgCfg.Examples,
			OutputPath:  GetPackageOutputPath(cfg, &pkgCfg),
		}
		b.backend.BuildPackage(ctx, buildCfg)
	}
	return nil
}

func GetPackageOutputPath(cfg *v1.Config, pkgCfg *v1.XPackageConfig) string {
	return filepath.Join(cfg.Dist, pkgCfg.ID, pkgCfg.NameTemplate)
}
