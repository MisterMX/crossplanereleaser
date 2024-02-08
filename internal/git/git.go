package git

import (
	"fmt"

	"github.com/pkg/errors"
)

type Ref string

const (
	RefHead Ref = "HEAD"
)

type Client interface {
	GetTag(ref Ref) (string, error)
	GetCommitSHA(ref Ref) (string, error)
	GetCommitCount(ref Ref) (int, error)
	GetCurrentBranch() (string, error)
	IsGitTreeClean() (bool, error)
}

var (
	errNoTag    = errors.New("errNoTag")
	errNoBranch = errors.New("errNoBranch")
)

func IsErrNoTag(err error) bool {
	return err == errNoTag
}

func IsErrNoBranch(err error) bool {
	return err == errNoBranch
}

func GetVersion(git Client, ref Ref) (string, error) {
	tag, err := git.GetTag(ref)
	if err == nil || !IsErrNoTag(err) {
		return tag, err
	}
	return getUntaggedVersion(git, ref)
}

func getUntaggedVersion(git Client, ref Ref) (string, error) {
	base := "v0.0.0"
	revCount, err := git.GetCommitCount(ref)
	if err != nil {
		return "", err
	}
	sha, err := git.GetCommitSHA(ref)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s-%d-g%s", base, revCount, sha[:8]), nil
}
