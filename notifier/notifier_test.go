package notifier

import (
	"github.com/tompscanlan/q3errord"
	"testing"
)

var notifyTests = []struct {
	in    q3errord.ServiceError
	valid bool
}{
	{q3errord.ServiceError{Service: "testing", Message: "some messages 1"},
		true},

	{q3errord.ServiceError{Service: "testing", Message: "some messages 2"},
		true},

	{q3errord.ServiceError{Service: "testing", Message: "some messages 3"},
		true},
}

func TestNotifierSend(t *testing.T) {
	Notify := LogNotifier{}
	for _, nt := range notifyTests {
		se := nt.in
		Notify.Send(se)
	}
}
