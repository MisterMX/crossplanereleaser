package config

import (
	"bytes"
	"text/template"

	v1 "github.com/mistermx/crossplanereleaser/config/v1"
)

func RenderConfigTemplates(cfg *v1.Config, props *ProjectProperties) error {
	for pi, push := range cfg.Pushes {
		for ii, imgTmpl := range push.ImageTemplates {
			var err error
			cfg.Pushes[pi].ImageTemplates[ii], err = renderTemplate(imgTmpl, props)
			if err != nil {
				return err
			}
		}
	}
	for pi, build := range cfg.Builds {
		var err error
		cfg.Builds[pi].NameTemplate, err = renderTemplate(build.NameTemplate, props)
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
