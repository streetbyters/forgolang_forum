package api

import (
	"fmt"
	"forgolang_forum/database/model"
	"github.com/valyala/fasthttp"
	"testing"
)

type PostPolicyTest struct {
	*Suite
}

func (s PostPolicyTest) SetupSuite() {
	SetupSuite(s.Suite)
}

func (s PostPolicyTest) Test_Should_403Err_CreatePostWithValidParamsIfUnauthorized() {
	s.Suite.Auth = struct {
		User  *model.User
		Token string
	}{User: nil, Token: ""}

	postDep := new(model.PostDEP)
	postDep.Title.SetValid("Post title")
	postDep.Content.SetValid("Post content")
	postDep.Description.SetValid("Post description")

	response := s.JSON(Post, "/api/v1/post", postDep)

	s.Equal(response.Status, fasthttp.StatusForbidden)
}

func (s PostPolicyTest) Test_CreatePostWithValidParamsAndUserRole() {
	UserAuth(s.Suite, "user")

	postDep := new(model.PostDEP)
	postDep.Title.SetValid("Post title 2")
	postDep.Content.SetValid("Post content")
	postDep.Description.SetValid("Post description")

	response := s.JSON(Post, "/api/v1/post", postDep)

	s.Equal(response.Status, fasthttp.StatusCreated)

	defaultLogger.LogInfo("Create post with valid params and user role")
}

func (s PostPolicyTest) Test_CreatePostWithValidParamsAndModeratorRole() {
	UserAuth(s.Suite, "moderator")

	postDep := new(model.PostDEP)
	postDep.Title.SetValid("Post title 3")
	postDep.Content.SetValid("Post content")
	postDep.Description.SetValid("Post description")

	response := s.JSON(Post, "/api/v1/post", postDep)

	s.Equal(response.Status, fasthttp.StatusCreated)

	defaultLogger.LogInfo("Create post with valid params and moderator role")
}

func (s PostPolicyTest) Test_DeletePostWithGivenIdentifierAndUserRole() {
	UserAuth(s.Suite, "user")

	post := model.NewPost(s.Auth.User.ID)
	err := s.API.App.Database.Insert(new(model.Post), post, "id")
	s.Nil(err)

	response := s.JSON(Delete, fmt.Sprintf("/api/v1/post/%d", post.ID), nil)

	s.Equal(response.Status, fasthttp.StatusNoContent)

	defaultLogger.LogInfo("Delete post with given identifier and user role")
}

func (s PostPolicyTest) Test_Should_403Err_DeletePostWithGivenIdentifierAndUserRoleIfPostAuthorOtherUser() {
	UserAuth(s.Suite, "user")

	pwd := "123456"
	user := model.NewUser(&pwd)
	user.Username = "post-user"
	user.Email = "post-user@mail.com"
	user.IsActive = true
	err := s.API.App.Database.Insert(new(model.User),
		user,
		"id", "inserted_at", "updated_at")
	s.Nil(err)

	post := model.NewPost(user.ID)
	err = s.API.App.Database.Insert(new(model.Post), post, "id")
	s.Nil(err)

	response := s.JSON(Delete, fmt.Sprintf("/api/v1/post/%d", post.ID), nil)

	s.Equal(response.Status, fasthttp.StatusForbidden)

	defaultLogger.LogInfo("Should 403 error delete post with given identifier and " +
		"user role if post author other user")
}

func (s PostPolicyTest) Test_DeletePostWithGivenIdentifierAndModeratorRole() {
	UserAuth(s.Suite, "moderator")

	post := model.NewPost(s.Auth.User.ID)
	err := s.API.App.Database.Insert(new(model.Post), post, "id")
	s.Nil(err)

	response := s.JSON(Delete, fmt.Sprintf("/api/v1/post/%d", post.ID), nil)

	s.Equal(response.Status, fasthttp.StatusNoContent)

	defaultLogger.LogInfo("Delete post with given identifier and user role")
}

func (s PostPolicyTest) Test_Should_403Err_DeletePostWithGivenIdentifierAndModeratorRoleIfPostAuthorOtherUser() {
	UserAuth(s.Suite, "moderator")

	pwd := "123456"
	user := model.NewUser(&pwd)
	user.Username = "post-user-2"
	user.Email = "post-user-2@mail.com"
	user.IsActive = true
	err := s.API.App.Database.Insert(new(model.User),
		user,
		"id", "inserted_at", "updated_at")
	s.Nil(err)

	post := model.NewPost(user.ID)
	err = s.API.App.Database.Insert(new(model.Post), post, "id")
	s.Nil(err)

	response := s.JSON(Delete, fmt.Sprintf("/api/v1/post/%d", post.ID), nil)

	s.Equal(response.Status, fasthttp.StatusForbidden)

	defaultLogger.LogInfo("Should 403 error delete post with given identifier and " +
		"moderator role if post author other user")
}

func (s PostPolicyTest) TearDownSuite() {
	TearDownSuite(s.Suite)
}

func Test_PostPolicy(t *testing.T) {
	s := PostPolicyTest{NewSuite()}
	Run(t, s)
}
