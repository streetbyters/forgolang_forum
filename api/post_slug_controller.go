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
	"forgolang_forum/database"
	model2 "forgolang_forum/database/model"
	"forgolang_forum/model"
	"github.com/fate-lovely/phi"
	"github.com/valyala/fasthttp"
	"strconv"
)

// PostSlugController post links api controller
type PostSlugController struct {
	Controller
	*API
}

//
func (c PostSlugController) Create(ctx *fasthttp.RequestCtx) {
	postID, err := strconv.ParseInt(phi.URLParam(ctx, "postID"), 10, 64)
	if err != nil {
		c.JSONResponse(ctx, model.ResponseError{
			Errors: nil,
			Detail: fasthttp.StatusMessage(fasthttp.StatusBadRequest),
		}, fasthttp.StatusBadRequest)
		return
	}

	postSlug := new(model2.PostSlug)
	c.JSONBody(ctx, &postSlug)
	postSlug.PostID = postID
	postSlug.SourceUserID.SetValid(c.GetAuthContext(ctx).ID)
	if errs, err := database.ValidateStruct(postSlug); err != nil {
		c.JSONResponse(ctx, model.ResponseError{
			Errors: errs,
			Detail: fasthttp.StatusMessage(fasthttp.StatusUnprocessableEntity),
		}, fasthttp.StatusUnprocessableEntity)
		return
	}

	err = c.GetDB().Insert(new(model2.PostSlug), postSlug, "id", "inserted_at")
	if errs, err := database.ValidateConstraint(err, postSlug); err != nil {
		c.JSONResponse(ctx, model.ResponseError{
			Errors: errs,
			Detail: fasthttp.StatusMessage(fasthttp.StatusUnprocessableEntity),
		}, fasthttp.StatusUnprocessableEntity)
		return
	}

	c.JSONResponse(ctx, model.ResponseSuccessOne{
		Data: postSlug,
	}, fasthttp.StatusCreated)
}
