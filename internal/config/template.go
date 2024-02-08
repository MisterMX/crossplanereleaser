package config

import (
	"bytes"
	"text/template"

	v1 "github.com/mistermx/crossplanereleaser/config/v1"
)

func RenderConfigTemplates(cfg *v1.Config, props *ProjectProperties) error {
	for di, docker := range cfg.Dockers {
		for ii, imgTmpl := range docker.ImageTemplates {
			var err error
			cfg.Dockers[di].ImageTemplates[ii], err = renderTemplate(imgTmpl, props)
			if err != nil {
				return err
			}
		}
	}
	for pi, pkgCfg := range cfg.XPackages {
		var err error
		cfg.XPackages[pi].NameTemplate, err = renderTemplate(pkgCfg.NameTemplate, props)
		if err != nil {
			return err
		}
	}
	return nil
}

func renderTemplate(tmplStr string, props *ProjectProperties) (string, error) {
	tmpl, err := template.New("").Parse(tmplStr)
	if err != nil {
		return "", err
	}
	buf := &bytes.Buffer{}
	if err := tmpl.Execute(buf, props); err != nil {
		return "", err
	}
	return buf.String(), nil
}
