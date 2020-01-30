package api

import (
	"context"
	"fmt"
	"forgolang_forum/database/model"
	"github.com/gosimple/slug"
	"github.com/valyala/fasthttp"
	"strconv"
	"testing"
)

type PostDetailControllerTest struct {
	*Suite
}

func (s PostDetailControllerTest) SetupSuite() {
	SetupSuite(s.Suite)
	UserAuth(s.Suite)
}

func (s PostDetailControllerTest) Test_CreatePostDetailWithValidParams() {
	post := model.NewPost(s.Auth.User.ID)
	err := s.API.App.Database.Insert(new(model.Post), post, "id")
	s.Nil(err)

	postDetail := new(model.PostDetail)
	postDetail.Title = "Post Title 1"
	postDetail.Content = "Post Content"
	postDetail.Description.SetValid("Post Description")

	postDep := new(model.PostDEP)
	postDep.ID = post.ID
	postDep.AuthorID = s.Auth.User.ID
	postDep.InsertedAt = post.InsertedAt

	_, err = s.API.App.ElasticClient.Index().
		Index("posts").
		Id(strconv.FormatInt(post.ID, 10)).
		BodyJson(postDep).
		Do(context.TODO())
	s.Nil(err)

	response := s.JSON(Post, fmt.Sprintf("/api/v1/post/%d/detail", post.ID), postDetail)

	s.Equal(response.Status, fasthttp.StatusCreated)

	data, _ := response.Success.Data.(map[string]interface{})
	s.Greater(data["id"], float64(0))
	s.Equal(data["post_id"], float64(post.ID))
	s.Equal(data["title"], "Post Title 1")
	s.Equal(data["description"], "Post Description")
	s.Equal(data["source_user_id"], float64(s.Auth.User.ID))
	s.NotNil(data["inserted_at"])

	defaultLogger.LogInfo("Create a post detail with valid params")
}

func (s PostDetailControllerTest) Test_Should_422Err_CreatePostDetailWithInvalidParams() {
	post := model.NewPost(s.Auth.User.ID)
	err := s.API.App.Database.Insert(new(model.Post), post, "id")
	s.Nil(err)

	postDetail := new(model.PostDetail)
	postDetail.Content = "Content"

	response := s.JSON(Post, fmt.Sprintf("/api/v1/post/%d/detail", post.ID), postDetail)

	s.Equal(response.Status, fasthttp.StatusUnprocessableEntity)

	defaultLogger.LogInfo("Should be 422 error create a post detail with " +
		"invalid params")
}

func (s PostDetailControllerTest) Test_Should_422Err_CreatePostWithValidParamsIfSlugNotUnique() {
	post := model.NewPost(s.Auth.User.ID)
	err := s.API.App.Database.Insert(new(model.Post), post, "id")
	s.Nil(err)

	postSlug := model.NewPostSlug(post.ID, s.Auth.User.ID)
	postSlug.Slug = slug.Make("post-title-slug")
	err = s.API.App.Database.Insert(new(model.PostSlug), postSlug, "id")
	s.Nil(err)

	post = model.NewPost(s.Auth.User.ID)
	err = s.API.App.Database.Insert(new(model.Post), post, "id")
	s.Nil(err)

	postDetail := new(model.PostDetail)
	postDetail.Title = "Post Title Slug"
	postDetail.Content = "Post Content"

	response := s.JSON(Post, fmt.Sprintf("/api/v1/post/%d/detail", post.ID), postDetail)

	s.Equal(response.Status, fasthttp.StatusUnprocessableEntity)
	data, _ := response.Error.Errors.(map[string]interface{})
	s.Equal(data["title"], "has been already taken")

	defaultLogger.LogInfo("Should be 422 error create a post detail with valid " +
		"params if slug has been already taken")
}

func (s PostDetailControllerTest) TearDownSuite() {
	TearDownSuite(s.Suite)
}

func Test_PostDetailController(t *testing.T) {
	s := PostDetailControllerTest{NewSuite()}
	Run(t, s)
}
