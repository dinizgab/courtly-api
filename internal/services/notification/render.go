package notification

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net/url"
)

//go:embed templates/*.html
var templateFs embed.FS

type Renderer interface {
	Render(naame string, booking any) (string, error)
}

type htmlRendererImpl struct {
	tpl *template.Template
}

func NewHTMLRender(fsys fs.FS) (Renderer, error) {
	if fsys == nil {
		fsys = templateFs
	}

	funcMap := template.FuncMap{
		"urlquery": url.QueryEscape,
	}

    tpls, err := template.New("").
		Funcs(funcMap).
		ParseFS(fsys, "templates/*.html")
	if err != nil {
		return nil, fmt.Errorf("Renderer.Render - failed to read template file: %w", err)
	}

	return &htmlRendererImpl{
		tpl: tpls,
	}, nil
}

func (r *htmlRendererImpl) Render(name string, booking any) (string, error) {
	var buf bytes.Buffer
	if err := r.tpl.ExecuteTemplate(&buf, name, booking); err != nil {
		return "", fmt.Errorf("Renderer.Render - failed to execute template: %w", err)
	}

	return buf.String(), nil
}
