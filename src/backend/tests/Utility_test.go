package tests

import (
	utils "t_webdevapp/utils"
	"testing"
)

func TestMailer(t *testing.T) {
	data := map[string]interface{}{
		"User": "RDev",
	}

	file := "../template/emails/welcome.html"
	if content, err := utils.ParseHTML(file, data); err == nil {
		utils.Mailer("", "Welcome!", content)
	}
}
