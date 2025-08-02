package domain

type Sender interface {
	Send(task Task) error
}
