package send_grid

import (
	"fmt"
	"log"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

const (
	KeySendGrid = ""
)

func TestGrid(subject, body, mailRecive string) error {
	from := mail.NewEmail("TEST", "")
	to := mail.NewEmail("recive", mailRecive)
	plainTextContent := "test grid"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, body)
	client := sendgrid.NewSendClient(KeySendGrid)
	_, err := client.Send(message)
	return err
}

func SendMany(mails []string) {
	m := mail.NewV3Mail()

	from := mail.NewEmail("test", "")
	content := mail.NewContent("text/html", "<p> %fname% : %CustomerID% - Personalizations are awesome!</p>")

	m.SetFrom(from)
	m.AddContent(content)

	// create new *Personalization
	personalization := mail.NewPersonalization()
	var emails = make([]*mail.Email, 0)
	for _, val := range mails {
		emails = append(emails, mail.NewEmail("recive", val))
	}

	personalization.AddTos(emails...)
	personalization.SetSubstitution("%fname%", "recipient")
	personalization.SetSubstitution("%CustomerID%", "CUSTOMER ID GOES HERE")
	personalization.Subject = "Having fun learning about personalizations?"

	// add `personalization` to `m`
	m.AddPersonalizations(personalization)

	request := sendgrid.GetRequest(KeySendGrid, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = mail.GetRequestBody(m)
	response, err := sendgrid.API(request)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
}
