package api

import (
	"fmt"
	"forgolang_forum/cmn"
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

func (s PostControllerTest) Test_ListAllPosts() {
	s.API.App.Cache.Del(cmn.GetRedisKey("post", "count"))

	for i := 0; i < 50; i++ {
		post := model.NewPost(s.Auth.User.ID)
		err := s.API.App.Database.Insert(new(model.Post), post, "id")
		s.Nil(err)
		if err != nil {
			break
		}
		postDetail := model.NewPostDetail(post.ID, s.Auth.User.ID)
		postDetail.Title = "Post"
		postDetail.Description.SetValid("Post Detail")
		postDetail.Content = "Post Context"
		err = s.API.App.Database.Insert(new(model.PostDetail), postDetail, "id")
		s.Nil(err)
		if err != nil {
			break
		}
	}

	response := s.JSON(Get, "/api/v1/post", nil)

	s.Equal(response.Status, fasthttp.StatusOK)
	s.Greater(response.Success.TotalCount, int64(49))
	data, _ := response.Success.Data.([]interface{})
	s.Equal(len(data), 40)

	defaultLogger.LogInfo("List all posts")
}

func (s PostControllerTest) Test_ListAllPostsWithPaginationParams() {
	s.API.App.Cache.Del(cmn.GetRedisKey("post", "count"))

	for i := 0; i < 50; i++ {
		post := model.NewPost(s.Auth.User.ID)
		err := s.API.App.Database.Insert(new(model.Post), post, "id")
		s.Nil(err)
		if err != nil {
			break
		}
		postDetail := model.NewPostDetail(post.ID, s.Auth.User.ID)
		postDetail.Title = "Post"
		postDetail.Description.SetValid("Post Detail")
		postDetail.Content = "Post Context"
		err = s.API.App.Database.Insert(new(model.PostDetail), postDetail, "id")
		s.Nil(err)
		if err != nil {
			break
		}
	}

	response := s.JSON(Get, fmt.Sprintf("/api/v1/post?limit=20&offset=10"), nil)

	s.Equal(response.Status, fasthttp.StatusOK)
	s.Greater(response.Success.TotalCount, int64(49))
	data, _ := response.Success.Data.([]interface{})
	s.Equal(len(data), 20)

	defaultLogger.LogInfo("List all posts with pagination params")
}

func (s PostControllerTest) Test_ShowPostWithGivenSlug() {
	post := model.NewPost(s.Auth.User.ID)
	err := s.API.App.Database.Insert(new(model.Post), post, "id")
	s.Nil(err)

	postDetail := model.NewPostDetail(post.ID, s.Auth.User.ID)
	postDetail.Title = "Post"
	postDetail.Description.SetValid("Post Detail")
	postDetail.Content = "Post Context"
	err = s.API.App.Database.Insert(new(model.PostDetail), postDetail, "id")
	s.Nil(err)

	postSlug := model.NewPostSlug(post.ID, s.Auth.User.ID)
	postSlug.Slug = "slug-1"
	err = s.API.App.Database.Insert(new(model.PostSlug), postSlug, "id")
	s.Nil(err)

	response := s.JSON(Get, fmt.Sprintf("/api/v1/post/%s", postSlug.Slug), nil)

	s.Equal(response.Status, fasthttp.StatusOK)
	data, _ := response.Success.Data.(map[string]interface{})
	s.Equal(data["id"], float64(post.ID))
	s.Equal(data["author_id"], float64(s.Auth.User.ID))
	s.Equal(data["author_username"], s.Auth.User.Username)
	s.Equal(data["slug"], postSlug.Slug)
	s.Equal(data["title"], "Post")
	s.Equal(data["description"], "Post Detail")
	s.Equal(data["content"], "Post Context")
	s.NotNil(data["inserted_at"])

	defaultLogger.LogInfo("Show a post with given slug")
}

func (s PostControllerTest) Test_ShowPostWithGivenIdentifier() {
	post := model.NewPost(s.Auth.User.ID)
	err := s.API.App.Database.Insert(new(model.Post), post, "id")
	s.Nil(err)

	postDetail := model.NewPostDetail(post.ID, s.Auth.User.ID)
	postDetail.Title = "Post"
	postDetail.Description.SetValid("Post Detail")
	postDetail.Content = "Post Context"
	err = s.API.App.Database.Insert(new(model.PostDetail), postDetail, "id")
	s.Nil(err)

	postSlug := model.NewPostSlug(post.ID, s.Auth.User.ID)
	postSlug.Slug = "slug-2"
	err = s.API.App.Database.Insert(new(model.PostSlug), postSlug, "id")
	s.Nil(err)

	response := s.JSON(Get, fmt.Sprintf("/api/v1/post/%d", post.ID), nil)

	s.Equal(response.Status, fasthttp.StatusOK)
	data, _ := response.Success.Data.(map[string]interface{})
	s.Equal(data["id"], float64(post.ID))
	s.Equal(data["author_id"], float64(s.Auth.User.ID))
	s.Equal(data["author_username"], s.Auth.User.Username)
	s.Equal(data["slug"], postSlug.Slug)
	s.Equal(data["title"], "Post")
	s.Equal(data["description"], "Post Detail")
	s.Equal(data["content"], "Post Context")
	s.NotNil(data["inserted_at"])

	defaultLogger.LogInfo("Show a post with given identifier")
}

func (s PostControllerTest) Test_Should_404Err_ShowPostWithGivenIdentifierIfNotExists() {
	response := s.JSON(Get, "/api/v1/post/999999", nil)

	s.Equal(response.Status, fasthttp.StatusNotFound)

	defaultLogger.LogInfo("Should be 404 error show a post with given identifier" +
		" if does not exists")
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
	err := s.API.App.Database.Insert(new(model.Post), post, "id")
	s.Nil(err)

	postSlug := model.NewPostSlug(post.ID, s.Auth.User.ID)
	postSlug.Slug = slug.Make("Post title slug")
	err = s.API.App.Database.Insert(new(model.PostSlug), postSlug, "id")
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
	err := s.API.App.Database.Insert(new(model.Post), post, "id")
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
