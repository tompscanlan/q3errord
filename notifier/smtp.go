package notifier

import (
	"fmt"
	"github.com/tompscanlan/q3errord"
	"log"
	"net/smtp"
	"strings"
)

// MailNotifier allows notification via smtp
type MailNotifier struct {
	From string
	To   []string
	Auth smtp.Auth
}

func (mn MailNotifier) Send(se q3errord.ServiceError) {

	if q3errord.Verbose {
		log.Printf("Notify via mail: %s, (%s)", se.Service, se.Message)
	}
	body := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n\n%s", mn.From, strings.Join(mn.To, ","), se.Service, se.Message)
	// Set up authentication information.

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	err := smtp.SendMail(
		"smtp.sendgrid.net:587",
		mn.Auth,
		mn.From,
		mn.To,
		[]byte(body),
	)
	if err != nil {
		log.Fatal(err)
	}
}
