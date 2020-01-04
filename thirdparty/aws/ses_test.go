package aws

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

func TestNewSES(t *testing.T) {
	workingPath, _ := os.Getwd()
	dirs := strings.SplitAfter(workingPath, "forgolang_forum")
	workingPath = dirs[0]

	c, _ := Config(workingPath, "test.env")

	s, err := NewSES(c)

	assert.Nil(t, err)
	assert.NotNil(t, s.Svc)
}

func TestSES_Send(t *testing.T) {
	workingPath, _ := os.Getwd()
	dirs := strings.SplitAfter(workingPath, "forgolang_forum")
	workingPath = dirs[0]

	c, _ := Config(workingPath, "test.env")

	s, err := NewSES(c)
	assert.Nil(t, err)

	s.Subject = "Test"
	s.Recipient = []string{
		"akdilsiz@tecpor.com",
	}
	s.Body = "Body"

	err = s.Send()

	assert.Nil(t, err)
	assert.NotNil(t, s.Result["MessageId"])
	assert.NotNil(t, s.Result["String"])
}

func TestSES_SendWithInvalidSender(t *testing.T) {
	workingPath, _ := os.Getwd()
	dirs := strings.SplitAfter(workingPath, "forgolang_forum")
	workingPath = dirs[0]

	c, _ := Config(workingPath, "test.env")

	s, err := NewSES(c)
	assert.Nil(t, err)
	s.Sender = "invalid@domain.com"
	s.Body = "Body"
	s.Recipient = []string{
		"akdilsiz@tecpor.com",
	}

	err = s.Send()
	assert.NotNil(t, err)
	aerr := err.(awserr.Error)
	assert.Equal(t, aerr.Code(), ses.ErrCodeMessageRejected)
}
