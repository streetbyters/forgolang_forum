package api

import (
	"encoding/json"
	"fmt"
	"forgolang_forum/cmn"
	"forgolang_forum/database/model"
	"github.com/gosimple/slug"
	"github.com/valyala/fasthttp"
	"testing"
	"time"
)

type CategoryControllerTest struct {
	*Suite
}

func (s CategoryControllerTest) SetupSuite() {
	SetupSuite(s.Suite)
	UserAuth(s.Suite)
}

func (s CategoryControllerTest) Test_ListAllCategories() {
	for i := 0; i < 50; i++ {
		category := model.NewCategory()
		category.Title = fmt.Sprintf("Category %d", i)
		category.Description.SetValid("Category Description")
		category.Slug = slug.Make(category.Title)
		err := s.API.GetDB().Insert(model.NewCategory(),
			category,
			"id", "inserted_at", "updated_at")
		s.Nil(err)
	}

	resp := s.JSON(Get, "/api/v1/category", nil)

	s.Equal(resp.Status, fasthttp.StatusOK)
	s.Greater(resp.Success.TotalCount, int64(49))

	defaultLogger.LogInfo("List all categories")
}

func (s CategoryControllerTest) Test_ListAllCategoriesAndCategoryLanguages() {
	for i := 0; i < 50; i++ {
		category := model.NewCategory()
		category.Title = fmt.Sprintf("Category 2 %d", i)
		category.Description.SetValid("Category Description")
		category.Slug = slug.Make(category.Title)
		err := s.API.GetDB().Insert(model.NewCategory(),
			category,
			"id", "inserted_at", "updated_at")
		s.Nil(err)
		if err != nil {
			break
		}

		for j, l := range s.API.Languages {
			categoryLanguage := model.NewCategoryLanguage(category.ID)
			categoryLanguage.LanguageID = l.ID
			categoryLanguage.Title = fmt.Sprintf("Category Language %d %d", i, j)
			categoryLanguage.Slug = slug.Make(categoryLanguage.Title)
			err := s.API.GetDB().Insert(new(model.CategoryLanguage), categoryLanguage, "id")
			s.Nil(err)
			if err != nil {
				break
			}
		}
	}

	resp := s.JSON(Get, "/api/v1/category", nil)

	s.Equal(resp.Status, fasthttp.StatusOK)
	s.Greater(resp.Success.TotalCount, int64(49))

	defaultLogger.LogInfo("List all categories and category languages")
}

func (s CategoryControllerTest) Test_ListAllCachedCategories() {
	for i := 0; i < 50; i++ {
		category := model.NewCategory()
		category.Title = fmt.Sprintf("Category Cache %d", i)
		category.Description.SetValid("Category Description")
		category.Slug = slug.Make(category.Title)
		err := s.API.GetDB().Insert(model.NewCategory(),
			category,
			"id", "inserted_at", "updated_at")
		s.Nil(err)
		s.API.App.Cache.SAdd(cmn.GetRedisKey("category", "all"), category.ToJSON())
	}

	resp := s.JSON(Get, "/api/v1/category", nil)

	s.Equal(resp.Status, fasthttp.StatusOK)
	s.Greater(resp.Success.TotalCount, int64(49))

	defaultLogger.LogInfo("List all categories")
}

func (s CategoryControllerTest) Test_ListAllCachedCategoriesAndCategoryLanguages() {
	for i := 0; i < 50; i++ {
		category := model.NewCategory()
		category.Title = fmt.Sprintf("Category Cache 2 %d", i)
		category.Description.SetValid("Category Description")
		category.Slug = slug.Make(category.Title)
		err := s.API.GetDB().Insert(model.NewCategory(),
			category,
			"id", "inserted_at", "updated_at")
		s.Nil(err)

		for j, l := range s.API.Languages {
			categoryLanguage := model.NewCategoryLanguage(category.ID)
			categoryLanguage.LanguageID = l.ID
			categoryLanguage.Title = fmt.Sprintf("Category Language 2 %d %d", i, j)
			categoryLanguage.Slug = slug.Make(categoryLanguage.Title)
			err := s.API.GetDB().Insert(new(model.CategoryLanguage), categoryLanguage, "id")
			s.Nil(err)
			if err != nil {
				break
			}
			s.API.GetCache().SAdd(fmt.Sprintf("%s:%d",
				cmn.GetRedisKey("category", "languages"),
				category.ID),
				categoryLanguage.ToJSON())
		}

		s.API.App.Cache.SAdd(cmn.GetRedisKey("category", "all"), category.ToJSON())
	}

	resp := s.JSON(Get, "/api/v1/category", nil)

	s.Equal(resp.Status, fasthttp.StatusOK)
	s.Greater(resp.Success.TotalCount, int64(49))

	defaultLogger.LogInfo("List all cached categories and category languages")
}

