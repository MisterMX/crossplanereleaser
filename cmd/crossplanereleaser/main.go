package main

import (
	"github.com/alecthomas/kong"
	"github.com/spf13/afero"
)

var cli struct {
	Build   buildCmd   `cmd:"build" help:"Build artifacts"`
	Release releaseCmd `cmd:"release" help:"Release and publish artifacts"`

	Version versionCmd `cmd:"version" help:"Print version information"`
}

var _ = kong.Must(&cli)

func main() {
	fs := afero.NewOsFs()

	ctx := kong.Parse(&cli,
		kong.Name("crossplanereleaser"),
		kong.Description("CLI utility to deal with certain tasks around CNP@DBNetz."),
		kong.BindTo(fs, (*afero.Fs)(nil)),
	)

	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
