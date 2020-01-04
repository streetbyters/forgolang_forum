// Copyright 2019 Abdulkadir DILSIZ
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
	"forgolang_forum/database"
	"forgolang_forum/database/model"
	model2 "forgolang_forum/model"
	"github.com/fate-lovely/phi"
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

	var categories []model.Category
	c.App.Database.QueryWithModel(fmt.Sprintf(`
		SELECT c.* FROM %s AS c ORDER BY $1 $1
	`, c.Model.TableName()),
		&categories,
		paginate.OrderField,
		paginate.OrderBy)

	var count int64
	c.App.Database.DB.Get(&count,
		fmt.Sprintf("SELECT count(*) FROM %s", c.Model.TableName()))

	c.JSONResponse(ctx, model2.ResponseSuccess{
		Data:       categories,
		TotalCount: 0,
	}, fasthttp.StatusOK)
}

// Show a category with given identifier
func (c CategoryController) Show(ctx *fasthttp.RequestCtx) {
	var category model.Category
	c.App.Database.QueryRowWithModel(fmt.Sprintf(`
			SELECT c.* FROM %s AS c WHERE id = $1 
		`, c.Model.TableName()),
		&category,
		phi.URLParam(ctx, "categoryID")).Force()

	c.JSONResponse(ctx, model2.ResponseSuccessOne{
		Data: category,
	}, fasthttp.StatusOK)
}

// Create category with valid params
func (c CategoryController) Create(ctx *fasthttp.RequestCtx) {
	var category model.Category
	c.JSONBody(ctx, &category)

	if errs, err := database.ValidateStruct(category); err != nil {
		c.JSONResponse(ctx, model2.ResponseError{
			Errors: errs,
			Detail: fasthttp.StatusMessage(fasthttp.StatusUnprocessableEntity),
		}, fasthttp.StatusUnprocessableEntity)
	}

	err := c.App.Database.Insert(c.Model, category, "id", "inserted_at", "updated_at")
	if errs, err := database.ValidateConstraint(err, category); err != nil {
		c.JSONResponse(ctx, model2.ResponseError{
			Errors: errs,
			Detail: fasthttp.StatusMessage(fasthttp.StatusUnprocessableEntity),
		}, fasthttp.StatusUnprocessableEntity)
	}

	c.JSONResponse(ctx, model2.ResponseSuccessOne{
		Data: category,
	}, fasthttp.StatusCreated)
}

// Update category with given identifier and valid params
func (c CategoryController) Update(ctx *fasthttp.RequestCtx) {
	var category model.Category
	c.App.Database.QueryRowWithModel(fmt.Sprintf(`
			SELECT c.* FROM %s AS c WHERE id = $1
		`, c.Model.TableName()),
		&category,
		phi.URLParam(ctx, "categoryID")).Force()

	var categoryRequest model.Category
	c.JSONBody(ctx, &categoryRequest)

	if errs, err := database.ValidateStruct(category); err != nil {
		c.JSONResponse(ctx, model2.ResponseError{
			Errors: errs,
			Detail: fasthttp.StatusMessage(fasthttp.StatusUnprocessableEntity),
		}, fasthttp.StatusUnprocessableEntity)
	}

	err := c.App.Database.Update(category, categoryRequest, nil, "id", "inserted_at", "updated_at")
	if errs, err := database.ValidateConstraint(err, category); err != nil {
		c.JSONResponse(ctx, model2.ResponseError{
			Errors: errs,
			Detail: fasthttp.StatusMessage(fasthttp.StatusUnprocessableEntity),
		}, fasthttp.StatusUnprocessableEntity)
	}

	c.JSONResponse(ctx, model2.ResponseSuccessOne{
		Data: categoryRequest,
	}, fasthttp.StatusCreated)
}

// Delete category with given identifier
func (c CategoryController) Delete(ctx *fasthttp.RequestCtx) {
	c.App.Database.Delete(fmt.Sprintf("SELECT * FROM %s",
		c.Model.TableName()),
		"id = $1",
		phi.URLParam(ctx, "categoryID")).Force()

	c.JSONResponse(ctx, model2.ResponseSuccessOne{
		Data: nil,
	}, fasthttp.StatusNoContent)
}
