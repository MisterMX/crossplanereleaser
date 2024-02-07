package release

import (
	"bytes"
	"text/template"

	"github.com/mistermx/xpreleaser/internal/git"
)

type templateData struct {
	Branch      string
	Version     string
	FullCommit  string
	ShortCommit string
}

func buildTemplateData(g git.Backend, ref git.Ref) (*templateData, error) {
	branch, err := g.GetCurrentBranch()
	if err != nil {
		return nil, err
	}
	version, err := git.GetVersion(g, ref)
	if err != nil {
		return nil, err
	}
	commitSha, err := g.GetCommitSHA(ref)
	if err != nil {
		return nil, err
	}
	return &templateData{
		Branch:      branch,
		Version:     version,
		FullCommit:  commitSha,
		ShortCommit: commitSha[:8],
	}, nil
}

func renderImageTemplate(imgTemplate string, data *templateData) (string, error) {
	tmpl, err := template.New("image_template").Parse(imgTemplate)
	if err != nil {
		return "", err
	}
	buf := bytes.Buffer{}
	err = tmpl.Execute(&buf, data)
	return buf.String(), err
}
