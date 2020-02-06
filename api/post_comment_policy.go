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
	"fmt"
	"forgolang_forum/database/model"
	"github.com/fate-lovely/phi"
	"github.com/valyala/fasthttp"
)

// PostCommentPolicy post comment authorization
type PostCommentPolicy struct {
	Policy
	*API
}

// Create post comment authorization
func (p PostCommentPolicy) Create(next phi.HandlerFunc) phi.HandlerFunc {
	return p.API.Authorization.Apply(next, "PostCommentController", "Create",
		func(ctx *fasthttp.RequestCtx) bool {
			return true
		})
}

// Delete post comment authorization
func (p PostCommentPolicy) Delete(next phi.HandlerFunc) phi.HandlerFunc {
	return p.API.Authorization.Apply(next, "PostCommentController", "Delete",
		func(ctx *fasthttp.RequestCtx) bool {
			if comment := p.GetComment(ctx); comment != nil && comment.UserID == p.GetAuthContext(ctx).ID {
				return true
			}
			return false
		})
}

// GetComment get comment
func (p PostCommentPolicy) GetComment(ctx *fasthttp.RequestCtx) *model.PostComment {
	postComment := new(model.PostComment)

	p.GetDB().QueryRowWithModel(fmt.Sprintf(`
		SELECT c.* FROM %s AS c
		WHERE c.post_id = $1 AND c.id = $2
	`, postComment.TableName()),
		postComment,
		phi.URLParam(ctx, "postID"),
		phi.URLParam(ctx, "commentID"))

	return postComment
}
