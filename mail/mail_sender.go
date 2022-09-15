package mail

type Sender interface {
	SendTo(to []string, title, content string) error
}
