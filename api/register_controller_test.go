package api

import (
	model2 "forgolang_forum/database/model"
	"forgolang_forum/model"
	"github.com/valyala/fasthttp"
	"testing"
)

type RegisterControllerTest struct {
	*Suite
}

func (s RegisterControllerTest) SetupSuite() {
	SetupSuite(s.Suite)
}

func (s RegisterControllerTest) Test_PostRegisterWithValidParams() {
	register := new(model.RegisterRequest)
	register.Username = "register-user"
	register.Email = "akdilsiz@tecpor.com"
	register.Password = "123456"

	response := s.JSON(Post, "/api/v1/auth/register", register)

	s.Equal(response.Status, fasthttp.StatusCreated)

	data, _ := response.Success.Data.(map[string]interface{})

	s.Greater(data["user_id"], float64(0))
	s.Equal(data["state"], "wait_for_confirmation")
	s.NotNil(data["inserted_at"])

	defaultLogger.LogInfo("Successfully post register with valid params")
}

func (s RegisterControllerTest) Test_Should_422Err_PostRegisterWithInalidParams() {
	register := new(model.RegisterRequest)
	register.Username = "register-user"
	register.Email = "akdilsiz"

	response := s.JSON(Post, "/api/v1/auth/register", register)

	s.Equal(response.Status, fasthttp.StatusUnprocessableEntity)

	defaultLogger.LogInfo("Should be 422 error post register with invalid params")
}

func (s RegisterControllerTest) Test_Should_422Err_PostRegisterWithValidParamsIfUsernameNotUnique() {
	pwd := "12345"
	user := model2.NewUser(&pwd)
	user.Email = "akdilsiz-2@tecpor.com"
	user.Username = "akdilsiz"
	user.IsActive = true
	err := s.API.GetDB().Insert(new(model2.User), user, "id")
	s.Nil(err)

	register := new(model.RegisterRequest)
	register.Username = "akdilsiz"
	register.Email = "register-user@tecpor.com"
	register.Password = "123456"

	response := s.JSON(Post, "/api/v1/auth/register", register)

	s.Equal(response.Status, fasthttp.StatusUnprocessableEntity)

	defaultLogger.LogInfo("Should be 422 error post register with valid params" +
		" if username has been already taken")
}

func (s RegisterControllerTest) TearDownSuite() {
	TearDownSuite(s.Suite)
}

func Test_RegisterController(t *testing.T) {
	s := RegisterControllerTest{NewSuite()}
	Run(t, s)
}