func (s CategoryControllerTest) Test_ShowCategoryWithGivenIdentifier() {
	category := model.NewCategory()
	category.Title = "Show Category"
	category.Description.SetValid("Show Category Description")
	category.Slug = slug.Make("Show category")
	err := s.API.GetDB().Insert(model.NewCategory(),
		category,
		"id", "inserted_at", "updated_at")
	s.Nil(err)

	resp := s.JSON(Get, fmt.Sprintf("/api/v1/category/%d", category.ID), nil)

	s.Equal(resp.Status, fasthttp.StatusOK)

	data, _ := resp.Success.Data.(map[string]interface{})
	s.Equal(data["id"], float64(category.ID))
	s.Equal(data["title"], "Show Category")
	s.Equal(data["description"], "Show Category Description")
	s.Equal(data["slug"], "show-category")
	s.Equal(data["inserted_at"], category.InsertedAt.UTC().Format(time.RFC3339Nano))
	s.Equal(data["updated_at"], category.UpdatedAt.UTC().Format(time.RFC3339Nano))

	defaultLogger.LogInfo("Show a category with given identifier")
}

func (s CategoryControllerTest) Test_ShowCategoryAndCategoryLanguagesWithGivenIdentifier() {
	category := model.NewCategory()
	category.Title = "Show Category Language"
	category.Description.SetValid("Show Category Description")
	category.Slug = slug.Make(category.Title)
	err := s.API.GetDB().Insert(model.NewCategory(),
		category,
		"id", "inserted_at", "updated_at")
	s.Nil(err)

	for i, l := range s.API.Languages {
		categoryLanguage := model.NewCategoryLanguage(category.ID)
		categoryLanguage.LanguageID = l.ID
		categoryLanguage.Title = fmt.Sprintf("Show Category Language %d", i)
		categoryLanguage.Slug = slug.Make(categoryLanguage.Title)
		err := s.API.GetDB().Insert(new(model.CategoryLanguage), categoryLanguage, "id")
		s.Nil(err)
		if err != nil {
			break
		}
	}

	resp := s.JSON(Get, fmt.Sprintf("/api/v1/category/%d", category.ID), nil)

	s.Equal(resp.Status, fasthttp.StatusOK)

	data, _ := resp.Success.Data.(map[string]interface{})
	s.Equal(data["id"], float64(category.ID))
	s.Equal(data["title"], "Show Category Language")
	s.Equal(data["description"], "Show Category Description")
	s.Equal(data["slug"], "show-category-language")
	s.Equal(data["inserted_at"], category.InsertedAt.UTC().Format(time.RFC3339Nano))
	s.Equal(data["updated_at"], category.UpdatedAt.UTC().Format(time.RFC3339Nano))
	s.Equal(len(data["languages"].([]interface{})), 2)

	defaultLogger.LogInfo("Show a category and category languages " +
		"with given identifier")
}

