// Copyright 2019 Forgolang Community
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
	"forgolang_forum/cmn"
	"forgolang_forum/database"
	model2 "forgolang_forum/database/model"
	"forgolang_forum/model"
	"github.com/valyala/fasthttp"
)

// RegisterController user register controller
type RegisterController struct {
	Controller
	*API
}

// Create user register method
func (c RegisterController) Create(ctx *fasthttp.RequestCtx) {
	var registerRequest model.RegisterRequest

	c.JSONBody(ctx, &registerRequest)
	if errs, err := database.ValidateStruct(registerRequest); err != nil {
		c.JSONResponse(ctx, model.ResponseError{
			Errors: errs,
			Detail: fasthttp.StatusMessage(fasthttp.StatusUnprocessableEntity),
		}, fasthttp.StatusUnprocessableEntity)
		return
	}

	user := model2.NewUser(&registerRequest.Password)
	user.Email = registerRequest.Email
	user.Username = registerRequest.Username
	user.IsActive = true
	err := c.App.Database.Insert(new(model2.User),
		user,
		"id", "inserted_at", "updated_at")
	if errs, err := database.ValidateConstraint(err, user); err != nil {
		c.JSONResponse(ctx, model.ResponseError{
			Errors: errs,
			Detail: fasthttp.StatusMessage(fasthttp.StatusUnprocessableEntity),
		}, fasthttp.StatusUnprocessableEntity)
		return
	}

	otc := model2.NewUserOneTimeCode(user.ID)
	otc.Type = database.Confirmation
	c.App.Database.Insert(new(model2.UserOneTimeCode), otc, "id")

	userState := model2.NewUserState(user.ID)
	userState.State = database.WaitForConfirmation
	userState.SourceUserID.SetValid(user.ID)
	c.App.Database.Insert(new(model2.UserState), userState, "id")

	registerResponse := new(model.RegisterResponse)
	registerResponse.UserID = user.ID
	registerResponse.State = string(database.WaitForConfirmation)
	registerResponse.InsertedAt = user.InsertedAt

	roleAssignment := model2.NewUserRoleAssignment(user.ID, 3)
	roleAssignment.SourceUserID.SetValid(user.ID)
	c.App.Database.Insert(new(model2.UserRoleAssignment), roleAssignment, "id")

	go func() {
		c.App.Queue.Email.Publish(cmn.QueueEmailBody{
			Recipients: []string{user.Email},
			Subject:    "Forgolang.com | Activation Required",
			Type:       "confirmation",
			Template:   "confirmation",
			Params: struct {
				UserID int64
				Code   string
			}{
				UserID: user.ID,
				Code:   otc.Code,
			},
		}.ToJSON())
	}()

	c.JSONResponse(ctx, model.ResponseSuccessOne{
		Data: registerResponse,
	}, fasthttp.StatusCreated)
}
