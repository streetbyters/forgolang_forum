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
	"fmt"
	"forgolang_forum/cmn"
	"forgolang_forum/database"
	"forgolang_forum/database/model"
	"github.com/valyala/fasthttp"
	"testing"
	"time"
)

type UserControllerTest struct {
	*Suite
}

func (s UserControllerTest) SetupSuite() {
	SetupSuite(s.Suite)
	UserAuth(s.Suite)
}

func (s UserControllerTest) Test_ListAllUsers() {
	for i := 0; i < 50; i++ {
		pwd := "12345"
		user := model.NewUser(&pwd)
		user.Username = fmt.Sprintf("new-user%d", i)
		user.Email = fmt.Sprintf("new-user%d@mail.com", i)
		user.IsActive = true

		err := s.API.App.Database.Insert(new(model.User), user, "id")
		s.Nil(err)

		roleAssignment := model.NewUserRoleAssignment(user.ID, 3)
		roleAssignment.SourceUserID.SetValid(user.ID)
		err = s.API.App.Database.Insert(new(model.UserRoleAssignment), roleAssignment, "id")
		s.Nil(err)

		userState := model.NewUserState(user.ID)
		userState.State = database.Active
		userState.SourceUserID.SetValid(user.ID)
		err = s.API.App.Database.Insert(new(model.UserState), userState, "id")
		s.Nil(err)
	}

	resp := s.JSON(Get, "/api/v1/user", nil)

	s.Equal(resp.Status, fasthttp.StatusOK)
	s.Greater(resp.Success.TotalCount, int64(49))

	defaultLogger.LogInfo("List all users")
}

func (s UserControllerTest) Test_ShowUserWithGivenIdentifier() {
	pwd := "123456"
	user := model.NewUser(&pwd)
	user.Username = "show-user"
	user.Email = "show-user@mail.com"
	user.IsActive = true
	err := s.API.App.Database.Insert(new(model.User),
		user,
		"id", "inserted_at", "updated_at")
	s.Nil(err)

	resp := s.JSON(Get, fmt.Sprintf("/api/v1/user/%d", user.ID), nil)

	s.Equal(resp.Status, fasthttp.StatusOK)

	data, _ := resp.Success.Data.(map[string]interface{})

	s.Equal(data["id"], float64(user.ID))
	s.Equal(data["username"], "show-user")
	s.Equal(data["email"], "show-user@mail.com")
	s.Equal(data["is_active"], true)
	s.Equal(data["inserted_at"], user.InsertedAt.Format(time.RFC3339Nano))
	s.Equal(data["updated_at"], user.UpdatedAt.Format(time.RFC3339Nano))

	defaultLogger.LogInfo("Show a user with given identifier")
}

func (s UserControllerTest) Test_ShowCachedUserWithGivenIdentifier() {
	pwd := "123456"
	user := model.NewUser(&pwd)
	user.Username = "show-user2"
	user.Email = "show-user2@mail.com"
	user.IsActive = true
	err := s.API.App.Database.Insert(new(model.User),
		user,
		"id", "inserted_at", "updated_at")
	s.Nil(err)

	s.API.App.Cache.Set(fmt.Sprintf("%s:%d",
		cmn.GetRedisKey("user", "one"), user.ID),
		user.ToJSON(), 0)

	resp := s.JSON(Get, fmt.Sprintf("/api/v1/user/%d", user.ID), nil)

	s.Equal(resp.Status, fasthttp.StatusOK)

	data, _ := resp.Success.Data.(map[string]interface{})

	s.Equal(data["id"], float64(user.ID))
	s.Equal(data["username"], "show-user2")
	s.Equal(data["email"], "show-user2@mail.com")
	s.Equal(data["is_active"], true)
	s.Equal(data["inserted_at"], user.InsertedAt.Format(time.RFC3339Nano))
	s.Equal(data["updated_at"], user.UpdatedAt.Format(time.RFC3339Nano))

	defaultLogger.LogInfo("Show a cached user with given identifier")
}

func (s UserControllerTest) Test_Should_404Err_ShowUserWithGivenIdentifierIfNotExists() {
	resp := s.JSON(Get, "/api/v1/user/999999999", nil)

	s.Equal(resp.Status, fasthttp.StatusNotFound)

	defaultLogger.LogInfo("Should be 404 error show a user with given identifier " +
		"if does not exists")
}

func (s UserControllerTest) Test_CreateUserWithValidParams() {
	user := new(model.User)
	user.Email = "create-user@mail.com"
	user.Username = "create-user"
	user.Password = "123456"
	user.IsActive = true

	resp := s.JSON(Post, "/api/v1/user", user)

	s.Equal(resp.Status, fasthttp.StatusCreated)

	data, _ := resp.Success.Data.(map[string]interface{})

	s.Greater(data["id"], float64(0))
	s.Equal(data["email"], "create-user@mail.com")
	s.Equal(data["username"], "create-user")
	s.Equal(data["password"], "****")
	s.Equal(data["is_active"], true)
	s.NotNil(data["inserted_at"])
	s.NotNil(data["updated_at"])

	defaultLogger.LogInfo("Create a user with valid params")
}

func (s UserControllerTest) Test_Should_422Error_CreateUserWithInvalidParams() {
	user := new(model.User)
	user.Email = "create-user"
	user.Username = "c"
	user.IsActive = true

	resp := s.JSON(Post, "/api/v1/user", user)

	s.Equal(resp.Status, fasthttp.StatusUnprocessableEntity)

	defaultLogger.LogInfo("Should be 422 error create a user with invalid params")
}

