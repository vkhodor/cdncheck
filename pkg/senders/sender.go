package senders

type Sender interface {
	Send(string) error
}
