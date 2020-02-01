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
	"fmt"
	"forgolang_forum/cmn"
	"forgolang_forum/database"
	"forgolang_forum/database/model"
	model2 "forgolang_forum/model"
	"github.com/fate-lovely/phi"
	"github.com/valyala/fasthttp"
	"strconv"
)

// PostCommentController post comment api controller
type PostCommentController struct {
	Controller
	*API
	Model model.PostComment
}

// Index list all post comments
func (c PostCommentController) Index(ctx *fasthttp.RequestCtx) {
	paginate, _, _ := c.Paginate(ctx, "id", "inserted_at")

	var comments []model.PostComment
	c.GetDB().QueryWithModel(fmt.Sprintf(`
		SELECT * FROM %s AS c
		WHERE c.post_id = $1
		ORDER BY %s %s
		LIMIT $2 OFFSET $3
	`, c.Model.TableName(), paginate.OrderField, paginate.OrderBy),
		&comments,
		phi.URLParam(ctx, "postID"),
		paginate.Limit,
		paginate.Offset)

	var count int64
	if count, _ = c.GetCache().Get(fmt.Sprintf("%s:%s",
		cmn.GetRedisKey("comment", "count"),
		phi.URLParam(ctx, "postID"))).Int64(); count <= 0 {
		c.GetDB().DB.Get(&count, fmt.Sprintf(`
			SELECT count(c.id) FROM %s AS c
			WHERE c.post_id = $1
		`, c.Model.TableName()),
			phi.URLParam(ctx, "postID"))

		c.GetCache().Set(fmt.Sprintf("%s:%s",
			cmn.GetRedisKey("comment", "count"),
			phi.URLParam(ctx, "postID")),
			count, 0)
	}

	c.JSONResponse(ctx, model2.ResponseSuccess{
		Data:       comments,
		TotalCount: count,
	}, fasthttp.StatusOK)
}

// Create post comment
func (c PostCommentController) Create(ctx *fasthttp.RequestCtx) {
	postID, err := strconv.ParseInt(phi.URLParam(ctx, "postID"), 10, 64)
	if err != nil {
		c.JSONResponse(ctx, model2.ResponseError{
			Detail: fasthttp.StatusMessage(fasthttp.StatusBadRequest),
		}, fasthttp.StatusBadRequest)
		return
	}

	postComment := new(model.PostComment)
	c.JSONBody(ctx, &postComment)
	postComment.PostID = postID
	postComment.UserID = c.GetAuthContext(ctx).ID

	err = c.GetDB().Insert(new(model.PostComment), postComment, "id", "inserted_at")
	if errs, err := database.ValidateConstraint(err, postComment); err != nil {
		c.JSONResponse(ctx, model2.ResponseError{
			Errors: errs,
			Detail: fasthttp.StatusMessage(fasthttp.StatusUnprocessableEntity),
		}, fasthttp.StatusUnprocessableEntity)
		return
	}

	c.GetCache().Incr(fmt.Sprintf("%s:%d",
		cmn.GetRedisKey("comment", "count"),
		postID))

	c.JSONResponse(ctx, model2.ResponseSuccessOne{
		Data: postComment,
	}, fasthttp.StatusCreated)
}

// Delete post comment
func (c PostCommentController) Delete(ctx *fasthttp.RequestCtx) {
	c.GetDB().Delete(c.Model.TableName(), "post_id = $1 AND id = $2",
		phi.URLParam(ctx, "postID"),
		phi.URLParam(ctx, "commentID")).Force()

	c.GetCache().Decr(fmt.Sprintf("%s:%s",
		cmn.GetRedisKey("comment", "count"),
		phi.URLParam(ctx, "postID")))

	c.JSONResponse(ctx, nil, fasthttp.StatusNoContent)
}
