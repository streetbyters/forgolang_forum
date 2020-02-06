// Copyright 2019 StreetByters Community
// Licensed to the Apache Software Foundation (ASF) under one or more
// contributor license agreements.  See the NOTICE file distributed with
// this work for additional information regarding copyright ownership.
// The ASF licenses this file to You under the Apache License, Version 2.0
// (the "License"); you may not use this file except in compliance with
// the License.  You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aws

import (
	"forgolang_forum/thirdparty"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

// SES third-party aws ses structure
type SES struct {
	thirdparty.Email
	Config    map[string]string
	Sender    string   `json:"sender" validate:"required"`
	Recipient []string `json:"recipient" validate:"required"`
	Subject   string   `json:"subject" validate:"required"`
	Body      string   `json:"body" validate:"required"`
	Charset   string   `json:"charset"`
	Svc       *ses.SES
	Result    map[string]interface{}
}

// NewSES generate ses
func NewSES(config map[string]string) (*SES, error) {
	s := new(SES)
	s.Charset = "UTF-8"
	s.Config = config
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(s.Config["ses_region"]),
		Credentials: credentials.NewStaticCredentials(s.Config["ses_access_key_id"], s.Config["ses_secret_access_key"], ""),
	})
	s.Sender = config["ses_source"]

	if err != nil {
		return nil, err
	}

	s.Svc = ses.New(sess)

	return s, nil
}

// Send email with given parameters
func (e *SES) Send() error {
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: aws.StringSlice(e.Recipient),
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(e.Charset),
					Data:    aws.String(e.Body),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(e.Charset),
				Data:    aws.String(e.Subject),
			},
		},
		Source: aws.String(e.Sender),
	}

	result, err := e.Svc.SendEmail(input)

	if err != nil {
		return err
	}

	e.Result = make(map[string]interface{})
	e.Result["Recipient"] = e.Recipient
	e.Result["MessageId"] = aws.StringValue(result.MessageId)
	e.Result["String"] = result.GoString()

	return nil
}
