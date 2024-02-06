package build

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"io/fs"

	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/empty"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
	"github.com/spf13/afero"
)

const (
	// MetaFile is the name of a Crossplane package metadata file.
	MetaFile string = "crossplane.yaml"

	// StreamFile is the name of the file in a Crossplane package image that
	// contains its YAML stream.
	StreamFile string = "package.yaml"

	// StreamFileMode determines the permissions on the stream file.
	StreamFileMode fs.FileMode = 0o644

	// XpkgExtension is the extension for compiled Crossplane packages.
	XpkgExtension string = ".xpkg"

	// XpkgMatchPattern is the match pattern for identifying compiled Crossplane packages.
	XpkgMatchPattern string = "*" + XpkgExtension

	// XpkgExamplesFile is the name of the file in a Crossplane package image
	// that contains the examples YAML stream.
	XpkgExamplesFile string = ".up/examples.yaml"

	// AnnotationKey is the key value for xpkg annotations.
	AnnotationKey string = "io.crossplane.xpkg"

	// PackageAnnotation is the annotation value used for the package.yaml
	// layer.
	PackageAnnotation string = "base"

	// ExamplesAnnotation is the annotation value used for the examples.yaml
	// layer.
	// TODO(lsviben) Consider changing this to "examples".
	ExamplesAnnotation string = "upbound"

	// DefaultRegistry is the registry name that will be used when no registry
	// is provided.
	DefaultRegistry string = "xpkg.upbound.io"
)

type ImageBackend struct {
	image  v1.Image
	cfg    *v1.Config
	layers []v1.Layer
}

func NewImageBackend() (*ImageBackend, error) {
	base := empty.Image
	cfgFile, err := base.ConfigFile()
	if err != nil {
		return nil, err
	}
	cfg := &cfgFile.Config
	cfg.Labels = map[string]string{}
	return &ImageBackend{
		image: base,
		cfg:   cfg,
	}, nil
}

// WriteTarball writes the image as a tarball to a file.
func (b *ImageBackend) WriteTarball(fsys afero.Fs, filename string, ref name.Reference) error {
	img, err := mutate.AppendLayers(b.image, b.layers...)
	if err != nil {
		return err
	}
	img, err = mutate.Config(img, *b.cfg)
	if err != nil {
		return err
	}
	// _, err = img.Digest()
	// if err != nil {
	// 	return err
	// }
	file, err := fsys.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	return tarball.Write(ref, img, file)
}

func (b *ImageBackend) WritePackage(content io.Reader) error {
	return b.addLayerSingleFile(PackageAnnotation, StreamFile, content)
}

func (b *ImageBackend) WriteExamples(content io.Reader) error {
	return b.addLayerSingleFile(ExamplesAnnotation, XpkgExamplesFile, content)
}

func (b *ImageBackend) addLayerSingleFile(layerAnnotation string, filename string, content io.Reader) error {
	layer, err := newLayerSingleFile(layerAnnotation, filename, content, StreamFileMode, b.cfg)
	if err != nil {
		return err
	}
	b.layers = append(b.layers, layer)
	// b.image, err = mutate.AppendLayers(b.image, layer)
	return err
}

func newLayerSingleFile(layerAnnotation string, filename string, content io.Reader, mode fs.FileMode, cfg *v1.Config) (v1.Layer, error) {
	tarBuf := &bytes.Buffer{}
	tarw := tar.NewWriter(tarBuf)
	contentBytes, err := io.ReadAll(content)
	if err != nil {
		return nil, err
	}
	h := tar.Header{
		Name: filename,
		Mode: int64(mode),
		Size: int64(len(contentBytes)),
	}
	if err := tarw.WriteHeader(&h); err != nil {
		return nil, err
	}
	if _, err := tarw.Write(contentBytes); err != nil {
		return nil, err
	}
	if err := tarw.Close(); err != nil {
		return nil, err
	}
	layer, err := tarball.LayerFromOpener(func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewReader(tarBuf.Bytes())), nil
	})
	if err != nil {
		return nil, err
	}
	d, err := layer.Digest()
	if err != nil {
		return nil, err
	}
	if layerAnnotation != "" {
		cfg.Labels[xpkgLabel(d.String())] = layerAnnotation
	}
	return layer, nil
}

// Label constructs a specially formated label using the annotationKey.
func xpkgLabel(annotation string) string {
	return fmt.Sprintf("%s:%s", AnnotationKey, annotation)
}
