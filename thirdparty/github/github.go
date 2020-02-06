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

package github

import (
	"forgolang_forum/model"
	baseGithub "github.com/google/go-github/v28/github"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
)

// Github forgolang.com github integration structure
type Github struct {
	Client      *baseGithub.Client
	OauthConfig *oauth2.Config
	State       string
}

// NewGithub generate github structure
func NewGithub(config *model.Config) *Github {
	c := new(Github)
	c.State = "forgolang.com"
	c.OauthConfig = &oauth2.Config{
		ClientID:     config.GithubClientID,
		ClientSecret: config.GithubClientSecret,
		Scopes:       []string{"read:user", "user:email", "user:follow"},
		Endpoint:     githuboauth.Endpoint,
	}

	return c
}

// URL generate oauth url
func (c *Github) URL() string {
	return c.OauthConfig.AuthCodeURL(c.State, oauth2.AccessTypeOnline)
}
