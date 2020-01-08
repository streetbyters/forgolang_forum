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
	"github.com/fate-lovely/phi"
	"github.com/valyala/fasthttp"
	"strconv"
)

// UserPolicy category authorization
type UserPolicy struct {
	Policy
	*API
}

// Index method for user api authorization
func (p UserPolicy) Index(next phi.HandlerFunc) phi.HandlerFunc {
	return p.API.Authorization.Apply(next, "UserController", "Index",
		func(ctx *fasthttp.RequestCtx) bool {
			return true
		})
}

// Show method for user api authorization
func (p UserPolicy) Show(next phi.HandlerFunc) phi.HandlerFunc {
	return p.API.Authorization.Apply(next, "UserController", "Show",
		func(ctx *fasthttp.RequestCtx) bool {
			if i, err := strconv.ParseInt(phi.URLParam(ctx, "userID"), 10, 64); err == nil && i == p.Auth.ID {
				return true
			}
			return false
		})
}

// Create method for user api authorization
func (p UserPolicy) Create(next phi.HandlerFunc) phi.HandlerFunc {
	return p.API.Authorization.Apply(next, "UserController", "Create",
		func(ctx *fasthttp.RequestCtx) bool {
			return true
		})
}

// Update method for user api authorization
func (p UserPolicy) Update(next phi.HandlerFunc) phi.HandlerFunc {
	return p.API.Authorization.Apply(next, "UserController", "Update",
		func(ctx *fasthttp.RequestCtx) bool {
			if i, err := strconv.ParseInt(phi.URLParam(ctx, "userID"), 10, 64); err == nil && i == p.Auth.ID {
				return true
			}
			return false
		})
}

// Delete method for user api authorization
func (p UserPolicy) Delete(next phi.HandlerFunc) phi.HandlerFunc {
	return p.API.Authorization.Apply(next, "UserController", "Delete",
		func(ctx *fasthttp.RequestCtx) bool {
			return true
		})
}
