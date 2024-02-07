package release

import (
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
	"github.com/pkg/errors"

	configv1 "github.com/mistermx/crossplanereleaser/config/v1"
	"github.com/mistermx/crossplanereleaser/internal/xpkg/build"
)

type ImagePublisherBackend interface {
	Put(ref name.Reference, img v1.Image) error
}

func PublishPackages(cfg *configv1.Config) error {
	for _, docker := range cfg.Dockers {
		publishForDocker(cfg, &docker)
	}
	return nil
}

func publishForDocker(cfg *configv1.Config, docker *configv1.DockerConfig) error {
	keychain, err := BuildKeyChainFromConfig(docker.Logins)
	if err != nil {
		return errors.Wrap(err, "cannot setup docker keychain")
	}
	backend := NewDockerRegistryBackend(keychain)
	pkgCfgs, err := selectXPackageConfigsByIDs(cfg, docker.IDs)
	if err != nil {
		return err
	}
	imgRefs, err := imageNamesAsRefs(docker.ImageTemplates)
	if err != nil {
		return err
	}
	for _, pkgCfg := range pkgCfgs {
		tarballPath := build.GetXPackageOutputPath(cfg, &pkgCfg)
		img, err := tarball.ImageFromPath(tarballPath, nil)
		if err != nil {
			return errors.Wrap(err, "cannot load image from tarball")
		}
		for _, ref := range imgRefs {
			if err := backend.Put(ref, img); err != nil {
				return errors.Wrapf(err, "cannot publish image %q", ref.String())
			}
		}
	}
	return nil
}

func selectXPackageConfigsByIDs(cfg *configv1.Config, pkgIDs []string) ([]configv1.XPackageConfig, error) {
	if len(pkgIDs) == 0 {
		return cfg.XPackages, nil
	}
	selected := make([]configv1.XPackageConfig, 0, len(pkgIDs))
	for _, id := range pkgIDs {
		found := false
		for _, pkgCfg := range cfg.XPackages {
			if pkgCfg.ID == id {
				selected = append(selected, pkgCfg)
				found = true
				break
			}
		}
		if !found {
			return nil, errors.Errorf("no xpackage with ID %q", id)
		}
	}
	return selected, nil
}

func imageNamesAsRefs(imageTemplates []string) ([]name.Reference, error) {
	refs := make([]name.Reference, len(imageTemplates))
	for i, refStr := range imageTemplates {
		ref, err := name.ParseReference(refStr)
		if err != nil {
			return nil, errors.Wrapf(err, "%d", i)
		}
		refs[i] = ref
	}
	return refs, nil
}
