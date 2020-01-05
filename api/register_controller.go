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

	user := model2.NewUser(registerRequest.Password)
	user.Email = registerRequest.Email
	user.Username = registerRequest.Username
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

	c.JSONResponse(ctx, model.ResponseSuccessOne{
		Data: user,
	}, fasthttp.StatusCreated)
}
