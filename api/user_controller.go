// Copyright 2019 Abdulkadir Dilsiz - Çağatay Yücelen
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
	"fmt"
	"forgolang_forum/cmn"
	"forgolang_forum/database"
	"forgolang_forum/database/model"
	model2 "forgolang_forum/model"
	"github.com/akdilsiz/agente/utils"
	"github.com/fate-lovely/phi"
	"github.com/valyala/fasthttp"
	"strconv"
	"time"
)

// UserController Users api structure
type UserController struct {
	Controller
	*API
}

// Index List all users with paginate params
func (c UserController) Index(ctx *fasthttp.RequestCtx) {
	paginate, _, _ := c.Paginate(ctx, "id", "inserted_at")

	user := new(model.User)

	var users []model.User
	c.App.Database.QueryWithModel(fmt.Sprintf(`%s
		ORDER BY %s %s
		LIMIT $1 OFFSET $2
	`, user.Query(false), paginate.OrderField, paginate.OrderBy),
		&users,
		paginate.Limit, paginate.Offset)

	var count int64
	c.App.Database.DB.Get(&count,
		fmt.Sprintf("SELECT count(u.id) FROM %s AS u", user.TableName()))

	c.JSONResponse(ctx, model2.ResponseSuccess{
		Data:       users,
		TotalCount: count,
	}, fasthttp.StatusOK)
}

// Show user with given identifier
func (c UserController) Show(ctx *fasthttp.RequestCtx) {
	var user model.User

	var us string
	if err := c.App.Cache.Get(fmt.Sprintf("%s:%s",
		cmn.GetRedisKey("user", "one"),
		phi.URLParam(ctx, "userID"))).Scan(&us); err == nil {
		json.Unmarshal([]byte(us), &user)

		c.JSONResponse(ctx, model2.ResponseSuccessOne{
			Data: user,
		}, fasthttp.StatusOK)
		return
	}

	c.App.Database.QueryRowWithModel(fmt.Sprintf("%s AND u.id = $1", user.Query(false)),
		&user,
		phi.URLParam(ctx, "userID")).Force()

	c.App.Cache.Set(fmt.Sprintf("%s:%d",
		cmn.GetRedisKey("user", "one"), user.ID),
		user.ToJSON(), time.Minute*30)

	c.JSONResponse(ctx, model2.ResponseSuccessOne{
		Data: user,
	}, fasthttp.StatusOK)
}

// Create user with valid params
func (c UserController) Create(ctx *fasthttp.RequestCtx) {
	user := new(model.User)
	c.JSONBody(ctx, &user)

	if errs, err := database.ValidateStruct(user); err != nil {
		c.JSONResponse(ctx, model2.ResponseError{
			Errors: errs,
			Detail: fasthttp.StatusMessage(fasthttp.StatusUnprocessableEntity),
		}, fasthttp.StatusUnprocessableEntity)
		return
	}

	user.PasswordDigest.SetValid(utils.HashPassword(user.Password, 11))
	if user.Password == "" {
		c.JSONResponse(ctx, model2.ResponseError{
			Errors: map[string]string{
				"password": "is not nil",
			},
			Detail: fasthttp.StatusMessage(fasthttp.StatusUnprocessableEntity),
		}, fasthttp.StatusUnprocessableEntity)
		return
	}

	err := c.App.Database.Insert(new(model.User), user, "id", "inserted_at", "updated_at")
	if errs, err := database.ValidateConstraint(err, user); err != nil {
		c.JSONResponse(ctx, model2.ResponseError{
			Errors: errs,
			Detail: fasthttp.StatusMessage(fasthttp.StatusUnprocessableEntity),
		}, fasthttp.StatusUnprocessableEntity)
		return
	}

	user.Password = "****"

	c.App.ElasticClient.Index().Index("users").
		Id(strconv.FormatInt(user.ID, 10)).
		BodyJson(user).
		Do(context.TODO())

	c.JSONResponse(ctx, model2.ResponseSuccessOne{
		Data: user,
	}, fasthttp.StatusCreated)
}

// Update user with given identifier and valid params
func (c UserController) Update(ctx *fasthttp.RequestCtx) {
	user := new(model.User)
	c.App.Database.QueryRowWithModel(fmt.Sprintf("%s AND u.id = $1", user.Query(false)),
		user,
		phi.URLParam(ctx, "userID")).Force()

	var userRequest model.User
	c.JSONBody(ctx, &userRequest)

	if errs, err := database.ValidateStruct(userRequest); err != nil {
		c.JSONResponse(ctx, model2.ResponseError{
			Errors: errs,
			Detail: fasthttp.StatusMessage(fasthttp.StatusUnprocessableEntity),
		}, fasthttp.StatusUnprocessableEntity)
		return
	}

	if c.GetAuthContext(ctx).Role != "superadmin" {
		userRequest.IsActive = user.IsActive
	}

	userRequest.InsertedAt = user.InsertedAt

	err := c.App.Database.Update(user, &userRequest, nil, "id", "inserted_at", "updated_at")
	if errs, err := database.ValidateConstraint(err, user); err != nil {
		c.JSONResponse(ctx, model2.ResponseError{
			Errors: errs,
			Detail: fasthttp.StatusMessage(fasthttp.StatusUnprocessableEntity),
		}, fasthttp.StatusUnprocessableEntity)
		return
	}

	user.Password = "****"

	c.App.Cache.Set(fmt.Sprintf("%s:%d",
		cmn.GetRedisKey("user", "one"), user.ID),
		user.ToJSON(), time.Minute*30)

	go c.App.ElasticClient.Index().Index("users").
		Id(strconv.FormatInt(user.ID, 10)).
		BodyJson(user).Do(context.TODO())

	c.JSONResponse(ctx, model2.ResponseSuccessOne{
		Data: user,
	}, fasthttp.StatusOK)
}

// Delete user with given identifier
func (c UserController) Delete(ctx *fasthttp.RequestCtx) {
	var user model.User

	c.App.Database.Delete(user.TableName(),
		"id = $1",
		phi.URLParam(ctx, "userID")).Force()

	c.App.Cache.Del(fmt.Sprintf("%s:%s",
		cmn.GetRedisKey("user", "one"),
		phi.URLParam(ctx, "userID")))
	go c.App.ElasticClient.Delete().Index("users").
		Id(strconv.FormatInt(user.ID, 10)).
		Do(context.TODO())

	c.JSONResponse(ctx, model2.ResponseSuccessOne{
		Data: nil,
	}, fasthttp.StatusNoContent)
}
