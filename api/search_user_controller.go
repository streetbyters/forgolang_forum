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
	"context"
	"encoding/json"
	"fmt"
	"forgolang_forum/database/model"
	model2 "forgolang_forum/model"
	"forgolang_forum/utils"
	"github.com/olivere/elastic/v7"
	"github.com/valyala/fasthttp"
	"strconv"
)

// SearchUserController user search api controller
type SearchUserController struct {
	Controller
	*API
}

func (c SearchUserController) Index(ctx *fasthttp.RequestCtx) {
	paginate, _, _ := c.Paginate(ctx,
		"_score", "id", "state", "inserted_at", "updated_at")
	queryParams := c.ParseQuery(ctx)

	searchQuery := utils.ToSearchString(queryParams["query"])

	var count int64
	var users []model.User
	var ascending bool = true

	if paginate.OrderBy == "desc" {
		ascending = false
	}

	q := elastic.NewBoolQuery()

	if b, _ := strconv.ParseBool(queryParams["multi"]); c.GetAuthContext(ctx).Role == "superadmin" && b {
		if b, _ := strconv.ParseBool(queryParams["is_active"]); b {
			q = q.Must(elastic.NewTermQuery("is_active", true))
		}
		if queryParams["username"] != "" {
			q = q.Must(elastic.NewMatchQuery("username", fmt.Sprintf("*%s*", utils.ToSearchString(queryParams["username"]))))
		}
		if queryParams["email"] != "" {
			q = q.Must(elastic.NewMatchQuery("email", fmt.Sprintf("*%s*", utils.ToSearchString(queryParams["email"]))))
		}
	} else {
		q = q.Must(elastic.NewQueryStringQuery(fmt.Sprintf("*%s*", searchQuery)))
	}

	query := elastic.NewConstantScoreQuery(q)

	results, _ := c.App.ElasticClient.Search("users").
		Query(query).
		Sort(paginate.OrderField, ascending).
		From(int(paginate.Offset)).
		Size(paginate.Limit).
		Do(context.TODO())

	count, _ = c.App.ElasticClient.Count("users").
		Query(query).
		Do(context.TODO())

	for _, v := range results.Hits.Hits {
		var u model.User
		b, _ := v.Source.MarshalJSON()
		json.Unmarshal(b, &u)
		if c.GetAuthContext(ctx).Role != "superadmin" && u.EmailHidden {
			u.Email = "***"
		}
		users = append(users, u)
	}

	c.JSONResponse(ctx, model2.ResponseSuccess{
		Data:       users,
		TotalCount: count,
	}, fasthttp.StatusOK)
}
