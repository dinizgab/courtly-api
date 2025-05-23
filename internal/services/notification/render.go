package notification

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"io/fs"

	"github.com/dinizgab/booking-mvp/internal/entity"
)

//go:embed templates/booking_confirmation.html
var templateFs embed.FS

type Renderer interface {
	Render(booking entity.BookingConfirmationDTO) (string, error)
}

type htmlRendererImpl struct {
	tpl *template.Template
}

func NewHTMLRender(fsys fs.FS) (Renderer, error) {
    if fsys == nil {
        fsys = templateFs
    }

	tpl, err := template.ParseFS(fsys, "templates/booking_confirmation.html")
	if err != nil {
		return nil, fmt.Errorf("Renderer.Render - failed to read template file: %w", err)
	}

	return &htmlRendererImpl{
        tpl: tpl,
    }, nil
}

func (r *htmlRendererImpl) Render(booking entity.BookingConfirmationDTO) (string, error) {
    fmt.Println("Rendering email template with booking details:", booking)
	var buf bytes.Buffer
	if err := r.tpl.Execute(&buf, booking); err != nil {
		return "", fmt.Errorf("Renderer.Render - failed to execute template: %w", err)
	}

	return buf.String(), nil
}
