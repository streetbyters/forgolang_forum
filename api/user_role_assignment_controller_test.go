package api

import (
	"fmt"
	"forgolang_forum/database/model"
	"github.com/valyala/fasthttp"
	"testing"
)

type UserRoleAssignmentControllerTest struct {
	*Suite
}

func (s UserRoleAssignmentControllerTest) SetupSuite() {
	SetupSuite(s.Suite)
	UserAuth(s.Suite)
}

func (s UserRoleAssignmentControllerTest) Test_CreateUserRoleAssignmentWithValidParams() {
	pwd := "12345"
	user := model.NewUser(&pwd)
	user.Email = "akdilsiz@tecpor.com"
	user.Username = "akdilsiz"
	err := s.API.GetDB().Insert(new(model.User), user, "id")
	s.Nil(err)

	roleAssignment := new(model.UserRoleAssignment)
	roleAssignment.RoleID = 2

	response := s.JSON(Post, fmt.Sprintf("/api/v1/user/%d/role_assignment", user.ID),
		roleAssignment)

	s.Equal(response.Status, fasthttp.StatusCreated)

	data, _ := response.Success.Data.(map[string]interface{})

	s.Greater(data["id"], float64(0))
	s.Equal(data["user_id"], float64(user.ID))
	s.Equal(data["role_id"], float64(2))
	s.NotNil(data["inserted_at"])

	defaultLogger.LogInfo("Create user role assignment with valid params")
}

func (s UserRoleAssignmentControllerTest) Test_Should_400Err_CreateUserRoleAssignmentWithValidParamsIfInvalidIdentifier() {
	roleAssignment := model.NewUserRoleAssignment(1, 3)
	response := s.JSON(Post, "/api/v1/user/userID/role_assignment",
		roleAssignment)

	s.Equal(response.Status, fasthttp.StatusBadRequest)

	defaultLogger.LogInfo("Should be 400 error create user role assignment " +
		"with valid params if invalid identifier")
}

func (s UserRoleAssignmentControllerTest) Test_Should_422Err_CreateUserRoleAssignmentWithInvalidParams() {
	pwd := "12345"
	user := model.NewUser(&pwd)
	user.Email = "akdilsiz-2@tecpor.com"
	user.Username = "akdilsiz-2"
	err := s.API.GetDB().Insert(new(model.User), user, "id")
	s.Nil(err)

	roleAssignment := new(model.UserRoleAssignment)

	response := s.JSON(Post, fmt.Sprintf("/api/v1/user/%d/role_assignment", user.ID),
		roleAssignment)

	s.Equal(response.Status, fasthttp.StatusUnprocessableEntity)

	defaultLogger.LogInfo("Should be 422 error create user role assignment with " +
		"invalid params")
}

func (s UserRoleAssignmentControllerTest) Test_Should_422Err_CreateUserRoleAssignmentWithValidParamsIfRelationalError() {
	pwd := "12345"
	user := model.NewUser(&pwd)
	user.Email = "akdilsiz-3@tecpor.com"
	user.Username = "akdilsiz-3"
	err := s.API.GetDB().Insert(new(model.User), user, "id")
	s.Nil(err)

	roleAssignment := new(model.UserRoleAssignment)
	roleAssignment.RoleID = 999999999

	response := s.JSON(Post, fmt.Sprintf("/api/v1/user/%d/role_assignment", user.ID),
		roleAssignment)

	s.Equal(response.Status, fasthttp.StatusUnprocessableEntity)

	defaultLogger.LogInfo("Should be 422 error create user role assignment with " +
		"valid params if relational error")
}

func (s UserRoleAssignmentControllerTest) TearDownSuite() {
	TearDownSuite(s.Suite)
}

func Test_UserRoleAssignmentController(t *testing.T) {
	s := UserRoleAssignmentControllerTest{NewSuite()}
	Run(t, s)
}
