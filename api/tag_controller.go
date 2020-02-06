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
	"encoding/json"
	"fmt"
	"forgolang_forum/cmn"
	"forgolang_forum/database/model"
	model2 "forgolang_forum/model"
	"github.com/valyala/fasthttp"
)

// TagController special classifications api controller
type TagController struct {
	Controller
	*API
	Model model.Tag
}

func (c TagController) Index(ctx *fasthttp.RequestCtx) {
	paginate, _, _ := c.Paginate(ctx, "count", "name")

	var count int64
	var tags []model.Tag

	if s, err := c.App.Cache.SMembers(cmn.GetRedisKey("tag", "all")).Result(); err == nil {
		var ts []string
		for _, v := range s {
			var t model.Tag
			json.Unmarshal([]byte(v), &t)
			tags = append(tags, t)
			ts = append(ts, fmt.Sprintf("%s:%d", cmn.GetRedisKey("tag", "count"), t.ID))
		}

		//count, _ := c.App.Cache.SCard(cmn.GetRedisKey("tag", "all")).Result()
		//
		//tss, _ := c.App.Cache.MGet(ts...).Result()
		//for _, tc := range tss {
		//	fmt.Println(tc)
		//}

		c.JSONResponse(ctx, model2.ResponseSuccess{
			Data:       tags,
			TotalCount: count,
		}, fasthttp.StatusOK)
		return
	}

	c.GetDB().QueryWithModel(fmt.Sprintf(`
		SELECT t.* FROM %s AS t
		ORDER BY %s %s
	`, c.Model.TableName(), paginate.OrderField, paginate.OrderBy),
		&tags)

	c.GetDB().DB.Get(&count, fmt.Sprintf("SELECT count(*) FROM %s", c.Model.TableName()))

	var ts []interface{}
	for _, t := range tags {
		ts = append(ts, t.ToJSON())
	}

	c.App.Cache.SAdd(cmn.GetRedisKey("tags", "all"), ts...)

	c.JSONResponse(ctx, model2.ResponseSuccess{
		Data:       tags,
		TotalCount: count,
	}, fasthttp.StatusOK)
}