func (s CategoryControllerTest) Test_ShowCachedCategoryWithGivenIdentifier() {
	category := model.NewCategory()
	category.Title = "Show Category"
	category.Description.SetValid("Show Category Description")
	category.Slug = slug.Make("Show category 2")
	err := s.API.GetDB().Insert(model.NewCategory(),
		category,
		"id", "inserted_at", "updated_at")
	s.Nil(err)

	s.API.App.Cache.Set(fmt.Sprintf("%s:%d",
		cmn.GetRedisKey("category", "one"),
		category.ID), category.ToJSON(), 0)

	resp := s.JSON(Get, fmt.Sprintf("/api/v1/category/%d", category.ID), nil)

	s.Equal(resp.Status, fasthttp.StatusOK)

	data, _ := resp.Success.Data.(map[string]interface{})
	s.Equal(data["id"], float64(category.ID))
	s.Equal(data["title"], "Show Category")
	s.Equal(data["description"], "Show Category Description")
	s.Equal(data["slug"], "show-category-2")
	s.Equal(data["inserted_at"], category.InsertedAt.Format(time.RFC3339Nano))
	s.Equal(data["updated_at"], category.UpdatedAt.Format(time.RFC3339Nano))

	defaultLogger.LogInfo("Show a cached category with given identifier")
}

func (s CategoryControllerTest) Test_ShowCachedCategoryAndCategoryLanguageWithGivenIdentifier() {
	category := model.NewCategory()
	category.Title = "Show Category 2 1"
	category.Description.SetValid("Show Category Description")
	category.Slug = slug.Make("Show category 2 1")
	err := s.API.GetDB().Insert(model.NewCategory(),
		category,
		"id", "inserted_at", "updated_at")
	s.Nil(err)

	s.API.App.Cache.Set(fmt.Sprintf("%s:%d",
		cmn.GetRedisKey("category", "one"),
		category.ID), category.ToJSON(), 0)

	for j, l := range s.API.Languages {
		categoryLanguage := model.NewCategoryLanguage(category.ID)
		categoryLanguage.LanguageID = l.ID
		categoryLanguage.Title = fmt.Sprintf("Show Category Language 2 1 %d", j)
		categoryLanguage.Slug = slug.Make(categoryLanguage.Title)
		err := s.API.GetDB().Insert(new(model.CategoryLanguage), categoryLanguage, "id")
		s.Nil(err)
		if err != nil {
			break
		}
		s.API.GetCache().SAdd(fmt.Sprintf("%s:%d",
			cmn.GetRedisKey("category", "languages"),
			category.ID),
			categoryLanguage.ToJSON())
	}

	resp := s.JSON(Get, fmt.Sprintf("/api/v1/category/%d", category.ID), nil)

	s.Equal(resp.Status, fasthttp.StatusOK)

	data, _ := resp.Success.Data.(map[string]interface{})
	s.Equal(data["id"], float64(category.ID))
	s.Equal(data["title"], "Show Category 2 1")
	s.Equal(data["description"], "Show Category Description")
	s.Equal(data["slug"], "show-category-2-1")
	s.Equal(data["inserted_at"], category.InsertedAt.Format(time.RFC3339Nano))
	s.Equal(data["updated_at"], category.UpdatedAt.Format(time.RFC3339Nano))
	s.Equal(len(data["languages"].([]interface{})), 2)

	defaultLogger.LogInfo("Show a cached category and category language " +
		" with given identifier")
}

func (s CategoryControllerTest) Test_Should_404Err_ShowCategoryWithGivenIdentifierIfNotExists() {
	resp := s.JSON(Get, "/api/v1/category/999999", nil)

	s.Equal(resp.Status, fasthttp.StatusNotFound)

	defaultLogger.LogInfo("Should be 404 error show a category with given identifier " +
		"if does not exists")
}

func (s CategoryControllerTest) Test_CreateCategoryWithValidParams() {
	category := model.NewCategory()
	category.Title = "Create Category"
	category.Description.SetValid("Category Description")
	category.Slug = "sss"

	resp := s.JSON(Post, "/api/v1/category", category)

	s.Equal(resp.Status, fasthttp.StatusCreated)

	data, _ := resp.Success.Data.(map[string]interface{})

	s.Greater(data["id"], float64(0))
	s.Equal(data["title"], "Create Category")
	s.Equal(data["description"], "Category Description")
	s.Equal(data["slug"], "create-category")

	cachedCategory := new(model.Category)
	var _s string
	err := s.API.App.Cache.Get(fmt.Sprintf("%s:%d",
		cmn.GetRedisKey("category", "one"),
		int64(data["id"].(float64)))).Scan(&_s)
	s.Nil(err)

	json.Unmarshal([]byte(_s), &cachedCategory)

	s.Equal(cachedCategory.ID, int64(data["id"].(float64)))
	s.Equal(cachedCategory.Title, data["title"])
	s.Equal(cachedCategory.Description.String, data["description"])
	s.Equal(cachedCategory.Slug, data["slug"])
	s.Equal(cachedCategory.InsertedAt.Format(time.RFC3339Nano), data["inserted_at"])
	s.Equal(cachedCategory.UpdatedAt.Format(time.RFC3339Nano), data["updated_at"])

	defaultLogger.LogInfo("Create category with valid params")
}

