package notifier

import (
	"github.com/tompscanlan/q3errord"
	"log"
)

type LogNotifier struct {
}

func (ln LogNotifier) Send(se q3errord.ServiceError) {
	log.Printf("Subject:%s, Message:%s", se.Service, se.Message)
}
