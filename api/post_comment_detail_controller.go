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
	"forgolang_forum/database"
	"forgolang_forum/database/model"
	model2 "forgolang_forum/model"
	"forgolang_forum/utils"
	"github.com/fate-lovely/phi"
	"github.com/valyala/fasthttp"
)

// PostCommentDetailController post comment detail api controller
type PostCommentDetailController struct {
	Controller
	*API
}

// Create post comment detail
func (c PostCommentDetailController) Create(ctx *fasthttp.RequestCtx) {
	var e []bool
	postID, notExists := utils.ParseInt(phi.URLParam(ctx, "postID"), 10, 64)
	e = append(e, notExists)
	commentID, notExists := utils.ParseInt(phi.URLParam(ctx, "commentID"), 10, 64)
	e = append(e, notExists)

	if exists, _ := utils.InArray(true, e); exists {
		c.JSONResponse(ctx, model2.ResponseError{
			Detail: fasthttp.StatusMessage(fasthttp.StatusBadRequest),
		}, fasthttp.StatusBadRequest)
		return
	}

	commentDetail := new(model.PostCommentDetail)
	c.JSONBody(ctx, &commentDetail)
	commentDetail.PostID = postID
	commentDetail.CommentID = commentID

	if errs, err := database.ValidateStruct(commentDetail); err != nil {
		c.JSONResponse(ctx, model2.ResponseError{
			Errors: errs,
			Detail: fasthttp.StatusMessage(fasthttp.StatusUnprocessableEntity),
		}, fasthttp.StatusUnprocessableEntity)
		return
	}

	commentDetail.Comment = c.App.TextPolicy.Sanitize(commentDetail.Comment)

	err := c.GetDB().Insert(new(model.PostCommentDetail), commentDetail, "id", "inserted_at")
	if errs, err := database.ValidateConstraint(err, commentDetail); err != nil {
		c.JSONResponse(ctx, model2.ResponseError{
			Errors: errs,
			Detail: fasthttp.StatusMessage(fasthttp.StatusUnprocessableEntity),
		}, fasthttp.StatusUnprocessableEntity)
		return
	}

	c.JSONResponse(ctx, model2.ResponseSuccessOne{
		Data: commentDetail,
	}, fasthttp.StatusCreated)
}
