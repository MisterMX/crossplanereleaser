package build

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/pkg/errors"
)

type CrankCLIBackend struct {
	cmd string
}

func NewCrankCLIBackend() *CrankCLIBackend {
	return &CrankCLIBackend{
		cmd: "crank",
	}
}

func (c CrankCLIBackend) exec(ctx context.Context, args ...string) error {
	cmd := exec.CommandContext(ctx, c.cmd, args...)
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return errors.Wrap(err, c.cmd)
}

func (c *CrankCLIBackend) BuildPackage(ctx context.Context, cfg *PackageBuildConfig) error {
	args := []string{
		"xpkg", "build",
		fmt.Sprintf("--package-root=%s", cfg.PackageDir),
		fmt.Sprintf("--package-file=%s", cfg.OutputPath),
	}
	if cfg.ExamplesDir != "" {
		args = append(args, fmt.Sprintf("--examples-root=%s", cfg.ExamplesDir))
	}
	// TODO: Support setting --controller-tar option
	return c.exec(ctx, args...)
}
