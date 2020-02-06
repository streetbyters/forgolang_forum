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
	"fmt"
	"forgolang_forum/database/model"
	"github.com/fate-lovely/phi"
	"github.com/valyala/fasthttp"
)

// PostPolicy post authorization
type PostPolicy struct {
	Policy
	*API
}

// Create method for posts api authorization
func (p PostPolicy) Create(next phi.HandlerFunc) phi.HandlerFunc {
	return p.API.Authorization.Apply(next, "PostController", "Create",
		func(ctx *fasthttp.RequestCtx) bool {
			return true
		})
}

// Delete method for posts api authorization
func (p PostPolicy) Delete(next phi.HandlerFunc) phi.HandlerFunc {
	return p.API.Authorization.Apply(next, "PostController", "Delete",
		func(ctx *fasthttp.RequestCtx) bool {
			if post := p.GetPost(ctx); post != nil && post.AuthorID == p.GetAuthContext(ctx).ID {
				return true
			}
			return false
		})
}

func (p PostPolicy) GetPost(ctx *fasthttp.RequestCtx) *model.Post {
	var post model.Post
	var postSlug model.PostSlug

	p.App.Database.QueryRowWithModel(fmt.Sprintf(`
		SELECT 
			p.id, p.author_id, p.inserted_at
		FROM %s AS p
		LEFT OUTER JOIN %s AS ps ON p.id = ps.post_id
		LEFT OUTER JOIN %s AS ps2 ON ps.post_id = ps2.post_id AND ps.id < ps2.id
		WHERE ps2.id IS NULL AND (p.id::text = $1::text OR ps.slug = $1)
	`, post.TableName(), postSlug.TableName(), postSlug.TableName()),
		&post,
		phi.URLParam(ctx, "postID"))

	return &post
}
