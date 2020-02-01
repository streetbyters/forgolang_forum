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
	"encoding/json"
	"fmt"
	"forgolang_forum/cmn"
	"forgolang_forum/database"
	"forgolang_forum/database/model"
	model2 "forgolang_forum/model"
	"github.com/fate-lovely/phi"
	"github.com/gosimple/slug"
	"github.com/jmoiron/sqlx"
	"github.com/valyala/fasthttp"
)

// CategoryController Users discussion topic categories api structure
type CategoryController struct {
	Controller
	*API
	Model model.Category
}

// Index list all categories
func (c CategoryController) Index(ctx *fasthttp.RequestCtx) {
	paginate, _, _ := c.Paginate(ctx, "id", "inserted_at")

	var count int64
	var categories []model.Category

	if s, err := c.App.Cache.SMembers(cmn.GetRedisKey("category", "all")).Result(); err == nil && len(s) > 0 {
		var keys []string
		for _, v := range s {
			var c model.Category
			json.Unmarshal([]byte(v), &c)
			keys = append(keys, fmt.Sprintf("%s:%d",
				cmn.GetRedisKey("category", "languages"),
				c.ID))
			categories = append(categories, c)
		}

		categoryLanguages := make(map[int64][]model.CategoryLanguage)
		if len(keys) > 0 {
			s, _ = c.GetCache().SUnion(keys...).Result()
			for _, cachedLang := range s {
				var cl model.CategoryLanguage
				json.Unmarshal([]byte(cachedLang), &cl)
				categoryLanguages[cl.CategoryID] = append(categoryLanguages[cl.CategoryID], cl)
			}
		}

		var categories2 []model.Category
		for _, c := range categories {
			c.Languages = categoryLanguages[c.ID]
			categories2 = append(categories2, c)
		}

		count, _ := c.App.Cache.SCard(cmn.GetRedisKey("category", "all")).Result()

		c.JSONResponse(ctx, model2.ResponseSuccess{
			Data:       categories2,
			TotalCount: count,
		}, fasthttp.StatusOK)
		return
	}

	c.GetDB().QueryWithModel(fmt.Sprintf(`
		SELECT c.* FROM %s AS c ORDER BY %s %s
	`, c.Model.TableName(),
		paginate.OrderField,
		paginate.OrderBy),
		&categories)

	c.GetDB().DB.Get(&count,
		fmt.Sprintf("SELECT count(*) FROM %s", c.Model.TableName()))

	var ids []int64
	var cats []interface{}
	for _, ca := range categories {
		cats = append(cats, ca.ToJSON())
		ids = append(ids, ca.ID)
	}
	c.App.Cache.SAdd(cmn.GetRedisKey("category", "all"), cats...)

	var categoryLanguage model.CategoryLanguage
	var categoryLanguages []model.CategoryLanguage
	query, args, _ := sqlx.In(fmt.Sprintf(`
		SELECT cl.* FROM %s AS cl
		LEFT OUTER JOIN %s AS cl2 ON 
			cl.language_id = cl2.language_id AND
			cl.category_id = cl2.category_id AND
			cl.id < cl2.id
		WHERE cl2.id IS NULL AND cl.category_id IN (?)
	`, categoryLanguage.TableName(), categoryLanguage.TableName()), ids)
	query = c.GetDB().DB.Rebind(query)
	c.GetDB().QueryWithModel(query,
		&categoryLanguages,
		args...)

	categoryLanguagesMap := make(map[int64][]model.CategoryLanguage)
	for _, ca := range categoryLanguages {
		categoryLanguagesMap[ca.CategoryID] = append(categoryLanguagesMap[ca.CategoryID], ca)
	}

	var categories2 []model.Category
	for _, ca := range categories {
		ca.Languages = categoryLanguagesMap[ca.ID]
		categories2 = append(categories2, ca)
	}

	c.JSONResponse(ctx, model2.ResponseSuccess{
		Data:       categories2,
		TotalCount: count,
	}, fasthttp.StatusOK)
}

// Show a category with given identifier
func (c CategoryController) Show(ctx *fasthttp.RequestCtx) {
	var category model.Category

	var cs string
	err := c.App.Cache.Get(fmt.Sprintf("%s:%s",
		cmn.GetRedisKey("category", "one"),
		phi.URLParam(ctx, "categoryID"))).Scan(&cs)
	err2 := c.App.Cache.Get(fmt.Sprintf("%s:%s",
		cmn.GetRedisKey("category", "slug"),
		phi.URLParam(ctx, "categoryID"))).Scan(&cs)
	if err == nil || err2 == nil {
		json.Unmarshal([]byte(cs), &category)

		var categoryLanguages []model.CategoryLanguage
		s, _ := c.GetCache().SMembers(fmt.Sprintf("%s:%d",
			cmn.GetRedisKey("category", "languages"),
			category.ID)).Result()
		for _, cl := range s {
			var l model.CategoryLanguage
			json.Unmarshal([]byte(cl), &l)
			categoryLanguages = append(categoryLanguages, l)
		}

		category.Languages = categoryLanguages

		c.JSONResponse(ctx, model2.ResponseSuccessOne{
			Data: category,
		}, fasthttp.StatusOK)
		return
	}

	c.GetDB().QueryRowWithModel(fmt.Sprintf(`
			SELECT c.* FROM %s AS c WHERE id::text = $1::text OR slug::text = $1::text
		`, c.Model.TableName()),
		&category,
		phi.URLParam(ctx, "categoryID")).Force()

	c.App.Cache.Set(fmt.Sprintf("%s:%d",
		cmn.GetRedisKey("category", "one"),
		category.ID),
		category.ToJSON(),
		0)
	c.App.Cache.Set(fmt.Sprintf("%s:%s",
		cmn.GetRedisKey("category", "slug"),
		category.Slug),
		category.ToJSON(),
		0)

	var categoryLanguage model.CategoryLanguage
	var categoryLanguages []model.CategoryLanguage
	c.GetDB().QueryWithModel(fmt.Sprintf(`
		SELECT cl.* FROM %s AS cl
		LEFT OUTER JOIN %s AS cl2 ON 
			cl.language_id = cl2.language_id AND
			cl.category_id = cl2.category_id AND
			cl.id < cl2.id
		WHERE cl2.id IS NULL AND cl.category_id = $1
	`, categoryLanguage.TableName(), categoryLanguage.TableName()),
		&categoryLanguages,
		category.ID)
	category.Languages = categoryLanguages

	c.JSONResponse(ctx, model2.ResponseSuccessOne{
		Data: category,
	}, fasthttp.StatusOK)
}

