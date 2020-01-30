package api

import (
	"fmt"
	"forgolang_forum/cmn"
	"forgolang_forum/database/model"
	"github.com/gosimple/slug"
	"github.com/valyala/fasthttp"
	"testing"
)

type CategoryPostControllerTest struct {
	*Suite
}

func (s CategoryPostControllerTest) SetupSuite() {
	SetupSuite(s.Suite)
	UserAuth(s.Suite)
}

func (s CategoryPostControllerTest) Test_ListAllPosts() {
	s.API.App.Cache.Del(cmn.GetRedisKey("post", "count"))

	category := model.NewCategory()
	category.Title = "Category"
	category.Slug = slug.Make(category.Title)
	err := s.API.App.Database.Insert(new(model.Category), category, "id")
	s.Nil(err)
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

		postCategoryAssignment := model.NewPostCategoryAssignment(post.ID, category.ID, s.Auth.User.ID)
		err = s.API.App.Database.Insert(new(model.PostCategoryAssignment), postCategoryAssignment, "id")
		s.Nil(err)
		if err != nil {
			break
		}
	}

	response := s.JSON(Get, fmt.Sprintf("/api/v1/category/%d/post", category.ID), nil)

	s.Equal(response.Status, fasthttp.StatusOK)
	s.Greater(response.Success.TotalCount, int64(49))
	data, _ := response.Success.Data.([]interface{})
	s.Equal(len(data), 40)

	defaultLogger.LogInfo("List all posts")
}

func (s CategoryPostControllerTest) Test_ListAllPostsWithPaginationParams() {
	s.API.App.Cache.Del(cmn.GetRedisKey("post", "count"))
	category := model.NewCategory()
	category.Title = "Category 2"
	category.Slug = slug.Make(category.Title)
	err := s.API.App.Database.Insert(new(model.Category), category, "id")
	s.Nil(err)
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

		postCategoryAssignment := model.NewPostCategoryAssignment(post.ID, category.ID, s.Auth.User.ID)
		err = s.API.App.Database.Insert(new(model.PostCategoryAssignment), postCategoryAssignment, "id")
		s.Nil(err)
		if err != nil {
			break
		}
	}

	response := s.JSON(Get, fmt.Sprintf("/api/v1/category/%s/post?limit=20&offset=10", category.Slug), nil)

	s.Equal(response.Status, fasthttp.StatusOK)
	s.Greater(response.Success.TotalCount, int64(49))
	data, _ := response.Success.Data.([]interface{})
	s.Equal(len(data), 20)

	defaultLogger.LogInfo("List all posts with pagination params")
}

func (s CategoryPostControllerTest) Test_ShowPostWithGivenSlug() {
	category := model.NewCategory()
	category.Title = "Category 3"
	category.Slug = slug.Make(category.Title)
	err := s.API.App.Database.Insert(new(model.Category), category, "id")
	s.Nil(err)

	post := model.NewPost(s.Auth.User.ID)
	err = s.API.App.Database.Insert(new(model.Post), post, "id")
	s.Nil(err)

	postCategoryAssignment := model.NewPostCategoryAssignment(post.ID, category.ID, s.Auth.User.ID)
	err = s.API.App.Database.Insert(new(model.PostCategoryAssignment), postCategoryAssignment, "id")
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

	response := s.JSON(Get, fmt.Sprintf("/api/v1/category/%d/post/%s", category.ID, postSlug.Slug), nil)

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

func (s CategoryPostControllerTest) Test_ShowPostWithGivenIdentifier() {
	category := model.NewCategory()
	category.Title = "Category 4"
	category.Slug = slug.Make(category.Title)
	err := s.API.App.Database.Insert(new(model.Category), category, "id")
	s.Nil(err)

	post := model.NewPost(s.Auth.User.ID)
	err = s.API.App.Database.Insert(new(model.Post), post, "id")
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

	postCategoryAssignment := model.NewPostCategoryAssignment(post.ID, category.ID, s.Auth.User.ID)
	err = s.API.App.Database.Insert(new(model.PostCategoryAssignment), postCategoryAssignment, "id")
	s.Nil(err)

	response := s.JSON(Get, fmt.Sprintf("/api/v1/category/%d/post/%d", category.ID, post.ID), nil)

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

func (s CategoryPostControllerTest) Test_Should_404Err_ShowPostWithGivenIdentifierIfNotExists() {
	response := s.JSON(Get, "/api/v1/category/999999999/post/999999", nil)

	s.Equal(response.Status, fasthttp.StatusNotFound)

	defaultLogger.LogInfo("Should be 404 error show a post with given identifier" +
		" if does not exists")
}

func (s CategoryPostControllerTest) TearDownSuite() {
	TearDownSuite(s.Suite)
}

func Test_CategoryPostController(t *testing.T) {
	s := CategoryPostControllerTest{NewSuite()}
	Run(t, s)
}
