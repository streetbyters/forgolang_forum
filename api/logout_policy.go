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
	"forgolang_forum/utils"
	"github.com/fate-lovely/phi"
	"github.com/valyala/fasthttp"
)

// LogoutPolicy logout authorization
type LogoutPolicy struct {
	Policy
	*API
}

// Create method for logout api authorization
func (p LogoutPolicy) Create(next phi.HandlerFunc) phi.HandlerFunc {
	return p.API.Authorization.Apply(next, "LogoutController", "Create",
		func(ctx *fasthttp.RequestCtx) bool {
			userID, _ := utils.ParseInt(phi.URLParam(ctx, "userID"), 10, 64)

			if userID == p.Auth.ID {
				return true
			}

			return false
		})
}
