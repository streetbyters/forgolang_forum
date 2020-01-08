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
	"github.com/valyala/fasthttp"
	"net/http"
)

// TokenController user authentication token controller
type TokenController struct {
	Controller
	*API
}

// Create generate user jwt method
func (c TokenController) Create(ctx *fasthttp.RequestCtx) {
	tokenRequest := new(model.TokenRequest)

	c.JSONBody(ctx, &tokenRequest)
	if errs, err := database.ValidateStruct(tokenRequest); err != nil {
		c.JSONResponse(ctx, model.ResponseError{
			Errors: errs,
			Detail: http.StatusText(http.StatusUnprocessableEntity),
		}, http.StatusUnprocessableEntity)
		return
	}

	passphrase := new(model2.UserPassphrase)
	result := c.App.Database.QueryRowWithModel(passphrase.PassphraseQuery(c.App.Database),
		passphrase,
		tokenRequest.Passphrase)
	if result.Error != nil {
		c.JSONResponse(ctx, model.ResponseError{
			Detail: fasthttp.StatusMessage(fasthttp.StatusNotFound),
		}, fasthttp.StatusNotFound)
		return
	}

	user := new(model2.User)
	result = c.App.Database.QueryRowWithModel("SELECT u.* FROM "+user.TableName()+" AS u "+
		"WHERE u.id = $1 AND u.is_active = true",
		user,
		passphrase.UserID)
	if result.Error != nil {
		c.JSONResponse(ctx, model.ResponseError{
			Detail: fasthttp.StatusMessage(fasthttp.StatusNotFound),
		}, fasthttp.StatusNotFound)
		return
	}

	role := new(model2.Role)
	roleAssignment := new(model2.UserRoleAssignment)
	err := c.App.Database.DB.QueryRowx(
		fmt.Sprintf("SELECT r.code, ra.role_id FROM %s AS ra "+
			"LEFT OUTER JOIN %s AS ra2 ON ra.user_id = ra2.user_id and ra.id < ra2.id "+
			"INNER JOIN %s AS r ON ra.role_id = r.id "+
			"WHERE ra2.id IS NULL AND ra.user_id = $1",
			roleAssignment.TableName(),
			roleAssignment.TableName(),
			role.TableName()),
		user.ID,
	).Scan(&role.Code, &roleAssignment.RoleID)
	if err != nil {
		c.JSONResponse(ctx, model.ResponseError{
			Detail: fasthttp.StatusMessage(fasthttp.StatusNotFound),
		}, fasthttp.StatusNotFound)
		return
	}

	jwt, _ := c.API.JWTAuth.Generate(user.ID, roleAssignment.RoleID, role.Code)

	c.JSONResponse(ctx, model.ResponseSuccessOne{
		Data: model.ResponseToken{
			JWT:    jwt,
			UserID: user.ID,
			Role:   role.Code,
		},
	}, http.StatusCreated)
}
