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
	"context"
	"forgolang_forum/database"
	model2 "forgolang_forum/database/model"
	"forgolang_forum/model"
	"github.com/fate-lovely/phi"
	"github.com/olivere/elastic/v7"
	"github.com/valyala/fasthttp"
	"strconv"
)

type PostCategoryAssignmentController struct {
	Controller
	*API
}

func (c PostCategoryAssignmentController) Create(ctx *fasthttp.RequestCtx) {
	postID, err := strconv.ParseInt(phi.URLParam(ctx, "postID"), 10, 64)
	if err != nil {
		c.JSONResponse(ctx, model.ResponseError{
			Detail: fasthttp.StatusMessage(fasthttp.StatusBadRequest),
		}, fasthttp.StatusBadRequest)
		return
	}

	postCategoryAssignment := new(model2.PostCategoryAssignment)
	c.JSONBody(ctx, &postCategoryAssignment)
	postCategoryAssignment.PostID = postID

	if errs, err := database.ValidateStruct(postCategoryAssignment); err != nil {
		c.JSONResponse(ctx, model.ResponseError{
			Errors: errs,
			Detail: fasthttp.StatusMessage(fasthttp.StatusUnprocessableEntity),
		}, fasthttp.StatusUnprocessableEntity)
		return
	}

	err = c.GetDB().Insert(new(model2.PostCategoryAssignment),
		postCategoryAssignment,
		"id", "inserted_At")
	if errs, err := database.ValidateConstraint(err, postCategoryAssignment); err != nil {
		c.JSONResponse(ctx, model.ResponseError{
			Errors: errs,
			Detail: fasthttp.StatusMessage(fasthttp.StatusUnprocessableEntity),
		}, fasthttp.StatusUnprocessableEntity)
		return
	}

	_, err = c.App.ElasticClient.Update().
		Index("posts").
		Id(strconv.FormatInt(postID, 10)).
		Script(elastic.NewScript(`
			if(ctx._source.category_assignments) {
				ctx._source.category_assignments.append(params.cat)
			} else {
				ctx._source.category_assignments = []
				ctx._source.category_assignments.append(params.cat)
			}
			`).
			Params(map[string]interface{}{
				"cat": postCategoryAssignment.CategoryID,
			})).
		Do(context.TODO())

	c.JSONResponse(ctx, model.ResponseSuccessOne{
		Data: postCategoryAssignment,
	}, fasthttp.StatusCreated)
}
