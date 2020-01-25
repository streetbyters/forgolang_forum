package api

import (
	"fmt"
	"forgolang_forum/database/model"
	"github.com/valyala/fasthttp"
	"testing"
)

type PostDetailPolicyTest struct {
	*Suite
}

func (s PostDetailPolicyTest) SetupSuite() {
	SetupSuite(s.Suite)
}

func (s PostDetailPolicyTest) Test_CreatePostDetailWithModeratorRole() {
	UserAuth(s.Suite, "moderator")

	post := model.NewPost(s.Auth.User.ID)
	err := s.API.App.Database.Insert(new(model.Post), post, "id")
	s.Nil(err)

	postDetail := new(model.PostDetail)
	postDetail.Title = "Post Title Moderator"
	postDetail.Content = "post Content"

	response := s.JSON(Post, fmt.Sprintf("/api/v1/post/%d/detail", post.ID), postDetail)

	s.Equal(response.Status, fasthttp.StatusCreated)

	defaultLogger.LogInfo("Create a post detail with valid params and " +
		"moderator role")
}

func (s PostDetailPolicyTest) Test_Should_403Err_CreatePostDetailWithModeratorRoleIfAuthorOtherUser() {
	UserAuth(s.Suite, "moderator")

	pwd := "123456"
	user := model.NewUser(&pwd)
	user.Username = "policy-user"
	user.Email = "policy-user@mail.com"
	user.IsActive = true
	err := s.API.App.Database.Insert(new(model.User),
		user,
		"id", "inserted_at", "updated_at")
	s.Nil(err)

	post := model.NewPost(user.ID)
	err = s.API.App.Database.Insert(new(model.Post), post, "id")
	s.Nil(err)

	postDetail := new(model.PostDetail)
	postDetail.Title = "Post Title Moderator"
	postDetail.Content = "post Content"

	response := s.JSON(Post, fmt.Sprintf("/api/v1/post/%d/detail", post.ID), postDetail)

	s.Equal(response.Status, fasthttp.StatusForbidden)

	defaultLogger.LogInfo("Should be 403 error create a post detail with valid " +
		"params and moderator role if post author other user")
}

func (s PostDetailPolicyTest) Test_CreatePostDetailWithUserRole() {
	UserAuth(s.Suite, "user")

	post := model.NewPost(s.Auth.User.ID)
	err := s.API.App.Database.Insert(new(model.Post), post, "id")
	s.Nil(err)

	postDetail := new(model.PostDetail)
	postDetail.Title = "Post Title User"
	postDetail.Content = "post Content"

	response := s.JSON(Post, fmt.Sprintf("/api/v1/post/%d/detail", post.ID), postDetail)

	s.Equal(response.Status, fasthttp.StatusCreated)

	defaultLogger.LogInfo("Create a post detail with valid params and " +
		"user role")
}

func (s PostDetailPolicyTest) Test_Should_403Err_CreatePostDetailWithUserRoleIfAuthorOtherUser() {
	UserAuth(s.Suite, "moderator")

	pwd := "123456"
	user := model.NewUser(&pwd)
	user.Username = "policy-user-2"
	user.Email = "policy-user-2@mail.com"
	user.IsActive = true
	err := s.API.App.Database.Insert(new(model.User),
		user,
		"id", "inserted_at", "updated_at")
	s.Nil(err)

	post := model.NewPost(user.ID)
	err = s.API.App.Database.Insert(new(model.Post), post, "id")
	s.Nil(err)

	postDetail := new(model.PostDetail)
	postDetail.Title = "Post Title User"
	postDetail.Content = "post Content"

	response := s.JSON(Post, fmt.Sprintf("/api/v1/post/%d/detail", post.ID), postDetail)

	s.Equal(response.Status, fasthttp.StatusForbidden)

	defaultLogger.LogInfo("Should be 403 error create a post detail with valid " +
		"params and user role if post author other user")
}

func (s PostDetailPolicyTest) TearDownSuite() {
	TearDownSuite(s.Suite)
}

func Test_PostDetailPolicy(t *testing.T) {
	s := PostDetailPolicyTest{NewSuite()}
	Run(t, s)
}
