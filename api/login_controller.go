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
	"fmt"
	"forgolang_forum/database"
	model2 "forgolang_forum/database/model"
	"forgolang_forum/model"
	"github.com/akdilsiz/agente/utils"
	"github.com/valyala/fasthttp"
)

// LoginController user authentication controller
type LoginController struct {
	Controller
	*API
}

// Create user sign in method
func (c LoginController) Create(ctx *fasthttp.RequestCtx) {
	var loginRequest model.LoginRequest

	c.JSONBody(ctx, &loginRequest)
	if errs, err := database.ValidateStruct(loginRequest); err != nil {
		c.JSONResponse(ctx, model.ResponseError{
			Errors: errs,
			Detail: fasthttp.StatusMessage(fasthttp.StatusUnprocessableEntity),
		}, fasthttp.StatusUnprocessableEntity)
		return
	}

	roleAssignment := new(model2.UserRoleAssignment)
	userModel := new(model2.User)
	userState := new(model2.UserState)
	c.App.Database.QueryRowWithModel(fmt.Sprintf(
		"SELECT u.* FROM %s AS u "+
			"INNER JOIN %s AS ra ON u.id = ra.user_id "+
			"LEFT OUTER JOIN %s AS ra2 ON ra.user_id = ra2.user_id and ra.id < ra2.id "+
			"INNER JOIN %s AS us ON u.id = us.user_id "+
			"LEFT OUTER JOIN %s AS us2 ON us.user_id = us2.user_id and us.id < us2.id "+
			"WHERE ra2.id IS NULL and us2.id IS NULL and us.state = $1 and u.username = $2 OR u.email = $2",
		userModel.TableName(),
		roleAssignment.TableName(),
		roleAssignment.TableName(),
		userState.TableName(),
		userState.TableName(),
	), userModel, database.Active, loginRequest.ID).Force()

	if err := utils.ComparePassword([]byte(userModel.PasswordDigest.String), []byte(loginRequest.Password)); err != nil {
		c.JSONResponse(ctx, model.ResponseError{
			Detail: "authentication failed",
		}, fasthttp.StatusUnauthorized)
		return
	}

	userPassphrase := new(model2.UserPassphrase)
	userPassphraseModel := model2.NewUserPassphrase(userModel.ID)
	c.App.Database.Insert(userPassphrase,
		userPassphraseModel, "id", "inserted_at")

	c.JSONResponse(ctx, model.ResponseSuccessOne{
		Data: model.LoginResponse{
			PassphraseID: userPassphrase.ID,
			UserID:       userModel.ID,
			Passphrase:   userPassphrase.Passphrase,
		},
	}, fasthttp.StatusCreated)
}