func (s UserControllerTest) Test_Should_422Error_CreateUserWithValidParamsIFEmailNotUnique() {
	user := new(model.User)
	user.Email = "create-user2@mail.com"
	user.Username = "create-user2"
	user.Password = "123456"
	user.IsActive = true
	err := s.API.App.Database.Insert(new(model.User), user, "id")
	s.Nil(err)

	user = new(model.User)
	user.Email = "create-user2@mail.com"
	user.Username = "create-user-3"
	user.Password = "123456"
	user.IsActive = false

	resp := s.JSON(Post, "/api/v1/user", user)

	s.Equal(resp.Status, fasthttp.StatusUnprocessableEntity)

	defaultLogger.LogInfo("Should be 422 error create a user with valid params " +
		"if email has been already taken")
}

func (s UserControllerTest) Test_UpdateUserWithGivenIdentifierAndValidParams() {
	user := new(model.User)
	user.Email = "update-user@mail.com"
	user.Username = "update-user"
	user.Password = "123456"
	user.IsActive = true
	err := s.API.App.Database.Insert(new(model.User), user, "id")
	s.Nil(err)

	userRequest := new(model.User)
	userRequest.Email = "update-user@mail.com"
	userRequest.Username = "update-user-edit"
	userRequest.IsActive = true

	resp := s.JSON(Put, fmt.Sprintf("/api/v1/user/%d", user.ID), userRequest)

	s.Equal(resp.Status, fasthttp.StatusOK)

	data, _ := resp.Success.Data.(map[string]interface{})

	s.Equal(data["id"], float64(user.ID))
	s.Equal(data["email"], "update-user@mail.com")
	s.Equal(data["username"], "update-user-edit")
	s.Equal(data["is_active"], true)
	s.NotNil(data["inserted_at"])
	s.NotNil(data["updated_at"])

	defaultLogger.LogInfo("Update a user with given identifier and valid params")
}

func (s UserControllerTest) Test_Should_422Error_UpdateUserWithGivenIdentifierAndInalidParams() {
	user := new(model.User)
	user.Email = "update-user2@mail.com"
	user.Username = "update-user2"
	user.Password = "123456"
	user.IsActive = true
	err := s.API.App.Database.Insert(new(model.User), user, "id")
	s.Nil(err)

	userRequest := new(model.User)
	userRequest.Email = "update-user"
	userRequest.Username = "update"
	userRequest.IsActive = true

	resp := s.JSON(Put, fmt.Sprintf("/api/v1/user/%d", user.ID), userRequest)

	s.Equal(resp.Status, fasthttp.StatusUnprocessableEntity)

	defaultLogger.LogInfo("Should be 422 error update a user with given " +
		"identifier and invalid params")
}

func (s UserControllerTest) Test_Should_422Error_UpdateUserWithGivenIdentifierAndValidParamsIfUsernameNotUnique() {
	user := new(model.User)
	user.Email = "update-user3@mail.com"
	user.Username = "update-user3"
	user.Password = "123456"
	user.IsActive = true
	err := s.API.App.Database.Insert(new(model.User), user, "id")
	s.Nil(err)

	user = new(model.User)
	user.Email = "update-user4@mail.com"
	user.Username = "update-user4"
	user.Password = "123456"
	user.IsActive = true
	err = s.API.App.Database.Insert(new(model.User), user, "id")
	s.Nil(err)

	userRequest := new(model.User)
	userRequest.Email = "update-user-edit@mail.com"
	userRequest.Username = "update-user3"
	userRequest.IsActive = true

	resp := s.JSON(Put, fmt.Sprintf("/api/v1/user/%d", user.ID), userRequest)

	s.Equal(resp.Status, fasthttp.StatusUnprocessableEntity)

	defaultLogger.LogInfo("Should be 422 error update a user with given " +
		"identifier and valid params if username not unique")
}

func (s UserControllerTest) Test_DeleteUserWithGivenIdentifier() {
	pwd := "123456"
	user := model.NewUser(&pwd)
	user.Username = "delete-user"
	user.Email = "delete-user@mail.com"
	user.IsActive = true
	err := s.API.App.Database.Insert(new(model.User),
		user,
		"id", "inserted_at", "updated_at")
	s.Nil(err)

	s.API.App.Cache.Set(fmt.Sprintf("%s:%d",
		cmn.GetRedisKey("user", "one"), user.ID),
		user.ToJSON(), 0)

	resp := s.JSON(Delete, fmt.Sprintf("/api/v1/user/%d", user.ID), nil)

	s.Equal(resp.Status, fasthttp.StatusNoContent)

	cache := s.API.App.Cache.Get(fmt.Sprintf("%s:%d",
		cmn.GetRedisKey("user", "one"), user.ID))
	s.Equal(cache.Err().Error(), "redis: nil")

	defaultLogger.LogInfo("Delete user with given identifier")
}

func (s UserControllerTest) Test_Should_404Err_DeleteUserWithGivenIdentifierIfNotExists() {
	resp := s.JSON(Delete,"/api/v1/user/999999999", nil)

	s.Equal(resp.Status, fasthttp.StatusNotFound)

	defaultLogger.LogInfo("Should 404 error delete user with given identifier " +
		"if user does not exists")
}

func (s UserControllerTest) TearDownSuite() {
	TearDownSuite(s.Suite)
}

func Test_UserController(t *testing.T) {
	s := UserControllerTest{NewSuite()}
	Run(t, s)
}
