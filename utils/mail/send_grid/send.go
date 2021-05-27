package send_grid

import (
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func TestGrid(subject, body, mailRecive string) error {
	from := mail.NewEmail("TEST", "longtm@miraway.vn")
	to := mail.NewEmail("recive", mailRecive)
	plainTextContent := "test grid"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, body)
	client := sendgrid.NewSendClient("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")
	_, err := client.Send(message)
	return err
}
