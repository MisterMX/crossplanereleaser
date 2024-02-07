package config

import (
	"bytes"
	"text/template"

	v1 "github.com/mistermx/xpreleaser/config/v1"
	"github.com/pkg/errors"
)

func RenderConfigTemplates(cfg *v1.Config, props *ProjectProperties) error {
	if err := renderDockerImageTemplates(cfg, props); err != nil {
		return errors.Wrap(err, "cannot render docker image_templates")
	}
	return nil
}

func renderDockerImageTemplates(cfg *v1.Config, props *ProjectProperties) error {
	for di, docker := range cfg.Dockers {
		for ii, imgTmpl := range docker.ImageTemplates {
			tmpl, err := template.New("image_template").Parse(imgTmpl)
			if err != nil {
				return err
			}
			buf := &bytes.Buffer{}
			if err := tmpl.Execute(buf, props); err != nil {
				return err
			}
			cfg.Dockers[di].ImageTemplates[ii] = buf.String()
		}
	}
	return nil
}
