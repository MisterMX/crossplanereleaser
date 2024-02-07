package parse

import (
	"context"
	"io"

	"github.com/mistermx/go-utils/k8s/yaml"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/mistermx/crossplanereleaser/internal/xpkg"
)

type ParserBackend interface {
	// GetPackageFiles returns the individual package YAML file streams from the
	// parser backend.
	// Every file can contain multiple documents that are parsed into individual
	// objects.
	//
	// The function returns the filestreams as [io.ReadCloser].
	// It is in the receivers responsibility to [io.ReadCloser.Close] after
	// reading.
	GetPackageFiles(ctx context.Context) (chan io.ReadCloser, error)

	// GetExampleFiles returns the individual example YAML file streams from the
	// parser backend.
	// Every file can contain multiple documents that are parsed into individual
	// objects.
	//
	// The function returns the filestreams as [io.ReadCloser].
	// It is in the receivers responsibility to [io.ReadCloser.Close] after
	// reading.
	GetExampleFiles(ctx context.Context) (chan io.ReadCloser, error)
}

func Parse(ctx context.Context, backend ParserBackend) (*xpkg.Package, error) {
	packageFiles, err := backend.GetPackageFiles(ctx)
	if err != nil {
		return nil, err
	}
	packageObjects := []*unstructured.Unstructured{}
	metaObjects := []*unstructured.Unstructured{}
	for f := range packageFiles {
		objects, err := parseStream(f)
		if err != nil {
			return nil, err
		}
		for _, o := range objects {
			switch {
			case isMetaObject(o):
				metaObjects = append(metaObjects, o)
			case isPackageObject(o):
				packageObjects = append(packageObjects, o)
			default:
				return nil, errors.Errorf("object %q with kind %q is not a package or meta object", o.GetName(), o.GroupVersionKind().String())
			}
		}
	}
	if len(metaObjects) == 0 {
		return nil, errors.New("package has no meta descriptor")
	}

	exampleFiles, err := backend.GetExampleFiles(ctx)
	if err != nil {
		return nil, err
	}
	exampleObjects := []*unstructured.Unstructured{}
	for f := range exampleFiles {
		objects, err := parseStream(f)
		if err != nil {
			return nil, err
		}
		exampleObjects = append(exampleObjects, objects...)
	}

	return &xpkg.Package{
		PackageObjects: packageObjects,
		MetaObjects:    metaObjects,
		ExampleObjects: exampleObjects,
	}, nil
}

func isMetaObject(o *unstructured.Unstructured) bool {
	return o.GetObjectKind().GroupVersionKind().Group == "meta.pkg.crossplane.io"
}

func isPackageObject(o *unstructured.Unstructured) bool {
	return o.GetObjectKind().GroupVersionKind().Group == "apiextensions.crossplane.io"
}

func parseStream(in io.ReadCloser) ([]*unstructured.Unstructured, error) {
	defer in.Close()
	return yaml.UnmarshalObjectsReader[*unstructured.Unstructured](in)
}
