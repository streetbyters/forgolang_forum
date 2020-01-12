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
	"forgolang_forum/database/model"
	model2 "forgolang_forum/model"
	"github.com/fate-lovely/phi"
	"github.com/valyala/fasthttp"
)

// LogoutController user authentication invalidation controller
type LogoutController struct {
	Controller
	*API
}

// Create passphrase invalidation with params
func (c LogoutController) Create(ctx *fasthttp.RequestCtx) {
	userID := phi.URLParam(ctx, "userID")
	passphraseID := phi.URLParam(ctx, "passphraseID")

	var passphrase model.UserPassphrase
	var passphraseInvalidation model.UserPassphraseInvalidation
	c.App.Database.QueryRowWithModel(fmt.Sprintf(`
		SELECT p.* FROM %s AS p
		LEFT OUTER JOIN %s AS pi ON p.id = pi.passphrase_id
		WHERE pi.passphrase_id IS NULL AND p.id = $1 AND p.user_id = $2
	`, passphrase.TableName(),
		passphraseInvalidation.TableName()),
		&passphrase,
		passphraseID,
		userID).Force()

	passphraseInvalidation.PassphraseID = passphrase.ID
	passphraseInvalidation.SourceUserID.SetValid(c.Auth.ID)

	c.App.Database.Insert(new(model.UserPassphraseInvalidation), &passphraseInvalidation,
		"id")

	c.JSONResponse(ctx, model2.ResponseSuccessOne{
		Data: passphraseInvalidation,
	}, fasthttp.StatusCreated)
}
