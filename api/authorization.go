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
	"forgolang_forum/database/model"
	pluggableError "github.com/akdilsiz/agente/errors"
	"github.com/fate-lovely/phi"
	"github.com/valyala/fasthttp"
)

// Authorization middleware
type Authorization struct {
	*API
}

// NewAuthorization generate middleware
func NewAuthorization(api *API) *Authorization {
	return &Authorization{API: api}
}

// Apply module authorization
func (m *Authorization) Apply(next phi.HandlerFunc, controller, method string, cb func(ctx *fasthttp.RequestCtx) bool) phi.HandlerFunc {
	return func(ctx *fasthttp.RequestCtx) {
		if m.Auth.Role == "superadmin" {
			next(ctx)
			return
		}

		if !m.gen(controller, method) || !cb(ctx) {
			panic(pluggableError.New("forbidden",
				fasthttp.StatusForbidden,
				fasthttp.StatusMessage(fasthttp.StatusForbidden)))
		}

		next(ctx)
	}
}

func (m *Authorization) gen(controller, method string) bool {
	rK := cmn.RedisKeys["user"]
	rK = rK.(map[string]string)["permission"]

	role := new(model.Role)
	rolePermission := new(model.RolePermission)
	roleAssignment := new(model.UserRoleAssignment)
	roleAssignmentInvalidation := new(model.UserRoleAssignmentInvalidation)
	_, err := m.App.Cache.Get(fmt.Sprintf("%s:%s:%d:%s:%s",
		rK,
		m.API.Auth.Role,
		m.API.Auth.ID,
		controller,
		method)).Result()
	if err != nil {
		err := m.App.Database.DB.QueryRow(fmt.Sprintf("SELECT r.id, r.code FROM %s AS ra "+
			"LEFT OUTER JOIN %s AS ra2 ON ra.user_id = ra2.user_id and ra.id < ra2.id "+
			"LEFT OUTER JOIN %s AS rai ON ra.id = rai.assignment_id "+
			"INNER JOIN %s AS r ON ra.role_id = r.id "+
			"INNER JOIN %s AS rp ON rp.role_id = r.id "+
			"WHERE ra2.id IS NULL AND rai.assignment_id IS NULL AND "+
			"rp.controller = $1 AND rp.method = $2 AND ra.user_id = $3 AND ra.role_id = $4",
			roleAssignment.TableName(),
			roleAssignment.TableName(),
			roleAssignmentInvalidation.TableName(),
			role.TableName(),
			rolePermission.TableName()),
			controller,
			method,
			m.API.Auth.ID,
			m.API.Auth.RoleID).Scan(&role.ID, &role.Code)
		if err != nil {
			return false
		}

		defer func() {
			m.App.Cache.Set(fmt.Sprintf("%s:%s:%d:%s:%s",
				rK,
				m.API.Auth.Role,
				m.API.Auth.ID,
				controller,
				method),
				true,
				0)
		}()

		return true
	}

	return true
}
