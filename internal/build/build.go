package build

import (
	"context"
)

type PackageBuildConfig struct {
	PackageDir  string
	ExamplesDir string
	OutputPath  string
}

type BuilderBackend interface {
	BuildPackage(ctx context.Context, cfg *PackageBuildConfig) error
}
