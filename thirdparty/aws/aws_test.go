package aws

import (
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

func Test_ConfigWithoutSender(t *testing.T) {
	workingPath, _ := os.Getwd()
	dirs := strings.SplitAfter(workingPath, "forgolang_forum")
	workingPath = dirs[0]

	config, err := Config(workingPath, "test")

	assert.Nil(t, err)
	assert.NotNil(t, config["ses_access_key_id"])
	assert.NotNil(t, config["ses_secret_key_id"])
	assert.Equal(t, config["ses_region"], "eu-central-1")
	assert.Equal(t, config["ses_source"], "noreply@forgolang.com")
}

//
//func Test_ConfigWithSender(t *testing.T) {
//	var ses SES
//	ses.Sender = "akdilsiz@tecpor.com"
//	config, err := Config(&ses)
//
//	assert.Nil(t, err)
//	assert.NotNil(t, config["access_key_id"])
//	assert.NotNil(t, config["secret_key_id"])
//	assert.Equal(t, config["ses_region"], "eu-west-1")
//	assert.Equal(t, config["ses_source"], "akdilsiz@tecpor.com")
//}
//
//func Test_Config(t *testing.T) {
//	var email thirdparty.Email
//	config, err := Config(email)
//
//	assert.Nil(t, err)
//	assert.NotNil(t, config["access_key_id"])
//	assert.NotNil(t, config["secret_key_id"])
//	assert.Equal(t, config["ses_region"], "eu-west-1")
//	assert.Equal(t, config["ses_source"], "noreply@tecpor.com")
//}
