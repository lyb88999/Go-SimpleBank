package mail

import (
	"github.com/lyb88999/Go-SimpleBank/util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSendEmailWithQQEmail(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	config, err := util.LoadConfig("..")
	require.NoError(t, err)
	qqEmailSender := NewQQEmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)
	subject := "a test email"
	content := `
	<h1>Hello world!</h1>
	`
	to := []string{"354083501@qq.com"}
	attachFiles := []string{"../app.env"}
	err = qqEmailSender.SendEmail(subject, content, to, nil, nil, attachFiles)
	require.NoError(t, err)
}
