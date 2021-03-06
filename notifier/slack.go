package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/tompscanlan/q3errord"
)

type SlackMsg struct {
	Text string `json:"text,omitempty"`
}

// MailNotifier allows notification via smtp
type SlackNotifier struct {
	WebhookUrl string
	Channel    string
}

func (sn SlackNotifier) Send(se q3errord.ServiceError) {

	if q3errord.Verbose {
		log.Printf("Notify via slack: %s, (%s)", se.Service, se.Message)
	}
	sn.PostSlack(fmt.Sprintf("%s: (%s)", se.Service, se.Message))
}

func (sn SlackNotifier) PostSlack(message string) {

	msg := new(SlackMsg)
	msg.Text = message

	jsonStr, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("json: ", string(jsonStr[:]))
	req, err := http.NewRequest("POST", sn.WebhookUrl, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	log.Println("response Status:", resp.Status)
	log.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println("response Body:", string(body))
}
