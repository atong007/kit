package mail

import (
	"testing"
)

func TestSmtpSender_SendTo(t *testing.T) {
	type fields struct {
		host        string
		port        string
		defaultMail string
		passwd      string
	}
	type args struct {
		to      []string
		title   string
		content string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"success",
			fields{
				host:        "mail.example.com",
				port:        "25",
				defaultMail: "your mail",
				passwd:      "your password",
			},
			args{
				to:      []string{"user@example.com"},
				title:   "this is the mail subject",
				content: "This is the mail content",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSmtpSender(tt.fields.host,
				tt.fields.port,
				tt.fields.defaultMail,
				tt.fields.passwd,
			)
			if err := s.SendTo(tt.args.to, tt.args.title, tt.args.content); (err != nil) != tt.wantErr {
				t.Errorf("SendTo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
