package api

import (
	"fmt"
	"forgolang_forum/cmn"
	"forgolang_forum/database/model"
	"github.com/gosimple/slug"
	"github.com/valyala/fasthttp"
	"testing"
)

type CategoryLanguageControllerTest struct {
	*Suite
}

func (s CategoryLanguageControllerTest) SetupSuite() {
	SetupSuite(s.Suite)
	UserAuth(s.Suite)
}

func (s CategoryLanguageControllerTest) Test_CreateCategoryLanguageWithValidParams() {
	category := model.NewCategory()
	category.Title = "Category Language"
	category.Description.SetValid("Category Description")
	category.Slug = slug.Make(category.Title)
	err := s.API.GetDB().Insert(new(model.Category), category, "id")
	s.Nil(err)

	categoryLanguage := new(model.CategoryLanguage)
	categoryLanguage.LanguageID = s.API.GetLanguage("tr-TR").ID
	categoryLanguage.Title = "Kategori 1"
	categoryLanguage.Slug = slug.Make(categoryLanguage.Title)
	categoryLanguage.Description.SetValid("Kategori Aciklama")

	response := s.JSON(Post, fmt.Sprintf("/api/v1/category/%d/language",
		category.ID), categoryLanguage)
	s.Equal(response.Status, fasthttp.StatusCreated)

	data, _ := response.Success.Data.(map[string]interface{})

	s.Greater(data["id"], float64(0))
	s.Equal(data["category_id"], float64(category.ID))
	s.Equal(data["language_id"], float64(s.API.GetLanguage("tr-TR").ID))
	s.Equal(data["title"], "Kategori 1")
	s.Equal(data["description"], "Kategori Aciklama")
	s.Equal(data["slug"], "kategori-1")

	categoryLanguage = new(model.CategoryLanguage)
	result := s.API.GetDB().QueryRowWithModel(fmt.Sprintf(`
		SELECT * FROM %s AS cl WHERE cl.id = $1
	`, categoryLanguage.TableName()),
		categoryLanguage,
		data["id"])
	s.Nil(result.Error)

	cache := s.API.GetCache().SIsMember(fmt.Sprintf("%s:%d",
		cmn.GetRedisKey("category", "languages"),
		category.ID),
		categoryLanguage.ToJSON())
	exists, err := cache.Result()
	s.Nil(err)
	s.Equal(exists, true)

	defaultLogger.LogInfo("Create category language with valid params")
}

func (s CategoryLanguageControllerTest) Test_Should_422Err_CreateCategoryLanguageWithInvalidParams() {
	category := model.NewCategory()
	category.Title = "Category Language 2"
	category.Description.SetValid("Category Description")
	category.Slug = slug.Make(category.Title)
	err := s.API.GetDB().Insert(new(model.Category), category, "id")
	s.Nil(err)

	categoryLanguage := new(model.CategoryLanguage)
	categoryLanguage.Title = "K"

	response := s.JSON(Post, fmt.Sprintf("/api/v1/category/%d/language",
		category.ID), categoryLanguage)

	s.Equal(response.Status, fasthttp.StatusUnprocessableEntity)

	defaultLogger.LogInfo("Should be 422 error create category language with " +
		"valid params")
}

func (s CategoryLanguageControllerTest) Test_Should_422Err_CreateCategoryLanguageWithValidParamsIfSlugNotUnique() {
	category := model.NewCategory()
	category.Title = "Category Language 3"
	category.Description.SetValid("Category Description")
	category.Slug = slug.Make(category.Title)
	err := s.API.GetDB().Insert(new(model.Category), category, "id")
	s.Nil(err)

	categoryLanguage := new(model.CategoryLanguage)
	categoryLanguage.CategoryID = category.ID
	categoryLanguage.LanguageID = s.API.GetLanguage("tr-TR").ID
	categoryLanguage.Title = "Kategori Slug"
	categoryLanguage.Slug = slug.Make(categoryLanguage.Title)
	categoryLanguage.Description.SetValid("Kategori Aciklama")

	err = s.API.GetDB().Insert(new(model.CategoryLanguage), categoryLanguage, "id")
	s.Nil(err)

	category = model.NewCategory()
	category.Title = "Category Language 4"
	category.Description.SetValid("Category Description")
	category.Slug = slug.Make(category.Title)
	err = s.API.GetDB().Insert(new(model.Category), category, "id")
	s.Nil(err)

	categoryLanguage = new(model.CategoryLanguage)
	categoryLanguage.LanguageID = s.API.GetLanguage("tr-TR").ID
	categoryLanguage.Title = "Kategori Slug"
	categoryLanguage.Slug = slug.Make(categoryLanguage.Title)
	categoryLanguage.Description.SetValid("Kategori Aciklama")

	response := s.JSON(Post, fmt.Sprintf("/api/v1/category/%d/language",
		category.ID), categoryLanguage)

	s.Equal(response.Status, fasthttp.StatusUnprocessableEntity)

	data, _ := response.Error.Errors.(map[string]interface{})
	s.Equal(data["slug"], "has been already taken")

	defaultLogger.LogInfo("Should be 422 error create category language with " +
		"valid params if slug not unique")
}

func (s CategoryLanguageControllerTest) Test_Should_422Err_CreateCategoryLanguageWithValidParamsIfRelationalError() {
	category := model.NewCategory()
	category.Title = "Category Language 5"
	category.Description.SetValid("Category Description")
	category.Slug = slug.Make(category.Title)
	err := s.API.GetDB().Insert(new(model.Category), category, "id")
	s.Nil(err)

	categoryLanguage := new(model.CategoryLanguage)
	categoryLanguage.LanguageID = int64(999999999)
	categoryLanguage.Title = "Kategori 5"
	categoryLanguage.Slug = slug.Make(categoryLanguage.Title)
	categoryLanguage.Description.SetValid("Kategori Aciklama")

	response := s.JSON(Post, fmt.Sprintf("/api/v1/category/%d/language",
		category.ID), categoryLanguage)

	s.Equal(response.Status, fasthttp.StatusUnprocessableEntity)

	data, _ := response.Error.Errors.(map[string]interface{})
	s.Equal(data["language_id"], "does not exists")

	defaultLogger.LogInfo("Should be 422 error create category language with " +
		"valid params if relational error")
}

func (s CategoryLanguageControllerTest) TearDownSuite() {
	TearDownSuite(s.Suite)
}

func Test_CategoryLanguageController(t *testing.T) {
	s := CategoryLanguageControllerTest{NewSuite()}
	Run(t, s)
}
