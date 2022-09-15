package mail

type Sender interface {
	SendTo(mail string, content string) error
}
