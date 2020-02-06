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

package api

import (
	"context"
	"fmt"
	"forgolang_forum/cmn"
	"forgolang_forum/database"
	"forgolang_forum/database/model"
	model2 "forgolang_forum/model"
	"github.com/fate-lovely/phi"
	"github.com/valyala/fasthttp"
	"strconv"
	"time"
)

// UserRoleAssignmentController user role assignment api controller
type UserRoleAssignmentController struct {
	Controller
	*API
}

// Create method for user role assignment
func (c UserRoleAssignmentController) Create(ctx *fasthttp.RequestCtx) {
	var roleAssignment model.UserRoleAssignment
	c.JSONBody(ctx, &roleAssignment)

	i, err := strconv.ParseInt(phi.URLParam(ctx, "userID"), 10, 64)
	if err != nil {
		c.JSONResponse(ctx, model2.ResponseError{
			Detail: fasthttp.StatusMessage(fasthttp.StatusBadRequest),
		}, fasthttp.StatusBadRequest)
		return
	}
	roleAssignment.UserID = i

	if errs, err := database.ValidateStruct(roleAssignment); err != nil {
		c.JSONResponse(ctx, model2.ResponseError{
			Errors: errs,
			Detail: fasthttp.StatusMessage(fasthttp.StatusUnprocessableEntity),
		}, fasthttp.StatusUnprocessableEntity)
		return
	}

	roleAssignment.SourceUserID.SetValid(c.GetAuthContext(ctx).ID)
	err = c.GetDB().Insert(new(model.UserRoleAssignment), &roleAssignment, "id", "inserted_at")
	if errs, err := database.ValidateConstraint(err, &roleAssignment); err != nil {
		c.JSONResponse(ctx, model2.ResponseError{
			Errors: errs,
			Detail: fasthttp.StatusMessage(fasthttp.StatusUnprocessableEntity),
		}, fasthttp.StatusUnprocessableEntity)
		return
	}

	user := new(model.User)
	c.GetDB().QueryRowWithModel(fmt.Sprintf("%s AND u.id = $1", user.Query(false)),
		user,
		phi.URLParam(ctx, "userID")).Force()

	c.App.Cache.Set(fmt.Sprintf("%s:%d",
		cmn.GetRedisKey("user", "one"), user.ID),
		user.ToJSON(), time.Minute*30)

	c.App.ElasticClient.Index().Index("users").
		Id(strconv.FormatInt(user.ID, 10)).
		BodyJson(user).
		Do(context.TODO())

	c.JSONResponse(ctx, model2.ResponseSuccessOne{
		Data: roleAssignment,
	}, fasthttp.StatusCreated)
}
