package api

import (
	"context"
	"fmt"
	"forgolang_forum/database/model"
	"github.com/valyala/fasthttp"
	"strconv"
	"testing"
)

type PostCategoryAssignmentControllerTest struct {
	*Suite
}

func (s PostCategoryAssignmentControllerTest) SetupSuite() {
	SetupSuite(s.Suite)
	UserAuth(s.Suite)
}

func (s PostCategoryAssignmentControllerTest) Test_CreatePostCategoryAssignmentWithValidParams() {
	post := model.NewPost(s.Auth.User.ID)
	err := s.API.GetDB().Insert(new(model.Post), post, "id")
	s.Nil(err)

	postDep := new(model.PostDEP)
	postDep.ID = post.ID
	_, err = s.API.App.ElasticClient.Index().
		Index("posts").
		Id(strconv.FormatInt(post.ID, 10)).
		BodyJson(postDep).
		Do(context.TODO())
	s.Nil(err)
	fmt.Println(post.ID)

	category := model.NewCategory()
	category.Title = "category-1"
	category.Slug = "category-1"
	err = s.API.GetDB().Insert(new(model.Category), category, "id")
	s.Nil(err)

	categoryAssignment := new(model.PostCategoryAssignment)
	categoryAssignment.CategoryID = category.ID

	response := s.JSON(Post, fmt.Sprintf("/api/v1/post/%d/category_assignment", post.ID),
		categoryAssignment)

	s.Equal(response.Status, fasthttp.StatusCreated)

	defaultLogger.LogInfo("Create post category assignment with valid params")
}

func (s PostCategoryAssignmentControllerTest) TearDownSuite() {
	TearDownSuite(s.Suite)
}

func Test_PostCategoryAssignmentController(t *testing.T) {
	s := PostCategoryAssignmentControllerTest{NewSuite()}
	Run(t, s)
}
