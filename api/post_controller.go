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
	"fmt"
	"forgolang_forum/cmn"
	"forgolang_forum/database"
	"forgolang_forum/database/model"
	model2 "forgolang_forum/model"
	"github.com/fate-lovely/phi"
	"github.com/gosimple/slug"
	"github.com/valyala/fasthttp"
)

// PostController discussions api controller
type PostController struct {
	Controller
	*API
	Model model.Post
}

// Index list all discussions with filter params
func (c PostController) Index(ctx *fasthttp.RequestCtx) {
	paginate, _, _ := c.Paginate(ctx, "id", "inserted_at", "updated_at")

	var posts []model.PostDEP
	var postSlug model.PostSlug
	var postDetail model.PostDetail
	var user model.User
	c.App.Database.QueryWithModel(fmt.Sprintf(`
		SELECT 
			p.id as id, p.author_id as author_id, u.username as author_username, 
			p.inserted_at as inserted_at, ps.slug as slug, pd.title as title, 
			pd.description as description, pd.content as content
		FROM %s AS p
		LEFT OUTER JOIN %s AS ps ON p.id = ps.post_id
		LEFT OUTER JOIN %s AS ps2 ON ps.post_id = ps2.post_id AND ps.id < ps2.id
		INNER JOIN %s AS pd ON p.id = pd.post_id
		LEFT OUTER JOIN %s AS pd2 ON pd.post_id = pd2.post_id AND pd.id < pd2.id
		INNER JOIN %s AS u ON p.author_id = u.id
		WHERE ps2.id IS NULL AND pd2.id IS NULL
		ORDER BY %s %s
		LIMIt $1 OFFSET $2
	`, c.Model.TableName(), postSlug.TableName(), postSlug.TableName(), postDetail.TableName(),
		postDetail.TableName(), user.TableName(),
		paginate.OrderField,
		paginate.OrderBy),
		&posts,
		paginate.Limit,
		paginate.Offset)

	var count int64
	count, _ = c.App.Cache.Get(cmn.GetRedisKey("post", "count")).Int64()
	if count == 0 {
		c.App.Database.DB.Get(&count, fmt.Sprintf(`
			SELECT count(p.id) FROM %s AS p
			INNER JOIN %s AS pd ON p.id = pd.post_id
			LEFT OUTER JOIN %s AS pd2 ON pd.post_id = pd2.post_id AND pd.id < pd2.id
			WHERE pd2.id IS NULL
		`, c.Model.TableName(), postDetail.TableName(), postDetail.TableName()))
	}

	c.JSONResponse(ctx, model2.ResponseSuccess{
		Data:       posts,
		TotalCount: count,
	}, fasthttp.StatusOK)
}

// Show discussion with given identifier or slug
func (c PostController) Show(ctx *fasthttp.RequestCtx) {
	var post model.PostDEP
	var postSlug model.PostSlug
	var postDetail model.PostDetail
	var user model.User
	c.App.Database.QueryRowWithModel(fmt.Sprintf(`
		SELECT 
			p.id as id, p.author_id as author_id, u.username as author_username, 
			p.inserted_at as inserted_at, ps.slug as slug, pd.title as title, 
			pd.description as description, pd.content as content
		FROM %s AS p
		LEFT OUTER JOIN %s AS ps ON p.id = ps.post_id
		LEFT OUTER JOIN %s AS ps2 ON ps.post_id = ps2.post_id AND ps.id < ps2.id
		INNER JOIN %s AS pd ON p.id = pd.post_id
		LEFT OUTER JOIN %s AS pd2 ON pd.post_id = pd2.post_id AND pd.id < pd2.id
		INNER JOIN %s AS u ON p.author_id = u.id
		WHERE ps2.id IS NULL AND pd2.id IS NULL AND (p.id::text = $1::text OR ps.slug = $1)
	`, c.Model.TableName(), postSlug.TableName(), postSlug.TableName(), postDetail.TableName(),
		postDetail.TableName(), user.TableName()),
		&post,
		phi.URLParam(ctx, "postID")).Force()

	c.JSONResponse(ctx, model2.ResponseSuccessOne{
		Data: post,
	}, fasthttp.StatusOK)
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
	var errs map[string]string
	c.App.Database.Transaction(func(tx *database.Tx) error {
		post.AuthorID = c.GetAuthContext(ctx).ID
		err = tx.DB.Insert(new(model.Post), post, "id")
		if errs, err = database.ValidateConstraint(err, post); err != nil {
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
