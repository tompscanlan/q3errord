package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tompscanlan/q3errord"
	"io/ioutil"
	"log"
	"net/http"
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

func (sn SlackNotifier) PostSlack(message string) error {

	msg := new(SlackMsg)
	msg.Text = message

	jsonStr, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	log.Println("json: ", string(jsonStr[:]))
	req, err := http.NewRequest("POST", sn.WebhookUrl, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
	return nil
}
