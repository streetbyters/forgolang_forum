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
	"errors"
	"fmt"
	"forgolang_forum/database"
	model2 "forgolang_forum/database/model"
	"forgolang_forum/model"
	"github.com/fate-lovely/phi"
	"github.com/gosimple/slug"
	"github.com/valyala/fasthttp"
	"strconv"
)

type PostDetailController struct {
	Controller
	*API
}

func (c PostDetailController) Create(ctx *fasthttp.RequestCtx) {
	var err error

	postID, err := strconv.ParseInt(phi.URLParam(ctx, "postID"), 10, 64)
	if err != nil {
		c.JSONResponse(ctx, model.ResponseError{
			Detail: fasthttp.StatusMessage(fasthttp.StatusBadRequest),
		}, fasthttp.StatusBadRequest)
		return
	}

	postDetail := new(model2.PostDetail)
	c.JSONBody(ctx, &postDetail)
	postDetail.PostID = postID
	postDetail.SourceUserID.SetValid(c.GetAuthContext(ctx).ID)

	if errs, err := database.ValidateStruct(postDetail); err != nil {
		c.JSONResponse(ctx, model.ResponseError{
			Errors: errs,
			Detail: fasthttp.StatusMessage(fasthttp.StatusUnprocessableEntity),
		}, fasthttp.StatusUnprocessableEntity)
		return
	}

	postDetail.Title = c.App.TextPolicy.Sanitize(postDetail.Title)
	postDetail.Description.SetValid(c.App.TextPolicy.Sanitize(postDetail.Description.String))
	postDetail.Content = c.App.TextPolicy.Sanitize(postDetail.Content)

	postSlug := model2.NewPostSlug(postID, c.GetAuthContext(ctx).ID)
	errs := make(map[string]string)
	c.GetDB().Transaction(func(tx *database.Tx) error {
		result := c.GetDB().QueryRow(fmt.Sprintf(`
			SELECT * FROM %s AS ps
			LEFT OUTER JOIN %s AS ps2 ON ps.post_id = ps2.post_id AND ps.id < ps2.id
			WHERE ps2.id IS NULL AND ps.post_id != $1 AND ps.slug = $2
		`, postSlug.TableName(), postSlug.TableName()),
			postID,
			slug.Make(postDetail.Title))
		if result.Count > 0 {
			err = errors.New("slug has been already taken")
			// TODO: error field!!!
			errs["title"] = "has been already taken"
			return err
		}

		err = c.GetDB().Insert(new(model2.PostDetail), postDetail, "id", "inserted_at")
		if errs, err = database.ValidateConstraint(err, postDetail); err != nil {
			return err
		}

		postSlug.Slug = slug.Make(postDetail.Title)
		return c.GetDB().Insert(new(model2.PostSlug), postSlug, "id")
	})

	if err != nil {
		c.JSONResponse(ctx, model.ResponseError{
			Errors: errs,
			Detail: fasthttp.StatusMessage(fasthttp.StatusUnprocessableEntity),
		}, fasthttp.StatusUnprocessableEntity)
		return
	}

	c.App.ElasticClient.Update().
		Index("posts").
		Id(strconv.FormatInt(postID, 10)).
		Doc(map[string]interface{}{
			"title":       postDetail.Title,
			"content":     postDetail.Content,
			"description": postDetail.Description,
			"slug":        postSlug.Slug,
		}).
		Do(context.TODO())

	c.JSONResponse(ctx, model.ResponseSuccessOne{
		Data: postDetail,
	}, fasthttp.StatusCreated)
}
