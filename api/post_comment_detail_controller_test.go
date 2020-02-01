package api

import (
	"fmt"
	"forgolang_forum/database/model"
	"github.com/valyala/fasthttp"
	"testing"
)

type PostCommentDetailControllerTest struct {
	*Suite
}

func (s PostCommentDetailControllerTest) SetupSuite() {
	SetupSuite(s.Suite)
	UserAuth(s.Suite)
}

func (s PostCommentDetailControllerTest) Test_CreatePostCommentDetailWithValidParams() {
	post := model.NewPost(s.Auth.User.ID)
	err := s.API.GetDB().Insert(new(model.Post), post, "id")
	s.Nil(err)

	pwd := "12345"
	user := model.NewUser(&pwd)
	user.Username = "test-comment-user"
	user.Email = "test-comment-user@mail.com"
	err = s.API.GetDB().Insert(new(model.User), user, "id")
	s.Nil(err)

	postComment := model.NewPostComment(post.ID, user.ID)
	err = s.API.GetDB().Insert(new(model.PostComment), postComment, "id")
	s.Nil(err)

	commentDetail := new(model.PostCommentDetail)
	commentDetail.Comment = "Test Comment"

	response := s.JSON(Post, fmt.Sprintf("/api/v1/post/%d/comment/%d/detail",
		post.ID, postComment.ID), commentDetail)

	s.Equal(response.Status, fasthttp.StatusCreated)

	data, _ := response.Success.Data.(map[string]interface{})
	s.Greater(data["id"], float64(0))
	s.Equal(data["post_id"], float64(post.ID))
	s.Equal(data["comment_id"], float64(postComment.ID))
	s.Equal(data["comment"], "Test Comment")

	defaultLogger.LogInfo("Create post comment detail with valid params")
}

func (s PostCommentDetailControllerTest) Test_Should_422Err_CreatePostCommentDetailWithInvalidParams() {
	post := model.NewPost(s.Auth.User.ID)
	err := s.API.GetDB().Insert(new(model.Post), post, "id")
	s.Nil(err)

	pwd := "12345"
	user := model.NewUser(&pwd)
	user.Username = "test-comment-user-2"
	user.Email = "test-comment-user-2@mail.com"
	err = s.API.GetDB().Insert(new(model.User), user, "id")
	s.Nil(err)

	postComment := model.NewPostComment(post.ID, user.ID)
	err = s.API.GetDB().Insert(new(model.PostComment), postComment, "id")
	s.Nil(err)

	commentDetail := new(model.PostCommentDetail)
	commentDetail.Comment = "Te"

	response := s.JSON(Post, fmt.Sprintf("/api/v1/post/%d/comment/%d/detail",
		post.ID, postComment.ID), commentDetail)

	s.Equal(response.Status, fasthttp.StatusUnprocessableEntity)

	defaultLogger.LogInfo("Should be 422 error create post comment detail " +
		"with invalid params")
}

func (s PostCommentDetailControllerTest) Test_Should_422Err_CreatePostCommentDetailWithValidParamsIfRelationalError() {
	post := model.NewPost(s.Auth.User.ID)
	err := s.API.GetDB().Insert(new(model.Post), post, "id")
	s.Nil(err)

	pwd := "12345"
	user := model.NewUser(&pwd)
	user.Username = "test-comment-user-3"
	user.Email = "test-comment-user-3@mail.com"
	err = s.API.GetDB().Insert(new(model.User), user, "id")
	s.Nil(err)

	commentDetail := new(model.PostCommentDetail)
	commentDetail.Comment = "Test comment"

	response := s.JSON(Post, fmt.Sprintf("/api/v1/post/%d/comment/999999999/detail",
		post.ID), commentDetail)

	s.Equal(response.Status, fasthttp.StatusUnprocessableEntity)

	data, _ := response.Error.Errors.(map[string]interface{})
	s.Equal(data["comment_id"], "does not exists")

	defaultLogger.LogInfo("Should be 422 error create post comment detail " +
		"with valid params if relational error")
}

func (s PostCommentDetailControllerTest) TearDownSuite() {
	TearDownSuite(s.Suite)
}

func Test_PostCommentDetailController(t *testing.T) {
	s := PostCommentDetailControllerTest{NewSuite()}
	Run(t, s)
}
