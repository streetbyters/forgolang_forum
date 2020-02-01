// Copyright 2019 Forgolang Community
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
	"github.com/valyala/fasthttp"
	"testing"
)

type PostCommentControllerTest struct {
	*Suite
}

func (s PostCommentControllerTest) SetupSuite() {
	SetupSuite(s.Suite)
	UserAuth(s.Suite)
}

func (s PostCommentControllerTest) Test_ListAllPostComments() {
	post := model.NewPost(s.Auth.User.ID)
	err := s.API.GetDB().Insert(new(model.Post), post, "id")
	s.Nil(err)

	lastID := int64(0)
	for i := 0; i < 200; i++ {
		postComment := model.NewPostComment(post.ID, s.Auth.User.ID)
		err := s.API.GetDB().Insert(new(model.PostComment), postComment, "id")
		s.Nil(err)
		lastID = postComment.ID
	}

	response := s.JSON(Get, fmt.Sprintf("/api/v1/post/%d/comment", post.ID), nil)

	s.Equal(response.Status, fasthttp.StatusOK)
	s.Equal(response.Success.TotalCount, int64(200))
	s.Equal(len(response.Success.Data.([]interface{})), 40)
	data, _ := response.Success.Data.([]interface{})[0].(map[string]interface{})
	s.Equal(data["id"], float64(lastID))

	defaultLogger.LogInfo("List all post comments")
}

func (s PostCommentControllerTest) Test_ListAllPostCommentsWithPaginationParams() {
	post := model.NewPost(s.Auth.User.ID)
	err := s.API.GetDB().Insert(new(model.Post), post, "id")
	s.Nil(err)

	firstID := int64(0)
	for i := 0; i < 200; i++ {
		postComment := model.NewPostComment(post.ID, s.Auth.User.ID)
		err := s.API.GetDB().Insert(new(model.PostComment), postComment, "id")
		s.Nil(err)
		if i == 0 {
			firstID = postComment.ID
		}
	}

	response := s.JSON(Get, fmt.Sprintf("/api/v1/post/%d/comment?order_field=id&order_by=asc",
		post.ID), nil)

	s.Equal(response.Status, fasthttp.StatusOK)
	s.Equal(response.Success.TotalCount, int64(200))
	s.Equal(len(response.Success.Data.([]interface{})), 40)
	data, _ := response.Success.Data.([]interface{})[0].(map[string]interface{})
	s.Equal(data["id"], float64(firstID))

	defaultLogger.LogInfo("List all post comments with pagination params")
}

func (s PostCommentControllerTest) Test_CreatePostCommentWithValidParams() {
	post := model.NewPost(s.Auth.User.ID)
	err := s.API.GetDB().Insert(new(model.Post), post, "id")
	s.Nil(err)

	postComment := new(model.PostComment)

	response := s.JSON(Post, fmt.Sprintf("/api/v1/post/%d/comment", post.ID), postComment)

	s.Equal(response.Status, fasthttp.StatusCreated)
	data, _ := response.Success.Data.(map[string]interface{})
	s.Greater(data["id"], float64(0))
	s.Equal(data["post_id"], float64(post.ID))
	s.Equal(data["user_id"], float64(s.Auth.User.ID))

	count, _ := s.API.GetCache().Get(fmt.Sprintf("%s:%d",
		cmn.GetRedisKey("comment", "count"),
		post.ID)).Int64()
	s.Equal(count, int64(1))

	defaultLogger.LogInfo("Create post comment with valid params")
}

func (s PostCommentControllerTest) Test_CreatePostCommentWithValidParamsIfRelationalError() {
	postComment := new(model.PostComment)

	response := s.JSON(Post, "/api/v1/post/999999999/comment", postComment)

	s.Equal(response.Status, fasthttp.StatusUnprocessableEntity)

	defaultLogger.LogInfo("Should be 422 error create post comment with " +
		"valid params if relational error")
}

func (s PostCommentControllerTest) Test_DeletePostCommentWithGivenIdentifier() {
	post := model.NewPost(s.Auth.User.ID)
	err := s.API.GetDB().Insert(new(model.Post), post, "id")
	s.Nil(err)

	postComment := model.NewPostComment(post.ID, s.Auth.User.ID)
	err = s.API.GetDB().Insert(new(model.PostComment), postComment, "id")
	s.Nil(err)

	s.API.GetCache().Set(fmt.Sprintf("%s:%d",
		cmn.GetRedisKey("comment", "count"),
		post.ID),
		100, 0)

	response := s.JSON(Delete, fmt.Sprintf("/api/v1/post/%d/comment/%d",
		post.ID, postComment.ID), nil)

	s.Equal(response.Status, fasthttp.StatusNoContent)

	var count int64
	s.API.GetDB().DB.Get(&count, fmt.Sprintf(`
		SELECT count(c.id) FROM %s AS c
		WHERE c.id = $1
	`, postComment.TableName()),
		postComment.ID)
	s.Equal(count, int64(0))

	count, _ = s.API.GetCache().Get(fmt.Sprintf("%s:%d",
		cmn.GetRedisKey("comment", "count"),
		post.ID)).Int64()
	s.Equal(count, int64(99))

	defaultLogger.LogInfo("Delete post comment with given identifier")
}

func (s PostCommentControllerTest) TearDownSuite() {
	TearDownSuite(s.Suite)
}

func Test_PostCommentController(t *testing.T) {
	s := PostCommentControllerTest{NewSuite()}
	Run(t, s)
}
