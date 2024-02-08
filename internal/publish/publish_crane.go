package publish

import (
	"context"
	"os"
	"os/exec"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/pkg/errors"
)

type CraneCLIPublisher struct {
	cmd string
}

func NewCraneCLIPublisher() *CraneCLIPublisher {
	return &CraneCLIPublisher{
		cmd: "crane",
	}
}

func (c *CraneCLIPublisher) exec(ctx context.Context, args ...string) error {
	cmd := exec.CommandContext(ctx, c.cmd, args...)
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return errors.Wrap(err, c.cmd)
}

func (c *CraneCLIPublisher) PublishPackage(ctx context.Context, filename string, ref name.Reference) error {
	return c.exec(ctx, "push", filename, ref.String())
}
