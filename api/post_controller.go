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
	"forgolang_forum/database/model"
	model2 "forgolang_forum/model"
	"github.com/fate-lovely/phi"
	"github.com/gosimple/slug"
	"github.com/valyala/fasthttp"
	"strconv"
)

type PostController struct {
	Controller
	*API
}

// Create post with post deps
func (c PostController) Create(ctx *fasthttp.RequestCtx) {
	postReq := new(model.PostDEP)
	c.JSONBody(ctx, &postReq)

	postReq.AuthorID = c.GetAuthContext(ctx).ID
	if errs, err := database.ValidateStruct(postReq); err != nil {
		c.JSONResponse(ctx, model2.ResponseError{
			Errors: errs,
			Detail: fasthttp.StatusMessage(fasthttp.StatusUnprocessableEntity),
		}, fasthttp.StatusUnprocessableEntity)
		return
	}

	post := new(model.Post)
	postSlug := model.NewPostSlug(0, c.GetAuthContext(ctx).ID)
	postDetail := model.NewPostDetail(0, c.GetAuthContext(ctx).ID)

	var err error
	errs := make(map[string]string)
	c.App.Database.Transaction(func(tx *database.Tx) error {
		post.AuthorID = c.GetAuthContext(ctx).ID
		err = tx.DB.Insert(new(model.Post), post, "id")
		if errs, err = database.ValidateConstraint(err, post); err != nil {
			return err
		}

		result := c.App.Database.QueryRow(fmt.Sprintf(`
			SELECT * FROM %s AS ps
			LEFT OUTER JOIN %s AS ps2 ON ps.post_id = ps2.post_id AND ps.id < ps2.id
			WHERE ps2.id IS NULL and ps.slug = $1
		`, postSlug.TableName(), postSlug.TableName()),
			slug.Make(postReq.Title.String))
		if result.Count > 0 {
			err = errors.New("slug has been already taken")
			// TODO: error field!!!
			errs["slug"] = "has been already taken"
			//fmt.Println(errs)
			return err
		}

		postSlug.PostID = post.ID
		postSlug.Slug = slug.Make(postReq.Title.String)
		err = tx.DB.Insert(new(model.PostSlug), postSlug, "id")
		if errs, err = database.ValidateConstraint(err, postSlug); err != nil {
			return err
		}
		postReq.Slug.SetValid(postSlug.Slug)

		postDetail.Title = postReq.Title.String
		postDetail.Description = postReq.Description
		postDetail.Content = postReq.Content.String
		postDetail.PostID = post.ID
		err = tx.DB.Insert(new(model.PostDetail), postDetail, "id")
		if errs, err = database.ValidateConstraint(err, postDetail); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSONResponse(ctx, model2.ResponseError{
			Errors: errs,
			Detail: fasthttp.StatusMessage(fasthttp.StatusUnprocessableEntity),
		}, fasthttp.StatusUnprocessableEntity)
		return
	}

	postReq.ID = post.ID
	postReq.InsertedAt = post.InsertedAt

	c.App.ElasticClient.Index().
		Index("posts").
		Id(strconv.FormatInt(post.ID, 10)).
		BodyJson(postReq).
		Do(context.TODO())

	c.JSONResponse(ctx, model2.ResponseSuccessOne{
		Data: postReq,
	}, fasthttp.StatusCreated)
}

// Delete post with given identifier
func (c PostController) Delete(ctx *fasthttp.RequestCtx) {
	var post model.Post
	c.App.Database.Delete(post.TableName(), "id = $1",
		phi.URLParam(ctx, "postID")).Force()

	c.JSONResponse(ctx, nil, fasthttp.StatusNoContent)
}
