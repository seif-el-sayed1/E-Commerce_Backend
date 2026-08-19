// internal/email/sender.go
package email

import (
	"log"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"

	"seif-el-sayed1/E-Commerce_Backend.git/internal/config"
	"seif-el-sayed1/E-Commerce_Backend.git/internal/constants"
)

type SendOptions struct {
	Email   string
	Subject string
	Message string
	HTML    string
}

func Send(options SendOptions) error {
	from := mail.NewEmail("", config.Env.EmailUser)
	to := mail.NewEmail("", options.Email)

	message := mail.NewSingleEmail(from, options.Subject, to, options.Message, options.HTML)
	client := sendgrid.NewSendClient(config.Env.EmailPass)

	response, err := client.Send(message)
	if err != nil {
		log.Println(constants.Error("Email send error: "), err)
		return err
	}

	if response.StatusCode >= 400 {
		log.Println(constants.Error("Email send failed with status: "), response.StatusCode)
		return err
	}

	log.Println(constants.Success("Email sent to " + options.Email))
	return nil
}
