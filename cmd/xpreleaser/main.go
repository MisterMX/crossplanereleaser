package main

import (
	"os"

	"github.com/alecthomas/kong"
	"github.com/go-log/log"
	fmtLog "github.com/go-log/log/fmt"
	"github.com/spf13/afero"
)

var cli struct {
	Build buildCmd `cmd:"build" help:"Build artifacts"`

	Version versionCmd `cmd:"version" help:"Print version information"`
}

var _ = kong.Must(&cli)

func main() {
	fs := afero.NewOsFs()
	logger := fmtLog.NewFromWriter(os.Stderr)

	ctx := kong.Parse(&cli,
		kong.Name("cnpctl"),
		kong.Description("CLI utility to deal with certain tasks around CNP@DBNetz."),
		kong.BindTo(fs, (*afero.Fs)(nil)),
		kong.BindTo(logger, (*log.Logger)(nil)),
	)

	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
