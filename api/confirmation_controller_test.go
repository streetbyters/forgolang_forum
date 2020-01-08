package api

import (
	"fmt"
	"forgolang_forum/database"
	"forgolang_forum/database/model"
	"github.com/valyala/fasthttp"
	"testing"
	"time"
)

type ConfirmationControllerTest struct {
	*Suite
}

func (s ConfirmationControllerTest) SetupSuite() {
	SetupSuite(s.Suite)
}

func (s ConfirmationControllerTest) Test_PostConfirmationWithGivenIdentifiers() {
	pwd := "12345"
	user := model.NewUser(&pwd)
	user.Email = "akdilsiz@tecpor.com"
	user.Username = "akdilsiz-confirmation"
	user.IsActive = true
	err := s.API.App.Database.Insert(new(model.User), user, "id")
	s.Nil(err)

	otc := model.NewUserOneTimeCode(user.ID)
	otc.Type = database.Confirmation
	err = s.API.App.Database.Insert(new(model.UserOneTimeCode), otc, "id", "inserted_at")
	s.Nil(err)

	resp := s.JSON(Post,
		fmt.Sprintf("/api/v1/auth/confirmation/%d/%s", user.ID, otc.Code),
		nil)

	s.Equal(resp.Status, fasthttp.StatusCreated)

	defaultLogger.LogInfo("Successfully post confirmation with given identifiers")
}

func (s ConfirmationControllerTest) Test_Should_404Err_PostConfirmationWithGivenIdentifiersIfOTCTimeExceeded() {
	pwd := "12345"
	user := model.NewUser(&pwd)
	user.Email = "email@mail.com"
	user.Username = "akdilsiz-confirmation-2"
	user.IsActive = true
	err := s.API.App.Database.Insert(new(model.User), user, "id")
	s.Nil(err)

	otc := model.NewUserOneTimeCode(user.ID)
	otc.Type = database.Confirmation
	otc.InsertedAt = time.Now().UTC().Add(-time.Minute * 20)
	err = s.API.App.Database.Insert(new(model.UserOneTimeCode), otc, "id", "inserted_at")
	s.Nil(err)

	resp := s.JSON(Post,
		fmt.Sprintf("/api/v1/auth/confirmation/%d/%s", user.ID, otc.Code),
		nil)

	s.Equal(resp.Status, fasthttp.StatusNotFound)

	defaultLogger.LogInfo("Should be 404 error post confirmation with given " +
		"identifiers if otc time exceeded")
}

func (s ConfirmationControllerTest) Test_Should_404Err_PostConfirmationWithGivenIdentifiersIfUserAlreadyActivated() {
	pwd := "12345"
	user := model.NewUser(&pwd)
	user.Email = "email-2@mail.com"
	user.Username = "akdilsiz-confirmation-3"
	user.IsActive = true
	err := s.API.App.Database.Insert(new(model.User), user, "id")
	s.Nil(err)

	otc := model.NewUserOneTimeCode(user.ID)
	otc.Type = database.Confirmation
	err = s.API.App.Database.Insert(new(model.UserOneTimeCode), otc, "id", "inserted_at")
	s.Nil(err)

	userState := model.NewUserState(user.ID)
	userState.State = database.Active
	err = s.API.App.Database.Insert(new(model.UserState), userState, "id")
	s.Nil(err)

	resp := s.JSON(Post,
		fmt.Sprintf("/api/v1/auth/confirmation/%d/%s", user.ID, otc.Code),
		nil)

	s.Equal(resp.Status, fasthttp.StatusNotFound)

	defaultLogger.LogInfo("Should be 404 error post confirmation with given " +
		"identifiers if user already activated")
}

func (s ConfirmationControllerTest) TearDownSuite() {
	TearDownSuite(s.Suite)
}

func Test_ConfirmationController(t *testing.T) {
	s := ConfirmationControllerTest{NewSuite()}
	Run(t, s)
}
