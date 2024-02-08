package config

import (
	"os"
	"strings"

	v1 "github.com/mistermx/crossplanereleaser/config/v1"
	"github.com/mistermx/crossplanereleaser/internal/git"
)

type ProjectProperties struct {
	// Common
	ProjectName string
	Env         map[string]string

	// Git
	Branch       string
	Tag          string
	FullCommit   string
	ShortCommit  string
	IsGitDirty   bool
	IsGitClean   bool
	GitTreeState string
}

func BuildProjectProperties(g git.Client, cfg *v1.Config) (*ProjectProperties, error) {
	props := &ProjectProperties{
		ProjectName: cfg.ProjectName,
		Env:         getEnvMap(),
	}
	ref := git.RefHead

	var err error
	props.Branch, err = g.GetCurrentBranch()
	if err != nil && !git.IsErrNoBranch(err) {
		return nil, err
	}

	props.Tag, err = git.GetVersion(g, ref)
	if err != nil {
		return nil, err
	}

	props.FullCommit, err = g.GetCommitSHA(ref)
	if err != nil {
		return nil, err
	}
	props.ShortCommit = props.FullCommit[:8]

	props.IsGitClean, err = g.IsGitTreeClean()
	if err != nil {
		return nil, err
	}
	props.IsGitDirty = !props.IsGitClean
	if props.IsGitClean {
		props.GitTreeState = "clean"
	} else {
		props.GitTreeState = "dirty"
	}

	return props, nil
}

func getEnvMap() map[string]string {
	environ := os.Environ()
	m := make(map[string]string, len(environ))
	for _, s := range os.Environ() {
		split := strings.SplitN(s, "=", 1)
		m[split[0]] = split[1]
	}
	return m
}
