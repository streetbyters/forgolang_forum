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

package aws

import (
	"bytes"
	"github.com/spf13/viper"
	"io/ioutil"
	"path/filepath"
)

// Config aws package
func Config(workingPath string, envFile string) (map[string]string, error) {
	viper.SetConfigName(envFile)
	viper.AddConfigPath(workingPath)
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	file, err := ioutil.ReadFile(filepath.Join(workingPath, "secret.env"))
	if err != nil {
		panic(err)
	}
	buffer := bytes.NewReader(file)
	err = viper.MergeConfig(buffer)
	if err != nil {
		panic(err)
	}

	config := make(map[string]string)
	config["ses_access_key_id"] = viper.GetString("AWS_SES_ACCESS_KEY_ID")
	config["ses_secret_access_key"] = viper.GetString("AWS_SES_SECRET_ACCESS_KEY")
	config["ses_region"] = viper.GetString("AWS_SES_REGION")
	config["ses_source"] = viper.GetString("AWS_SES_SOURCE")

	config["s3_region"] = viper.GetString("AWS_S3_REGION")
	config["s3_bucket"] = viper.GetString("AWS_S3_BUCKET")
	config["s3_access_key_id"] = viper.GetString("AWS_S3_ACCESS_KEY_ID")
	config["s3_secret_access_key"] = viper.GetString("AWS_S3_SECRET_ACCESS_KEY")

	return config, nil
}
