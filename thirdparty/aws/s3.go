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
	"bytes"
	"errors"
	"forgolang_forum/thirdparty"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	awsS3 "github.com/aws/aws-sdk-go/service/s3"
	"mime/multipart"
	"net/http"
)

// S3 third-party aws s3 structure
type S3 struct {
	thirdparty.Storage
	Config map[string]string
	Svc    *awsS3.S3
}

// NewS3 generate aws s3 structure
func NewS3(config map[string]string) (*S3, error) {
	s := new(S3)
	s.Config = config
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(s.Config["s3_region"]),
		Credentials: credentials.NewStaticCredentials(s.Config["s3_access_key_id"], s.Config["s3_secret_access_key"], ""),
	})

	if err != nil {
		return nil, err
	}

	s.Svc = awsS3.New(sess)

	return s, nil
}

// Delete object s3 bucket
func (s *S3) Delete(args ...interface{}) error {
	if len(args) < 1 {
		return errors.New("args is not nil")
	}

	_, err := s.Svc.DeleteObject(&awsS3.DeleteObjectInput{
		Bucket: aws.String(s.Config["s3_bucket"]),
		Key:    aws.String(args[0].(string)),
	})

	return err
}

// Upload object s3 bucket
func (s *S3) Upload(args ...interface{}) error {
	if len(args) < 3 {
		return errors.New("args is not nil")
	}
	file := args[0].(*multipart.FileHeader)
	fileName := args[1].(string)
	acl := args[2].(string)

	f, _ := file.Open()
	buffer := make([]byte, file.Size)
	f.Read(buffer)
	defer f.Close()

	_, err := s.Svc.PutObject(&awsS3.PutObjectInput{
		Bucket:        aws.String(s.Config["s3_bucket"]),
		Key:           aws.String(fileName),
		ACL:           aws.String(acl),
		Body:          bytes.NewReader(buffer),
		ContentLength: aws.Int64(file.Size),
		ContentType:   aws.String(http.DetectContentType(buffer)),
	})

	return err
}
