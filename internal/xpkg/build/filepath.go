package build

import (
	"path/filepath"

	v1 "github.com/mistermx/crossplanereleaser/config/v1"
)

func GetXPackageOutputPath(cfg *v1.Config, pkgCfg *v1.XPackageConfig) string {
	return filepath.Join(cfg.Dist, pkgCfg.ID, "package.xpkg")
}