// Create category with valid params
func (c CategoryController) Create(ctx *fasthttp.RequestCtx) {
	category := new(model.Category)
	c.JSONBody(ctx, &category)

	if errs, err := database.ValidateStruct(category); err != nil {
		c.JSONResponse(ctx, model2.ResponseError{
			Errors: errs,
			Detail: fasthttp.StatusMessage(fasthttp.StatusUnprocessableEntity),
		}, fasthttp.StatusUnprocessableEntity)
		return
	}

	category.Slug = slug.Make(category.Title)

	err := c.GetDB().Insert(new(model.Category), category, "id", "inserted_at", "updated_at")
	if errs, err := database.ValidateConstraint(err, category); err != nil {
		c.JSONResponse(ctx, model2.ResponseError{
			Errors: errs,
			Detail: fasthttp.StatusMessage(fasthttp.StatusUnprocessableEntity),
		}, fasthttp.StatusUnprocessableEntity)
		return
	}

	c.App.Cache.Set(fmt.Sprintf("%s:%d",
		cmn.GetRedisKey("category", "one"),
		category.ID), category.ToJSON(), 0).Err()
	c.App.Cache.SAdd(cmn.GetRedisKey("category", "all"), category.ToJSON())
	c.App.Cache.Set(fmt.Sprintf("%s:%s",
		cmn.GetRedisKey("category", "slug"),
		category.Slug), category.ToJSON(), 0)

	c.JSONResponse(ctx, model2.ResponseSuccessOne{
		Data: category,
	}, fasthttp.StatusCreated)
}

// Update category with given identifier and valid params
func (c CategoryController) Update(ctx *fasthttp.RequestCtx) {
	category := new(model.Category)
	c.GetDB().QueryRowWithModel(fmt.Sprintf(`
			SELECT c.* FROM %s AS c WHERE c.id = $1
		`, c.Model.TableName()),
		category,
		phi.URLParam(ctx, "categoryID")).Force()

	var categoryRequest model.Category
	c.JSONBody(ctx, &categoryRequest)

	if errs, err := database.ValidateStruct(categoryRequest); err != nil {
		c.JSONResponse(ctx, model2.ResponseError{
			Errors: errs,
			Detail: fasthttp.StatusMessage(fasthttp.StatusUnprocessableEntity),
		}, fasthttp.StatusUnprocessableEntity)
		return
	}

	categoryRequest.Slug = slug.Make(categoryRequest.Title)
	categoryRequest.InsertedAt = category.InsertedAt

	c.App.Cache.SRem(cmn.GetRedisKey("category", "all"), category.ToJSON())

	err := c.GetDB().Update(category, &categoryRequest, nil,
		"id", "inserted_at", "updated_at")
	if errs, err := database.ValidateConstraint(err, category); err != nil {
		c.App.Cache.SAdd(cmn.GetRedisKey("category", "all"), category.ToJSON())
		c.JSONResponse(ctx, model2.ResponseError{
			Errors: errs,
			Detail: fasthttp.StatusMessage(fasthttp.StatusUnprocessableEntity),
		}, fasthttp.StatusUnprocessableEntity)
		return
	}

	c.App.Cache.Set(fmt.Sprintf("%s:%d",
		cmn.GetRedisKey("category", "one"),
		category.ID), category.ToJSON(), 0)
	c.App.Cache.Set(fmt.Sprintf("%s:%s",
		cmn.GetRedisKey("category", "slug"),
		category.Slug), category.ToJSON(), 0)

	c.App.Cache.SAdd(cmn.GetRedisKey("category", "all"), category.ToJSON())

	c.JSONResponse(ctx, model2.ResponseSuccessOne{
		Data: category,
	}, fasthttp.StatusOK)
}

// Delete category with given identifier
func (c CategoryController) Delete(ctx *fasthttp.RequestCtx) {
	var category model.Category
	c.GetDB().QueryRowWithModel(fmt.Sprintf("SELECT c.* FROM %s AS c WHERE c.id = $1",
		category.TableName()),
		&category,
		phi.URLParam(ctx, "categoryID")).Force()

	c.GetDB().Delete(c.Model.TableName(),
		"id = $1",
		phi.URLParam(ctx, "categoryID")).Force()

	c.App.Cache.SRem(cmn.GetRedisKey("category", "all"), category.ToJSON())
	c.App.Cache.Del(fmt.Sprintf("%s:%s",
		cmn.GetRedisKey("category", "one"),
		phi.URLParam(ctx, "categoryID")),
		fmt.Sprintf("%s:%s",
			cmn.GetRedisKey("category", "slug"),
			category.Slug))

	c.JSONResponse(ctx, model2.ResponseSuccessOne{
		Data: nil,
	}, fasthttp.StatusNoContent)
}
