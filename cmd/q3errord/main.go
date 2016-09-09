package main

import (
	"log"
	"net/http"
	"net/smtp"
	"sync"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/tompscanlan/q3errord"

	"github.com/tompscanlan/q3errord/notifier"

	"gopkg.in/alecthomas/kingpin.v2"
)

var verbose = kingpin.Flag("verbose", "Verbose mode.").Short('v').Bool()
var port = kingpin.Flag("port", "port to listen on").Default(listenPortDefault).OverrideDefaultFromEnvar("PORT").Short('l').Int()
var username = kingpin.Flag("username", "username for SMTP auth").Short('u').String()
var password = kingpin.Flag("password", "password for SMTP auth").Short('p').String()
var slackurl = kingpin.Flag("slack-webhook", "url for a slack input webhook").String()
var smtpserver = kingpin.Flag("smtp-server", "smtp server we can send mail through").Short('s').String()
var smtpport = kingpin.Flag("smtp-port", "smtp server we can send mail through").Default("587").Int()
var from = kingpin.Flag("from", "email to set as sender of notifications").Default("tscanlan@vmware.com").String()
var to = kingpin.Flag("to", "list of email addresses to send notifications to").Default("tscanlan@vmware.com", "tompscanlan@gmail.com").Strings()

// ServiceErrors is a list of all service errors
var ServiceErrors = []*q3errord.ServiceError{}
var lock = sync.RWMutex{}

// Notify is set to the type of notifier we want to use
var Notify q3errord.Notifier

const (
	listenPortDefault = "8080"
)

func init() {
	setupFlags()
	q3errord.Verbose = *verbose
}

func setupFlags() {
	kingpin.CommandLine.HelpFlag.Short('h')
	kingpin.Parse()
}
func main() {

	if *username != "" && *password != "" && *smtpserver != "" {
		log.Println("will use SMTP notifications")
		auth := smtp.PlainAuth(
			"",
			*username,
			*password,
			*smtpserver,
		)
		Notify = notifier.MailNotifier{From: *from, To: *to, Auth: auth}
	} else if *slackurl != "" {

		log.Println("will use slack notifications")
		Notify = notifier.SlackNotifier{WebhookUrl: *slackurl}
	} else {
		log.Println("will use log notifications")
		Notify = notifier.LogNotifier{}
	}

	api := rest.NewApi()

	statusMw := &rest.StatusMiddleware{}
	api.Use(statusMw)
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Get("/.status", func(w rest.ResponseWriter, r *rest.Request) {
			w.WriteJson(statusMw.GetStatus())
		}),

		rest.Get("/errors", GetAllErrors),
		rest.Post("/error", PostError),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	log.Fatal(http.ListenAndServe(":8080", api.MakeHandler()))
}

func GetAllErrors(w rest.ResponseWriter, r *rest.Request) {
	lock.RLock()
	w.WriteJson(&ServiceErrors)
	lock.RUnlock()
}

func PostError(w rest.ResponseWriter, r *rest.Request) {
	se := q3errord.ServiceError{}
	err := r.DecodeJsonPayload(&se)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if se.Service == "" {
		rest.Error(w, "Service name required", 400)
		return
	}
	if se.Message == "" {
		rest.Error(w, "Service error message required", 400)
		return
	}
	lock.Lock()
	ServiceErrors = append(ServiceErrors, &se)
	lock.Unlock()
	w.WriteJson(&se)

	go Notify.Send(se)
}
