package parse

import (
	"context"
	"io"
	"io/fs"
	"strings"

	"github.com/spf13/afero"
)

type FsDirBackend struct {
	fsys       afero.Fs
	packageDir string
	exampleDir string
}

func NewFSDirBackend(fsys afero.Fs, packageDir, exampleDir string) *FsDirBackend {
	return &FsDirBackend{
		fsys:       fsys,
		packageDir: packageDir,
		exampleDir: exampleDir,
	}
}

func (b *FsDirBackend) GetPackageFiles(_ context.Context) (chan io.ReadCloser, error) {
	files := make(chan io.ReadCloser, 8)
	go func() {
		afero.Walk(b.fsys, b.packageDir, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			// Ignore hidden files or directories
			if info.Name()[0] == '.' {
				return fs.SkipDir
			}
			if info.IsDir() {
				return nil
			}
			if !strings.HasSuffix(path, ".yaml") && !strings.HasSuffix(path, ".yml") {
				return nil
			}
			f, err := b.fsys.Open(path)
			if err != nil {
				return err
			}
			files <- f
			return nil
		})
		// TODO: handle errors in a way
		close(files)
	}()
	return files, nil
}

func (b *FsDirBackend) GetExampleFiles(_ context.Context) (chan io.ReadCloser, error) {
	files := make(chan io.ReadCloser, 8)
	if exists, _ := afero.DirExists(b.fsys, b.exampleDir); !exists {
		close(files)
		return files, nil
	}

	go func() {
		afero.Walk(b.fsys, b.exampleDir, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			// Ignore hidden files and directories
			isHidden := info.Name()[0] == '.'
			if info.IsDir() && info.Name()[0] == '.' {
				return fs.SkipDir
			}
			if isHidden {
				return nil
			}
			if !strings.HasSuffix(path, ".yaml") && !strings.HasSuffix(path, ".yml") {
				return nil
			}
			f, err := b.fsys.Open(path)
			if err != nil {
				return err
			}
			files <- f
			return nil
		})
		// TODO: handle errors in a way
		close(files)
	}()
	return files, nil
}
