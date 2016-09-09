package notifier

import "testing"

var notifyTests = []struct {
	in    ServiceError
	valid bool
}{
	{ServiceError{Service: "testing", Message: "some messages 1"},
		true},

	{ServiceError{Service: "testing", Message: "some messages 2"},
		true},

	{ServiceError{Service: "testing", Message: "some messages 3"},
		true},
}

func TestNotifierSend(t *testing.T) {
	Notify := LogNotifier{}
	for _, nt := range notifyTests {
		se := nt.in
		Notify.Send(se)
	}
}
