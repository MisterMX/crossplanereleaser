package build

import (
	"context"
)

type PackageBuildConfig struct {
	PackageDir      string
	ExamplesDir     string
	OutputPath      string
	RuntimeImageTar string
}

type BuilderBackend interface {
	BuildPackage(ctx context.Context, cfg *PackageBuildConfig) error
}
