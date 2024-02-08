package publish

import (
	"context"

	"github.com/google/go-containerregistry/pkg/name"
)

type PackagePublisher interface {
	PublishPackage(ctx context.Context, filename string, ref name.Reference) error
}
