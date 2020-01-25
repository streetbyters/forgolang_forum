package api

import (
	"fmt"
	"forgolang_forum/database/model"
	"github.com/valyala/fasthttp"
	"testing"
)

type PostSlugControllerTest struct {
	*Suite
}

func (s PostSlugControllerTest) SetupSuite() {
	SetupSuite(s.Suite)
	UserAuth(s.Suite)
}

func (s PostSlugControllerTest) Test_CretePostSlugWithValidParams() {
	post := new(model.Post)
	post.AuthorID = s.Auth.User.ID
	err := s.API.App.Database.Insert(new(model.Post), post, "id")
	s.Nil(err)

	postSlug := new(model.PostSlug)
	postSlug.Slug = "post-slug-create"

	response := s.JSON(Post, fmt.Sprintf("/api/v1/post/%d/slug", post.ID), postSlug)

	s.Equal(response.Status, fasthttp.StatusCreated)
	data, _ := response.Success.Data.(map[string]interface{})

	s.Greater(data["id"], float64(0))
	s.Equal(data["post_id"], float64(post.ID))
	s.Equal(data["source_user_id"], float64(s.Auth.User.ID))
	s.Equal(data["slug"], "post-slug-create")
	s.NotNil(data["inserted_at"])

	defaultLogger.LogInfo("Create post slug with valid params")
}

func (s PostSlugControllerTest) Test_Should_422Err_CreatePostSlugWithInvalidParams() {
	post := new(model.Post)
	post.AuthorID = s.Auth.User.ID
	err := s.API.App.Database.Insert(new(model.Post), post, "id")
	s.Nil(err)

	postSlug := new(model.PostSlug)

	response := s.JSON(Post, fmt.Sprintf("/api/v1/post/%d/slug", post.ID), postSlug)

	s.Equal(response.Status, fasthttp.StatusUnprocessableEntity)

	defaultLogger.LogInfo("Should be 422 error crete a post slug " +
		"with invalid params")
}

func (s PostSlugControllerTest) Test_Should_422Err_CreatePostSlugWithValidParamsIfRelationalError() {
	postSlug := new(model.PostSlug)

	response := s.JSON(Post, "/api/v1/post/999999999/slug", postSlug)

	s.Equal(response.Status, fasthttp.StatusUnprocessableEntity)

	defaultLogger.LogInfo("Should be 422 error crete a post slug " +
		"with valid params if relational error")
}

func (s PostSlugControllerTest) TearDownSuite() {
	TearDownSuite(s.Suite)
}

func Test_PostSlugController(t *testing.T) {
	s := PostSlugControllerTest{NewSuite()}
	Run(t, s)
}
