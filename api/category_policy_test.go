package api

import (
	"fmt"
	"forgolang_forum/database/model"
	"github.com/gosimple/slug"
	"github.com/valyala/fasthttp"
	"testing"
)

type CategoryPolicyTest struct {
	*Suite
}

func (s CategoryPolicyTest) SetupSuite() {
	SetupSuite(s.Suite)
}

func (s CategoryPolicyTest) Test_Should_403Error_CreateCategoryWithValidParamsAndUserRole() {
	UserAuth(s.Suite, "user")

	category := model.NewCategory()
	category.Title = "Category"
	category.Description = "Description"

	resp := s.JSON(Post, "/api/v1/category", category)

	s.Equal(resp.Status, fasthttp.StatusForbidden)

	defaultLogger.LogInfo("Should be 403 error create category with valid params " +
		"and user role")
}

func (s CategoryPolicyTest) Test_Should_UpdateCategoryWithValidParamsAndModeratorRole() {
	UserAuth(s.Suite, "moderator")

	category := model.NewCategory()
	category.Title = "Category / Edit"
	category.Description = "Description"
	category.Slug = slug.Make(category.Title)
	err := s.API.App.Database.Insert(new(model.Category), category, "id")
	s.Nil(err)

	categoryR := model.NewCategory()
	categoryR.Title = "Title / Edit"
	categoryR.Description = "Description / Edit"

	resp := s.JSON(Put, fmt.Sprintf("/api/v1/category/%d", category.ID), categoryR)

	s.Equal(resp.Status, fasthttp.StatusOK)

	defaultLogger.LogInfo("Update category with given identifier and valid params " +
		"and moderator role")
}


func (s CategoryPolicyTest) Test_Should_403Error_UpdateCategoryWithValidParamsAndUserRole() {
	UserAuth(s.Suite, "user")

	category := model.NewCategory()
	category.Title = "Category"
	category.Description = "Description"
	category.Slug = slug.Make(category.Title)
	err := s.API.App.Database.Insert(new(model.Category), category, "id")
	s.Nil(err)

	category = model.NewCategory()
	category.Title = "Title"

	resp := s.JSON(Put, fmt.Sprintf("/api/v1/category/%d", category.ID), category)

	s.Equal(resp.Status, fasthttp.StatusForbidden)

	defaultLogger.LogInfo("Should be 403 error update category with valid params " +
		"and user role")
}

func (s CategoryPolicyTest) Test_Should_403Error_DeleteCategoryWithGivenIdentifierAndUserRole() {
	UserAuth(s.Suite, "moderator")

	category := model.NewCategory()
	category.Title = "Category 2"
	category.Description = "Description"
	category.Slug = slug.Make(category.Title)
	err := s.API.App.Database.Insert(new(model.Category), category, "id")
	s.Nil(err)

	resp := s.JSON(Delete, fmt.Sprintf("/api/v1/category/%d", category.ID), nil)

	s.Equal(resp.Status, fasthttp.StatusForbidden)

	defaultLogger.LogInfo("Should be 403 error update category with valid params " +
		"and user role")
}

func (s CategoryPolicyTest) TearDownSuite() {
	TearDownSuite(s.Suite)
}

func Test_CategoryPolicy(t *testing.T) {
	s := CategoryPolicyTest{NewSuite()}
	Run(t, s)
}
