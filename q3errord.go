package q3errord

// Verbose lets other packages know we want to be noisey

var Verbose = true

type Notifier interface {
	Send(ServiceError)
}
type ServiceError struct {
	Service string `json:"service,omitempty"`
	Message string `json:"message,omitempty"`
}
