package api

import (
	"fmt"
	"forgolang_forum/database/model"
	"github.com/valyala/fasthttp"
	"testing"
)

type PostCommentPolicyTest struct {
	*Suite
}

func (s PostCommentPolicyTest) SetupSuite() {
	SetupSuite(s.Suite)
}

func (s PostCommentPolicyTest) Test_ListAllPostCommentsWithoutAuth() {
	pwd := "123456"
	user := model.NewUser(&pwd)
	user.Username = "test-comment-user"
	user.Email = "test-comment-user@mail.com"
	err := s.API.GetDB().Insert(new(model.User), user, "id")
	s.Nil(err)

	post := model.NewPost(user.ID)
	err = s.API.GetDB().Insert(new(model.Post), post, "id")
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

	defaultLogger.LogInfo("List all post comments without auth")
}

func (s PostCommentPolicyTest) Test_CreatePostCommentWithUserRole() {
	UserAuth(s.Suite, "user")

	post := model.NewPost(s.Auth.User.ID)
	err := s.API.GetDB().Insert(new(model.Post), post, "id")
	s.Nil(err)

	postComment := new(model.PostComment)

	response := s.JSON(Post, fmt.Sprintf("/api/v1/post/%d/comment", post.ID),
		postComment)

	s.Equal(response.Status, fasthttp.StatusCreated)

	defaultLogger.LogInfo("Create post comment with user role")
}

func (s PostCommentPolicyTest) Test_DeletePostCommentWithUserRole() {
	UserAuth(s.Suite, "user")

	pwd := "12345"
	user := model.NewUser(&pwd)
	user.Username = "test-comment-user-2"
	user.Email = "test-comment-user-2@mail.com"
	err := s.API.GetDB().Insert(new(model.User), user, "id")
	s.Nil(err)

	post := model.NewPost(user.ID)
	err = s.API.GetDB().Insert(new(model.Post), post, "id")
	s.Nil(err)

	postComment := model.NewPostComment(post.ID, s.Auth.User.ID)
	err = s.API.GetDB().Insert(new(model.PostComment), postComment, "id")
	s.Nil(err)

	response := s.JSON(Delete, fmt.Sprintf("/api/v1/post/%d/comment/%d",
		post.ID, postComment.ID), nil)

	s.Equal(response.Status, fasthttp.StatusNoContent)

	defaultLogger.LogInfo("Delete post comment with given identifier and " +
		"user role")
}

func (s PostCommentPolicyTest) Test_Should_403Err_DeletePostCommentWithUserRoleIfCommentOtherUser() {
	UserAuth(s.Suite, "user")

	pwd := "12345"
	user := model.NewUser(&pwd)
	user.Username = "test-comment-user-3"
	user.Email = "test-comment-user-3@mail.com"
	err := s.API.GetDB().Insert(new(model.User), user, "id")
	s.Nil(err)

	post := model.NewPost(user.ID)
	err = s.API.GetDB().Insert(new(model.Post), post, "id")
	s.Nil(err)

	postComment := model.NewPostComment(post.ID, user.ID)
	err = s.API.GetDB().Insert(new(model.PostComment), postComment, "id")
	s.Nil(err)

	response := s.JSON(Delete, fmt.Sprintf("/api/v1/post/%d/comment/%d",
		post.ID, postComment.ID), nil)

	s.Equal(response.Status, fasthttp.StatusForbidden)

	defaultLogger.LogInfo("Should be 403 error delete post comment with " +
		"given identifier and user role if comment other user")
}

func (s PostCommentPolicyTest) TearDownSuite() {
	TearDownSuite(s.Suite)
}

func Test_PostCommentPolicy(t *testing.T) {
	s := PostCommentPolicyTest{NewSuite()}
	Run(t, s)
}
