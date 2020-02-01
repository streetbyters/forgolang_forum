package api

import (
	"fmt"
	"forgolang_forum/database/model"
	"github.com/gosimple/slug"
	"github.com/valyala/fasthttp"
	"testing"
)

type PostControllerTest struct {
	*Suite
}

func (s PostControllerTest) SetupSuite() {
	SetupSuite(s.Suite)
	UserAuth(s.Suite)
}

func (s PostControllerTest) Test_CreatePostWithValidParams() {
	postDep := new(model.PostDEP)
	postDep.Title.SetValid("Post title")
	postDep.Content.SetValid("Post content")
	postDep.Description.SetValid("Post description")

	response := s.JSON(Post, "/api/v1/post", postDep)

	s.Equal(response.Status, fasthttp.StatusCreated)
	data, _ := response.Success.Data.(map[string]interface{})
	s.Greater(data["id"], float64(0))
	s.Equal(data["author_id"], float64(s.Auth.User.ID))
	s.Equal(data["title"], "Post title")
	s.Equal(data["content"], "Post content")
	s.Equal(data["description"], "Post description")
	s.Equal(data["slug"], slug.Make(postDep.Title.String))

	defaultLogger.LogInfo("Create post with valid params")
}

func (s PostControllerTest) Test_Should_422Err_CreatePostWithInvalidParams() {
	postDep := new(model.PostDEP)
	postDep.Title.SetValid("Po")
	postDep.Description.SetValid("Post description")

	response := s.JSON(Post, "/api/v1/post", postDep)

	s.Equal(response.Status, fasthttp.StatusUnprocessableEntity)

	defaultLogger.LogInfo("Should 422 error create post with invalid params")
}

func (s PostControllerTest) Test_Should_422Err_CreatePostWithValidParamsIfSlugNotUnique() {
	post := model.NewPost(s.Auth.User.ID)
	err := s.API.GetDB().Insert(new(model.Post), post, "id")
	s.Nil(err)

	postSlug := model.NewPostSlug(post.ID, s.Auth.User.ID)
	postSlug.Slug = slug.Make("Post title slug")
	err = s.API.GetDB().Insert(new(model.PostSlug), postSlug, "id")
	s.Nil(err)

	postDep := new(model.PostDEP)
	postDep.Title.SetValid("Post title slug")
	postDep.Content.SetValid("Post content")
	postDep.Description.SetValid("Post description")

	response := s.JSON(Post, "/api/v1/post", postDep)

	s.Equal(response.Status, fasthttp.StatusUnprocessableEntity)
	data, _ := response.Error.Errors.(map[string]interface{})
	s.Equal(data["slug"], "has been already taken")

	defaultLogger.LogInfo("Should be 422 error create post with valid params " +
		"if slug has been already taken")
}

func (s PostControllerTest) Test_DeletePostWithGivenIdentifier() {
	post := model.NewPost(s.Auth.User.ID)
	err := s.API.GetDB().Insert(new(model.Post), post, "id")
	s.Nil(err)

	response := s.JSON(Delete, fmt.Sprintf("/api/v1/post/%d", post.ID), nil)

	s.Equal(response.Status, fasthttp.StatusNoContent)

	defaultLogger.LogInfo("Delete post with given identifier")
}

func (s PostControllerTest) Test_Shoul_404Err_DeletePostWithGivenIdentifierIfNotExists() {
	response := s.JSON(Delete, "/api/v1/post/999999999", nil)

	s.Equal(response.Status, fasthttp.StatusNotFound)

	defaultLogger.LogInfo("Should be 404 error delete post with given identifier " +
		"if does not exists")
}

func (s PostControllerTest) TearDownSuite() {
	TearDownSuite(s.Suite)
}

func Test_PostController(t *testing.T) {
	s := PostControllerTest{NewSuite()}
	Run(t, s)
}
