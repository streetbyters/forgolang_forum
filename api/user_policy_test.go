package api

import (
	"fmt"
	"forgolang_forum/database/model"
	"github.com/valyala/fasthttp"
	"testing"
)

type UserPolicyTest struct {
	*Suite
}

func (s UserPolicyTest) SetupSuite() {
	SetupSuite(s.Suite)
}

func (s UserPolicyTest) Test_Should_403Err_ListAllUsersWithUserRole() {
	UserAuth(s.Suite, "user")

	resp := s.JSON(Get, "/api/v1/user", nil)

	s.Equal(resp.Status, fasthttp.StatusForbidden)

	defaultLogger.LogInfo("Should be 403 error list all users with user role")
}

func (s UserPolicyTest) Test_Should_403Err_ShowUserWithUserRoleIfOtherUserIdentifier() {
	UserAuth(s.Suite, "user")

	pwd := "123456"
	user := model.NewUser(&pwd)
	user.Username = "policy-user"
	user.Email = "policy-user@mail.com"
	user.IsActive = true
	err := s.API.App.Database.Insert(new(model.User),
		user,
		"id", "inserted_at", "updated_at")
	s.Nil(err)

	resp := s.JSON(Get, fmt.Sprintf("/api/v1/user/%d", user.ID), nil)

	s.Equal(resp.Status, fasthttp.StatusForbidden)

	defaultLogger.LogInfo("Should be 403 error show a user with given identifier " +
		"and user role if other user identifier")
}

func (s UserPolicyTest) Test_ShowUserWithUserRole() {
	UserAuth(s.Suite, "user")

	resp := s.JSON(Get, fmt.Sprintf("/api/v1/user/%d", s.Auth.User.ID), nil)

	s.Equal(resp.Status, fasthttp.StatusOK)

	defaultLogger.LogInfo("Show a user with given identifier and user role")
}

func (s UserPolicyTest) Test_Should_403Err_CreateUserWithValidParamsAndUserRole() {
	user := new(model.User)
	user.Username = "policy-user"
	user.Email = "policy-user@mail.com"
	user.Password = "123456"
	user.IsActive = true

	resp := s.JSON(Post, "/api/v1/user", user)

	s.Equal(resp.Status, fasthttp.StatusForbidden)

	defaultLogger.LogInfo("Should be 403 error create a user with valid params " +
		"and user role")
}

func (s UserPolicyTest) Test_Should_403Err_UpdateWithGivenIdentifierAndValidParamsAndUserRoleIfOtherUserIdentifier() {
	UserAuth(s.Suite, "user")

	pwd := "123456"
	user := model.NewUser(&pwd)
	user.Username = "policy-user-3"
	user.Email = "policy-user3@mail.com"
	user.IsActive = true
	err := s.API.App.Database.Insert(new(model.User),
		user,
		"id", "inserted_at", "updated_at")
	s.Nil(err)
	userID := user.ID

	user = new(model.User)
	user.Username = "policy-user"
	user.Email = "policy-user@mail.com"
	user.IsActive = true

	resp := s.JSON(Put, fmt.Sprintf("/api/v1/user/%d", userID), user)

	s.Equal(resp.Status, fasthttp.StatusForbidden)

	defaultLogger.LogInfo("Should be 403 error update a user with given identifier " +
		"and valid params and user role if other user identifier")
}

func (s UserPolicyTest) Test_UpdateWithGivenIdentifierAndValidParamsAndUserRole() {
	UserAuth(s.Suite, "user")

	user := new(model.User)
	user.Username = "policy-user-edit"
	user.Email = "policy-user-edit@mail.com"
	user.IsActive = true

	resp := s.JSON(Put, fmt.Sprintf("/api/v1/user/%d", s.Auth.User.ID), user)

	s.Equal(resp.Status, fasthttp.StatusOK)

	defaultLogger.LogInfo("Update a user with given identifier " +
		"and valid params and user role")
}

func (s UserPolicyTest) Test_Should_403Err_DeleteUserWithGivenIdentifierAndUserRole() {
	UserAuth(s.Suite, "user")

	pwd := "123456"
	user := model.NewUser(&pwd)
	user.Username = "policy-user-4"
	user.Email = "policy-user4@mail.com"
	user.IsActive = true
	err := s.API.App.Database.Insert(new(model.User),
		user,
		"id", "inserted_at", "updated_at")
	s.Nil(err)

	resp := s.JSON(Delete, fmt.Sprintf("/api/v1/user/%d", user.ID), nil)

	s.Equal(resp.Status, fasthttp.StatusForbidden)

	defaultLogger.LogInfo("Should be 403 error delete a user with given " +
		"identifier and user role")
}

func (s UserPolicyTest) TearDownSuite() {
	TearDownSuite(s.Suite)
}

func Test_UserPolicy(t *testing.T) {
	s := UserPolicyTest{NewSuite()}
	Run(t, s)
}
