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
	"context"
	"fmt"
	"forgolang_forum/cmn"
	"forgolang_forum/database"
	"forgolang_forum/database/model"
	model2 "forgolang_forum/model"
	"github.com/fate-lovely/phi"
	"github.com/valyala/fasthttp"
	"strconv"
	"time"
)

// ConfirmationController user activation api controller
type ConfirmationController struct {
	Controller
	*API
}

// Create method for user activation
func (c ConfirmationController) Create(ctx *fasthttp.RequestCtx) {
	otc := new(model.UserOneTimeCode)
	userState := new(model.UserState)

	c.App.Database.QueryRowWithModel(fmt.Sprintf(`
		SELECT otc.* FROM %s AS otc
		LEFT OUTER JOIN %s AS us ON otc.user_id = us.user_id
		LEFT OUTER JOIN %s AS us2 ON us.user_id = us2.user_id and us.id < us2.id
		WHERE us2.id IS NULL AND otc.user_id = $1 AND otc.code = $2 AND otc.type = 'confirmation' AND
			otc.inserted_at >= ((CURRENT_TIMESTAMP at time zone 'utc') - interval '15 minutes') AND 
			(us.id IS NULL OR us.state != 'active')
	`, otc.TableName(), userState.TableName(), userState.TableName()),
		otc,
		phi.URLParam(ctx, "userID"),
		phi.URLParam(ctx, "code")).Force()

	userState = model.NewUserState(otc.UserID)
	userState.State = database.Active
	userState.SourceUserID.SetValid(otc.UserID)
	c.App.Database.Insert(new(model.UserState), userState, "id")

	user := new(model.User)
	c.App.Database.QueryRowWithModel(fmt.Sprintf("%s AND u.id = $1", user.Query(false)),
		&user,
		phi.URLParam(ctx, "userID")).Force()

	c.App.Cache.Set(fmt.Sprintf("%s:%d",
		cmn.GetRedisKey("user", "one"), user.ID),
		user.ToJSON(), time.Minute*30)

	c.App.ElasticClient.Index().Index("users").
		Id(strconv.FormatInt(user.ID, 10)).
		BodyJson(user).
		Do(context.TODO())

	go func() {
		c.App.Queue.Email.Publish(cmn.QueueEmailBody{
			Subject:    "Welcome to Forgolang.com",
			Recipients: []string{user.Email},
			Type:       "welcome",
			Template:   "welcome",
			Params: struct {
				UserName string
			}{
				UserName: user.Username,
			},
		}.ToJSON())
	}()

	c.JSONResponse(ctx, model2.ResponseSuccessOne{
		Data: otc,
	}, fasthttp.StatusCreated)
}
