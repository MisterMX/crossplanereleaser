package git

import (
	"bytes"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// GitCLIBackend uses the system's "git" command to implement [GitBackend].
type GitCLIBackend struct{}

func NewGitCLIBackend() *GitCLIBackend {
	return &GitCLIBackend{}
}

func (g *GitCLIBackend) run(args ...string) (string, error) {
	stdout := bytes.Buffer{}
	cmd := exec.Command("git", args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = &stdout
	err := cmd.Run()
	return stdout.String(), errors.Wrapf(err, "cmd %q", cmd.String())
}

func (g *GitCLIBackend) hasTags() (bool, error) {
	out, err := g.run("tag")
	return out != "", err
}

func (g *GitCLIBackend) GetCommitCount(ref Ref) (int, error) {
	out, err := g.run("rev-list", "--count", string(ref))
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(strings.TrimSpace(out))
}

func (g *GitCLIBackend) GetTag(ref Ref) (string, error) {
	hasTag, err := g.hasTags()
	if err != nil {
		return "", err
	}
	if !hasTag {
		return "", errNoTag
	}
	tag, err := g.run("describe", "--tags", "--abbrev=8", string(ref))
	return strings.TrimSpace(tag), err
}

func (g *GitCLIBackend) GetCommitSHA(ref Ref) (string, error) {
	sha, err := g.run("rev-parse", string(ref))
	return strings.TrimSpace(sha), err
}

func (g *GitCLIBackend) GetCurrentBranch() (string, error) {
	out, err := g.run("rev-parse", "--abbrev-ref", string(RefHead))
	if err != nil {
		return "", err
	}
	ref := strings.TrimSpace(out)

	// rev-parse returns HEAD if we are not on a branch.
	if ref != string(RefHead) {
		return ref, nil
	}

	branches, err := g.run("branch", "-a", "--points-at", string(RefHead))
	if err != nil {
		return "", err
	}
	branchesSplit := strings.Split(branches, "\n")
	if len(branchesSplit) == 0 {
		return "", errNoBranch
	}
	return branchesSplit[0], nil
}

func (g *GitCLIBackend) IsGitTreeClean() (bool, error) {
	out, err := g.run("status", "--porcelain")
	if err != nil {
		return false, err
	}
	return out == "", nil
}
