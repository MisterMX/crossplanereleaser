package release

import (
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

type DockerRegistryBackend struct {
	keychain authn.Keychain
}

func NewDockerRegistryBackend(keychain authn.Keychain) *DockerRegistryBackend {
	return &DockerRegistryBackend{
		keychain: keychain,
	}
}

func (b *DockerRegistryBackend) Put(ref name.Reference, img v1.Image) error {
	return remote.Put(ref, img, remote.WithAuthFromKeychain(b.keychain))
}
