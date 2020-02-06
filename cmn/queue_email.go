// Copyright 2019 Street Byters Community
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

package cmn

import (
	"encoding/json"
	"fmt"
	"forgolang_forum/utils"
	"github.com/streadway/amqp"
	"path/filepath"
)

// Email queue
type QueueEmail struct {
	Queue       *Queue
	AMQPChannel *amqp.Channel
	AMQPQueue   amqp.Queue
}

// QueueEmailBody received message body structure
type QueueEmailBody struct {
	Subject    string      `json:"subject"`
	Recipients []string    `json:"recipients" validate:"required"`
	Type       string      `json:"type" validate:"required"`
	Template   string      `json:"template" validate:"required"`
	Params     interface{} `json:"params,omitempty"`
}

func (s QueueEmailBody) ToJSON() string {
	b, err := json.Marshal(s)
	if err != nil {
		return ""
	}
	return string(b)
}

// NewEmail generate email queue structure
func NewEmail(queue *Queue) *QueueEmail {
	e := &QueueEmail{Queue: queue}

	channel, err := e.Queue.App.Amqp.Channel()
	FailOnError(e.Queue.App.Logger, err)

	e.AMQPChannel = channel
	err = e.AMQPChannel.ExchangeDeclare(
		"send_email_exchange",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	FailOnError(queue.App.Logger, err)

	args := make(amqp.Table)
	args["x-dead-letter-exchange"] = "send_email_exchange"
	args["x-dead-letter-routing-key"] = "send_email_error"

	e.AMQPQueue, err = e.AMQPChannel.QueueDeclare(
		"send_email",
		true,
		false,
		false,
		false,
		args,
	)
	FailOnError(queue.App.Logger, err)
	//if err != nil {
	//	e.Queue.App.Logger.LogError(err, "rabbitMQ error channel queue declare")
	//}

	return e
}

func (e QueueEmail) Publish(message string) error {
	err := e.AMQPChannel.Publish(
		"send_email_exchange",
		"send_email_error",
		false,
		false,
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "application/json",
			ContentEncoding: "",
			Body:            []byte(message),
			DeliveryMode:    amqp.Transient,
			Priority:        9,
		},
	)
	if err != nil {
		e.Queue.App.Logger.LogError(err, fmt.Sprintf("%s: message publish error", e.AMQPQueue.Name))
		return err
	}
	e.Queue.App.Logger.LogInfo(fmt.Sprintf("%s: message publish", e.AMQPQueue.Name))
	return err
}

// Start email queue
func (e QueueEmail) Start() {
	err := e.AMQPChannel.QueueBind(
		e.AMQPQueue.Name,
		"send_email_error",
		"send_email_exchange",
		false,
		nil)
	if err != nil {
		e.Queue.App.Logger.LogError(err, "rabbitMQ error channel queue bind")
	}

	received, err := e.AMQPChannel.Consume(
		e.AMQPQueue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		e.Queue.App.Logger.LogError(err, "rabbitMQ error channel queue consume")
	}

	e.Queue.App.Logger.LogInfo("Start and subscribe email queue")

	go func() {
		for receive := range received {
			e.Receive(receive.Body)
		}
	}()

	<-e.Queue.App.Channel

	defer e.AMQPChannel.Close()
}

func (e QueueEmail) Receive(body []byte) {
	var receivedBody QueueEmailBody
	err := json.Unmarshal(body, &receivedBody)
	if err != nil {
		e.Queue.App.Logger.LogError(err, "QueueEmail body error")
		return
	}

	if errs, err := utils.ValidateStruct(receivedBody); err != nil {
		ss, _ := json.Marshal(errs)
		e.Queue.App.Logger.LogError(err, fmt.Sprintf("QueueEmail validate error: %s", string(ss)))
		return
	}

	template := utils.ReadFileToTemplate(filepath.Join(e.Queue.App.Config.Path,
		"mail", "template", fmt.Sprintf("%s.html", receivedBody.Template)),
		receivedBody.Params)

	if template == "" {
		return
	}

  	email := e.Queue.App.Email
	email.Recipient = receivedBody.Recipients
	email.Subject = receivedBody.Subject
	email.Body = template
	email.Send()
}
