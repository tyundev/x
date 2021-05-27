package mail_gunc

import (
	"context"
	//"fmt"
	"time"

	mailgun "github.com/mailgun/mailgun-go"
)

func SendMailGunc(subject, body, mailRecive string) error {
	var yourDomain = "sandbox3491fb598a84428a9c6a3a5d0ca18fc5.mailgun.org" // e.g. mg.yourcompany.com

	// You can find the Private API Key in your Account Menu, under "Settings":
	// (https://app.mailgun.com/app/account/security)
	var privateAPIKey = "c0980275d083c30b1d92d55b22dc5c2a-65b08458-38faa1bc"
	//var publicKey = "pubkey-de1e9e9cd301b73601c727ef1d5a348d"
	// Create an instance of the Mailgun Client
	mg := mailgun.NewMailgun(yourDomain, privateAPIKey)

	sender := "longtm@miraway.vn"

	// The message object allows you to add attachments and Bcc recipients
	message := mg.NewMessage(sender, subject, body, "longtm23@gmail.com")
	//message.SetTemplate("passwordReset")
	//message.AddTemplateVariable("passwordResetLink", "some link to your site unique to your user")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	// Send the message	with a 10 second timeout
	_, _, err := mg.Send(ctx, message)
	return err
}
