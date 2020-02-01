package api

import (
	"fmt"
	"forgolang_forum/database/model"
	"github.com/valyala/fasthttp"
	"testing"
)

type LogoutControllerTest struct {
	*Suite
}

func (s LogoutControllerTest) SetupSuite() {
	SetupSuite(s.Suite)
	UserAuth(s.Suite)
}

func (s LogoutControllerTest) Test_PostLogoutWithGivenIdentifier() {
	pwd := "123456"
	user := model.NewUser(&pwd)
	user.Email = "logout-user@mail.com"
	user.Username = "logout-user"
	user.IsActive = true
	err := s.API.GetDB().Insert(new(model.User), user, "id")
	s.Nil(err)

	passphrase := model.NewUserPassphrase(user.ID)
	err = s.API.GetDB().Insert(new(model.UserPassphrase), passphrase, "id")
	s.Nil(err)

	response := s.JSON(Post, fmt.Sprintf("/api/v1/user/%d/sign_out/%d",
		user.ID, passphrase.ID), nil)

	s.Equal(response.Status, fasthttp.StatusCreated)

	defaultLogger.LogInfo("Successfully post logout with given identifier")
}

func (s LogoutControllerTest) Test_Should_404Err_PostLogoutWithGivenIdentifierIfNotExists() {
	pwd := "123456"
	user := model.NewUser(&pwd)
	user.Email = "logout-user-2@mail.com"
	user.Username = "logout-user-2"
	user.IsActive = true
	err := s.API.GetDB().Insert(new(model.User), user, "id")
	s.Nil(err)

	response := s.JSON(Post, fmt.Sprintf("/api/v1/user/%d/sign_out/999999999",
		user.ID), nil)

	s.Equal(response.Status, fasthttp.StatusNotFound)

	defaultLogger.LogInfo("Should be 404 error post logout with given identifier " +
		"if passphrase does not exists")
}

func (s LogoutControllerTest) Test_Should_404Err_PostLogoutWithGivenIdentifierIfPassphraseAlreadyInvalidate() {
	pwd := "123456"
	user := model.NewUser(&pwd)
	user.Email = "logout-user-3@mail.com"
	user.Username = "logout-user-3"
	user.IsActive = true
	err := s.API.GetDB().Insert(new(model.User), user, "id")
	s.Nil(err)

	passphrase := model.NewUserPassphrase(user.ID)
	err = s.API.GetDB().Insert(new(model.UserPassphrase), passphrase, "id")
	s.Nil(err)

	passphraseInvalidation := model.NewUserPassphraseInvalidation()
	passphraseInvalidation.PassphraseID = passphrase.ID
	passphraseInvalidation.SourceUserID.SetValid(user.ID)
	err = s.API.GetDB().Insert(new(model.UserPassphraseInvalidation), passphraseInvalidation,
		"passphrase_id")
	s.Nil(err)

	response := s.JSON(Post, fmt.Sprintf("/api/v1/user/%d/sign_out/%d",
		user.ID, passphrase.ID), nil)

	s.Equal(response.Status, fasthttp.StatusNotFound)

	defaultLogger.LogInfo("Should be 404 error post logout with given identifier " +
		"if already invalidate")
}

func (s LogoutControllerTest) TearDownSuite() {
	TearDownSuite(s.Suite)
}

func Test_LogoutController(t *testing.T) {
	s := LogoutControllerTest{NewSuite()}
	Run(t, s)
}
