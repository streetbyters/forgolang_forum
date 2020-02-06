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
	"errors"
	"fmt"
	"forgolang_forum/cmn"
	"forgolang_forum/database"
	model2 "forgolang_forum/database/model"
	"forgolang_forum/model"
	"github.com/fate-lovely/phi"
	"github.com/gosimple/slug"
	"github.com/valyala/fasthttp"
	"strconv"
)

// CategoryLanguageController category localization(L10N) api structure
type CategoryLanguageController struct {
	Controller
	*API
}

// Create category language with valid params
func (c CategoryLanguageController) Create(ctx *fasthttp.RequestCtx) {
	var err error
	categoryID, err := strconv.ParseInt(phi.URLParam(ctx, "categoryID"), 10, 64)
	if err != nil {
		c.JSONResponse(ctx, model.ResponseError{
			Detail: fasthttp.StatusMessage(fasthttp.StatusBadRequest),
		}, fasthttp.StatusBadRequest)
		return
	}

	categoryLanguage := new(model2.CategoryLanguage)
	c.JSONBody(ctx, &categoryLanguage)
	categoryLanguage.CategoryID = categoryID
	categoryLanguage.SourceUserID.SetValid(c.GetAuthContext(ctx).ID)
	if errs, err := database.ValidateStruct(categoryLanguage); err != nil {
		c.JSONResponse(ctx, model.ResponseError{
			Errors: errs,
			Detail: fasthttp.StatusMessage(fasthttp.StatusUnprocessableEntity),
		}, fasthttp.StatusUnprocessableEntity)
		return
	}

	categoryLanguage.Slug = slug.Make(categoryLanguage.Title)

	errs := make(map[string]string)
	c.GetDB().Transaction(func(tx *database.Tx) error {
		catLang := new(model2.CategoryLanguage)
		result := tx.DB.QueryRow(fmt.Sprintf(`
			SELECT cl.* FROM %s AS cl
			LEFT OUTER JOIN %s AS cl2 ON cl.language_id = cl2.language_id AND cl.id < cl2.id
			WHERE cl2.id IS NULL AND cl.slug = $1 AND cl.language_id = $2
		`, catLang.TableName(), catLang.TableName()),
			categoryLanguage.Slug,
			categoryLanguage.LanguageID)
		if result.Count > 0 {
			err = errors.New("slug has been already taken")
			errs["slug"] = "has been already taken"
			return err
		}

		err = c.GetDB().Insert(new(model2.CategoryLanguage), categoryLanguage,
			"id", "inserted_at")
		if errs, err = database.ValidateConstraint(err, categoryLanguage); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		c.JSONResponse(ctx, model.ResponseError{
			Errors: errs,
			Detail: fasthttp.StatusMessage(fasthttp.StatusUnprocessableEntity),
		}, fasthttp.StatusUnprocessableEntity)
		return
	}

	c.GetCache().SAdd(fmt.Sprintf("%s:%d",
		cmn.GetRedisKey("category", "languages"),
		categoryID), categoryLanguage.ToJSON())

	c.JSONResponse(ctx, model.ResponseSuccessOne{
		Data: categoryLanguage,
	}, fasthttp.StatusCreated)
}
