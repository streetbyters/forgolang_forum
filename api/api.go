// Copyright 2019 StreetByters Community
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
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"forgolang_forum/cmn"
	"forgolang_forum/database"
	model2 "forgolang_forum/database/model"
	"forgolang_forum/model"
	"forgolang_forum/utils"
	"github.com/go-redis/redis"
	"github.com/olivere/elastic/v7"
	pluggableError "github.com/streetbyters/agente/errors"
	"github.com/valyala/fasthttp"
	"net/url"
	"strconv"
)

// API rest api structure
type API struct {
	App           *cmn.App
	Router        *Router
	JWTAuth       *JWTAuth
	Authorization *Authorization
	Languages     []model2.Language
}

// NewAPI building api
func NewAPI(app *cmn.App) *API {
	api := &API{App: app}
	api.JWTAuth = NewJWTAuth(api)
	api.Authorization = NewAuthorization(api)
	api.Router = NewRouter(api)

	var language model2.Language
	var languages []model2.Language
	result := api.GetDB().QueryWithModel(fmt.Sprintf(`
		SELECT * FROM %s AS l ORDER BY l.id ASC
	`, language.TableName()),
		&languages)
	cmn.FailOnError(api.App.Logger, result.Error)
	api.Languages = languages

	return api
}

// ParseQuery parse url query string
func (a *API) ParseQuery(ctx *fasthttp.RequestCtx) map[string]string {
	qs, _ := url.ParseQuery(string(ctx.URI().QueryString()))
	values := make(map[string]string)
	for key, val := range qs {
		values[key] = val[0]
	}

	return values
}

// Paginate request paginate build
func (a *API) Paginate(ctx *fasthttp.RequestCtx, orderFields ...string) (model.Pagination, map[string]string, error) {
	var err error
	errs := make(map[string]string)
	pagination := model.NewPagination()
	queryParams := a.ParseQuery(ctx)

	if val, ok := queryParams["limit"]; ok {
		pagination.Limit, err = strconv.Atoi(val)
		if err != nil {
			errs["limit"] = "is not valid"
		}
	}

	if val, ok := queryParams["offset"]; ok {
		pagination.Offset, err = strconv.ParseInt(val, 10, 64)
		if err != nil {
			errs["offset"] = "is not valid"
		}
	}

	if val, ok := queryParams["order_by"]; ok {
		pagination.OrderBy = val
	}

	if val, ok := queryParams["order_field"]; ok {
		pagination.OrderField = val
		if exists, _ := utils.InArray(val, orderFields); !exists {
			err = errors.New("order field is not valid")
			errs["order_field"] = "is not valid"
		}
	}

	if err != nil {
		panic(pluggableError.New(fasthttp.StatusMessage(fasthttp.StatusBadRequest),
			fasthttp.StatusBadRequest,
			"paginate params not valid",
			errs))
	}

	errs, err = pagination.Validate(orderFields...)

	if err != nil {
		panic(pluggableError.New(fasthttp.StatusMessage(fasthttp.StatusBadRequest),
			fasthttp.StatusBadRequest,
			"paginate params not valid",
			errs))
	}

	return pagination, errs, err
}

// GetAuthContext get auth context request
func (a *API) GetAuthContext(ctx *fasthttp.RequestCtx) *model.AuthContext {
	return ctx.UserValue("AuthContext").(*model.AuthContext)
}

// GetLanguageContext get default language context
func (a *API) GetLanguageContext(ctx *fasthttp.RequestCtx) *model2.Language {
	return ctx.UserValue("Language").(*model2.Language)
}

func (a *API) SetLanguageContext(code string, ctx *fasthttp.RequestCtx) *fasthttp.RequestCtx {
	if l := a.GetLanguage(code); l != nil {
		ctx.SetUserValue("Language", l)
		return ctx
	}

	ctx.SetUserValue("Language", a.GetLanguage(a.App.Config.Lang))
	return ctx
}

// GetDB api database getter
func (a *API) GetDB() *database.Database {
	return a.App.Database
}

// GetLanguage api language getter with given code
func (a *API) GetLanguage(code string) *model2.Language {
	exists := false
	language := model2.NewLanguage()
	for _, l := range a.Languages {
		if l.Code == code {
			language = &l
			exists = !exists
			break
		}
	}

	if exists {
		return language
	}

	return nil
}

// GetElastic api elastic search getter
func (a *API) GetElastic() *elastic.Client {
	return a.App.ElasticClient
}

// GetCache api redis getter
func (a *API) GetCache() *redis.Client {
	return a.App.Cache
}

// JSONBody parse given model request body
func (a *API) JSONBody(ctx *fasthttp.RequestCtx, model interface{}) {
	r := bytes.NewReader(ctx.PostBody())
	json.NewDecoder(r).Decode(&model)
}

// JSONResponse building json response
func (a *API) JSONResponse(ctx *fasthttp.RequestCtx, response model.ResponseInterface, status int) {
	ctx.Response.Header.Set("Content-Type", "application/json; charset=utf-8")
	if response != nil {
		ctx.SetBody([]byte(response.ToJSON()))
	}
	ctx.SetStatusCode(status)
}
