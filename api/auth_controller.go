// +build !test
// Copyright 2019 Abdulkadir DILSIZ
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

package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"forgolang_forum/database"
	model2 "forgolang_forum/database/model"
	github2 "forgolang_forum/thirdparty/github"
	"github.com/google/go-github/v28/github"
	"github.com/valyala/fasthttp"
)

// AuthController third-party authentication callback controller
type AuthController struct {
	Controller
	*API
}

// Github redirect github oauth page
func (c AuthController) Github(ctx *fasthttp.RequestCtx) {
	ctx.Redirect(c.App.Github.URL(), fasthttp.StatusTemporaryRedirect)
}

// GithubCallback github oauth callback method
func (c AuthController) GithubCallback(ctx *fasthttp.RequestCtx) {
	var thirdParty model2.ThirdParty

	c.App.Database.QueryRowWithModel(fmt.Sprintf("SELECT t.* FROM %s AS t "+
		"WHERE t.code = 'github' AND t.type = 'auth' AND t.is_active = true",
		thirdParty.TableName()),
		&thirdParty).Force()

	state := ctx.FormValue("state")
	if string(state) != c.App.Github.State {
		defaultLogger.LogError(errors.New("github auth state does not match"),
			fmt.Sprintf("%s: %s", state, c.App.Github.State))
		ctx.Redirect("/", fasthttp.StatusTemporaryRedirect)
		return
	}

	code := ctx.FormValue("code")
	token, err := c.App.Github.OauthConfig.Exchange(context.TODO(), string(code))
	if err != nil {
		defaultLogger.LogError(err, fmt.Sprintf("github exchange failed"))
		ctx.Redirect("/", fasthttp.StatusTemporaryRedirect)
		return
	}

	oauthClient := c.App.Github.OauthConfig.Client(context.TODO(), token)
	client := github.NewClient(oauthClient)
	githubUser, _, err := client.Users.Get(context.Background(), "")
	if err != nil {
		defaultLogger.LogError(err, fmt.Sprintf("github user get failed"))
		ctx.Redirect("/", fasthttp.StatusTemporaryRedirect)
		return
	}

	user := model2.NewUser(nil)
	user.Email = githubUser.GetEmail()
	user.Username = githubUser.GetLogin()
	user.Avatar.SetValid(githubUser.GetAvatarURL())
	user.IsActive = true

	passphrase := model2.NewUserPassphrase(0)
	db := c.App.Database.Transaction(func(tx *database.Tx) error {
		var currentUser model2.User
		tx.DB.QueryRowWithModel(fmt.Sprintf("SELECT u.* FROM %s AS u "+
			"WHERE u.email = $1",
			currentUser.TableName()),
			&currentUser,
			user.Email)

		var currentComebackApp model2.UserComebackApp
		if currentUser.ID > int64(0) {
			passphrase.UserID = currentUser.ID
			if err := tx.DB.Update(&currentUser, user, nil,
				"id", "inserted_at", "updated_at"); err != nil {
				defaultLogger.LogError(err, "current user error")
				return err
			}
			tx.DB.QueryRowWithModel(fmt.Sprintf("SELECT c.* FROM %s AS c "+
				"WHERE c.user_id = $1 AND c.tparty_id = $2",
				currentComebackApp.TableName()),
				&currentComebackApp,
				currentUser.ID,
				thirdParty.ID)
		} else {
			if err := tx.DB.Insert(new(model2.User), user, "id", "inserted_at", "updated_at"); err != nil {
				defaultLogger.LogError(err, "github user insert error")
				return err
			}
			passphrase.UserID = user.ID
		}

		if currentUser.ID == int64(0) {
			roleAssignment := model2.NewUserRoleAssignment(user.ID, 3)
			if err := tx.DB.Insert(new(model2.UserRoleAssignment), roleAssignment, "id"); err != nil {
				return err
			}

			userState := model2.NewUserState(user.ID)
			userState.State = database.Active
			userState.SourceUserID.SetValid(user.ID)
			if err := tx.DB.Insert(new(model2.UserState), userState, "id"); err != nil {
				return err
			}
		}

		comebackApp := model2.NewUserComebackApp(user.ID, thirdParty.ID)
		comebackApp.AccessToken = token.AccessToken
		comebackApp.RefreshToken.SetValid(token.RefreshToken)
		comebackApp.Expire.SetValid(token.Expiry.UnixNano())

		userInformation := new(github2.UserInformation)
		userInformation.Bio = githubUser.GetBio()
		userInformation.Followers = githubUser.GetFollowers()
		userInformation.Following = githubUser.GetFollowing()
		userInformation.PublicRepos = githubUser.GetPublicRepos()
		userInformation.PublicGists = githubUser.GetPublicGists()
		b, _ := json.Marshal(userInformation)
		comebackApp.Data.Scan(b)

		var err error
		if currentComebackApp.ID > int64(0) {
			err = tx.DB.Update(&currentComebackApp, comebackApp, nil,
				"id", "updated_at")
		} else {
			err = tx.DB.Insert(new(model2.UserComebackApp), comebackApp, "id")
		}

		if err != nil {
			return err
		}

		return tx.DB.Insert(new(model2.UserPassphrase), passphrase, "id")
	})

	if db.Error != nil {
		defaultLogger.LogError(err, fmt.Sprintf("github user get failed"))
		ctx.Redirect("https://forgolang.com/login?action=thirdy-party&type=github&status=failed", fasthttp.StatusTemporaryRedirect)
		return
	}

	ctx.Redirect(fmt.Sprintf("https://forgolang.com/login?token=%s&action=third-party&type=github&status=success",
		passphrase.Passphrase), fasthttp.StatusTemporaryRedirect)
}