func (s CategoryControllerTest) Test_Shoul_422Err_CreateCategoryWithValidParams() {
	category := model.NewCategory()
	category.Description.SetValid("Category Description")

	resp := s.JSON(Post, "/api/v1/category", category)

	s.Equal(resp.Status, fasthttp.StatusUnprocessableEntity)

	defaultLogger.LogInfo("Should be 422 error create category with invalid params")
}

func (s CategoryControllerTest) Test_Should_422Err_CreateCategoryWithValidParamsIfSlugNotUnique() {
	category := model.NewCategory()
	category.Title = "Create Category 2"
	category.Description.SetValid("Category Description")
	category.Slug = slug.Make(category.Title)
	err := s.API.GetDB().Insert(new(model.Category), category, "id")
	s.Nil(err)

	category = model.NewCategory()
	category.Title = "Create Category 2"
	category.Description.SetValid("Category Description")
	category.Slug = "sss"

	resp := s.JSON(Post, "/api/v1/category", category)

	s.Equal(resp.Status, fasthttp.StatusUnprocessableEntity)

	defaultLogger.LogInfo("Should be 422 error create category with valid params " +
		"if slug has been already taken")
}

func (s CategoryControllerTest) Test_UpdateCategoryWithGivenIdentifierAndValidParams() {
	category := model.NewCategory()
	category.Title = "Update Category"
	category.Description.SetValid("Category Description")
	category.Slug = slug.Make(category.Title)
	err := s.API.GetDB().Insert(new(model.Category),
		category,
		"id", "inserted_at", "updated_at")
	s.Nil(err)

	categoryR := model.NewCategory()
	categoryR.Title = "Update Category / Edit"
	categoryR.Description.SetValid("Category Description")

	resp := s.JSON(Put, fmt.Sprintf("/api/v1/category/%d", category.ID), categoryR)

	s.Equal(resp.Status, fasthttp.StatusOK)

	data, _ := resp.Success.Data.(map[string]interface{})

	s.Equal(data["id"], float64(category.ID))
	s.Equal(data["title"], "Update Category / Edit")
	s.Equal(data["description"], "Category Description")
	s.Equal(data["slug"], "update-category-edit")
	s.Equal(data["inserted_at"], category.InsertedAt.Format(time.RFC3339Nano))
	s.NotEqual(data["updated_at"], category.UpdatedAt.Format(time.RFC3339Nano))

	cachedCategory := new(model.Category)
	var _s string
	err = s.API.App.Cache.Get(fmt.Sprintf("%s:%d",
		cmn.GetRedisKey("category", "one"),
		int64(data["id"].(float64)))).Scan(&_s)
	s.Nil(err)

	json.Unmarshal([]byte(_s), &cachedCategory)

	s.Equal(cachedCategory.ID, int64(data["id"].(float64)))
	s.Equal(cachedCategory.Title, data["title"])
	s.Equal(cachedCategory.Description.String, data["description"])
	s.Equal(cachedCategory.Slug, data["slug"])
	s.Equal(cachedCategory.InsertedAt.Format(time.RFC3339Nano), data["inserted_at"])
	s.Equal(cachedCategory.UpdatedAt.Format(time.RFC3339Nano), data["updated_at"])

	defaultLogger.LogInfo("Update a category with given identifier and valid params")
}

