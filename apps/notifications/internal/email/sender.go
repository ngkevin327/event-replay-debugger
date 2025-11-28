package email

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"html/template"
	"io"
)

//go:embed templates/*.html
var templateFS embed.FS

// SESClient sends rendered email.
type SESClient interface {
	SendEmail(ctx context.Context, to, subject, htmlBody string) error
}

// Sender renders templates and delivers via SES.
type Sender struct {
	Client SESClient
	From   string
}

// SendTemplate renders name with data and sends to recipient.
func (s *Sender) SendTemplate(ctx context.Context, to, name string, data any) error {
	if s.Client == nil {
		return fmt.Errorf("email client not configured")
	}
	tpl, err := template.ParseFS(templateFS, "templates/"+name+".html")
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return err
	}
	subject := templateSubject(name)
	return s.Client.SendEmail(ctx, to, subject, buf.String())
}

func templateSubject(name string) string {
	switch name {
	case "incident_ready":
		return "Incident ready for review"
	case "replay_completed":
		return "Replay completed"
	default:
		return "Replay notification"
	}
}

// LogSES logs rendered email (staging/dev).
type LogSES struct {
	Out io.Writer
}

func (l *LogSES) SendEmail(_ context.Context, to, subject, htmlBody string) error {
	_, err := fmt.Fprintf(l.Out, "to=%s subject=%s bytes=%d\n", to, subject, len(htmlBody))
	return err
}
