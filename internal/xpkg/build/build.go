package build

import (
	"bytes"
	"context"
	"io"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/yaml"

	"github.com/mistermx/crossplanereleaser/internal/xpkg"
)

type BuilderBackend interface {
	WritePackage(in io.Reader) error
	WriteExamples(in io.Reader) error
}

func BuildImage(ctx context.Context, backend BuilderBackend, pkg *xpkg.Package) error {
	packageBuf := &bytes.Buffer{}
	examplesBuf := &bytes.Buffer{}
	if err := writeObjectListYaml(pkg.MetaObjects, packageBuf); err != nil {
		return errors.Wrap(err, "cannot write meta objects to YAML")
	}
	if err := writeObjectListYaml(pkg.PackageObjects, packageBuf); err != nil {
		return errors.Wrap(err, "cannot write package objects to YAML")
	}
	if err := writeObjectListYaml(pkg.ExampleObjects, examplesBuf); err != nil {
		return errors.Wrap(err, "cannot write package objects to YAML")
	}
	if err := backend.WritePackage(packageBuf); err != nil {
		return errors.Wrap(err, "cannot write package")
	}
	// Write examples if there are any
	if examplesBuf.Len() > 0 {
		if err := backend.WriteExamples(examplesBuf); err != nil {
			return errors.Wrap(err, "cannot write example")
		}
	}
	return nil
}

func writeObjectListYaml(list []*unstructured.Unstructured, out io.Writer) error {
	for _, u := range list {
		raw, err := yaml.Marshal(u)
		if err != nil {
			return err
		}
		io.WriteString(out, "---\n")
		out.Write(raw)
		io.WriteString(out, "\n")
	}
	return nil
}