func (s CategoryControllerTest) Test_Should_422Error_UpdateCategoryWithGivenIdentifierAndInvalidParams() {
	category := model.NewCategory()
	category.Title = "Update Category 2"
	category.Description.SetValid("Category Description")
	category.Slug = slug.Make(category.Title)
	err := s.API.GetDB().Insert(new(model.Category),
		category,
		"id", "inserted_at", "updated_at")
	s.Nil(err)

	categoryR := model.NewCategory()
	categoryR.Description.SetValid("Category Description")

	resp := s.JSON(Put, fmt.Sprintf("/api/v1/category/%d", category.ID), categoryR)

	s.Equal(resp.Status, fasthttp.StatusUnprocessableEntity)

	defaultLogger.LogInfo("Should be 422 error update a category with given " +
		"identifier and invalid params")
}

func (s CategoryControllerTest) Test_Should_404Error_UpdateCategoryWithGivenIdentifierAndValidParamsIfNotExists() {
	categoryR := model.NewCategory()
	categoryR.Title = "Update Category 4 / Edit"
	categoryR.Description.SetValid("Category Description")

	resp := s.JSON(Put, "/api/v1/category/9999999", categoryR)

	s.Equal(resp.Status, fasthttp.StatusNotFound)

	defaultLogger.LogInfo("Should be 404 error update a category with given " +
		"identifier and valid params if does not exists")
}

func (s CategoryControllerTest) Test_Should_422Error_UpdateCategoryWithGivenIdentifierAndValidParamsIfSlugNotUnique() {
	category := model.NewCategory()
	category.Title = "Update Category 4 / Edit"
	category.Description.SetValid("Category Description")
	category.Slug = slug.Make(category.Title)
	err := s.API.GetDB().Insert(new(model.Category),
		category,
		"id", "inserted_at", "updated_at")
	s.Nil(err)

	category = model.NewCategory()
	category.Title = "Update Category 3"
	category.Description.SetValid("Category Description")
	category.Slug = slug.Make(category.Title)
	err = s.API.GetDB().Insert(new(model.Category),
		category,
		"id", "inserted_at", "updated_at")
	s.Nil(err)

	categoryR := model.NewCategory()
	categoryR.Title = "Update Category 4 / Edit"
	categoryR.Description.SetValid("Category Description")

	resp := s.JSON(Put, fmt.Sprintf("/api/v1/category/%d", category.ID), categoryR)

	s.Equal(resp.Status, fasthttp.StatusUnprocessableEntity)

	defaultLogger.LogInfo("Should be 422 error update a category with given " +
		"identifier and valid params if slug has been already taken")
}

func (s CategoryControllerTest) Test_DeleteCategoryWithGivenIdentifier() {
	category := model.NewCategory()
	category.Title = "Delete Category"
	category.Description.SetValid("Category Description")
	category.Slug = slug.Make(category.Title)
	err := s.API.GetDB().Insert(new(model.Category),
		category,
		"id", "inserted_at", "updated_at")
	s.Nil(err)

	s.API.App.Cache.Set(fmt.Sprintf("%s:%d",
		cmn.GetRedisKey("category", "one"),
		category.ID), category.ToJSON(), 0)

	resp := s.JSON(Delete, fmt.Sprintf("/api/v1/category/%d", category.ID), nil)

	s.Equal(resp.Status, fasthttp.StatusNoContent)

	err = s.API.App.Cache.Get(fmt.Sprintf("%s:%d",
		cmn.GetRedisKey("category", "one"),
		category.ID)).Err()
	s.Equal(err.Error(), "redis: nil")

	defaultLogger.LogInfo("Delete a category with given identifier")
}

func (s CategoryControllerTest) Test_Should_404Err_DeleteCategoryWithGivenIdentifierIfNotExists() {
	resp := s.JSON(Delete, "/api/v1/category/999999999", nil)

	s.Equal(resp.Status, fasthttp.StatusNotFound)

	defaultLogger.LogInfo("Should be 404 error delete a category with given " +
		"identifier if does not exists")
}

func (s CategoryControllerTest) TearDownSuite() {
	TearDownSuite(s.Suite)
}

func Test_CategoryController(t *testing.T) {
	s := CategoryControllerTest{NewSuite()}
	Run(t, s)
}
