package api

import (
	"context"
	"fmt"
	"forgolang_forum/database/model"
	"github.com/valyala/fasthttp"
	"strconv"
	"testing"
	"time"
)

type SearchUserControllerTest struct {
	*Suite
}

func (s SearchUserControllerTest) SetupSuite() {
	SetupSuite(s.Suite)
}

func (s SearchUserControllerTest) Test_SearchUserWithGivenQuery() {
	UserAuth(s.Suite, "user")
	for i := 0; i < 200; i++ {
		user := model.NewUser(nil)
		if i < 50 {
			user.Email = fmt.Sprintf("search-user-%d@mail.com", i)
			user.Username = fmt.Sprintf("search-user-%d", i)

			if i < 20 {
				user.EmailHidden = true
			} else {
				user.EmailHidden = false
			}
		} else {
			user.Email = fmt.Sprintf("other-%d@mail.com", i)
			user.Username = fmt.Sprintf("other-%d", i)
		}
		err := s.API.App.Database.Insert(new(model.User), user, "id")
		s.Nil(err)
		s.API.App.ElasticClient.Index().Index("users").
			Id(strconv.Itoa(i)).
			BodyString(user.ToJSON()).
			Do(context.TODO())
	}

	ch := make(chan bool)

	time.AfterFunc(time.Second * 5, func() {
		ch<-true
	})

	<-ch

	response := s.JSON(Get, "/api/v1/search/user?query=search-user&order_by=asc", nil)

	s.Equal(response.Status, fasthttp.StatusOK)
	s.Equal(response.Success.TotalCount, int64(50))
	s.Equal(len(response.Success.Data.([]interface{})), 40)

	user := response.Success.Data.([]interface{})[0].(map[string]interface{})

	s.Equal(user["username"], "search-user-0")
	s.Equal(user["email"], "***")
	s.Equal(user["email_hidden"], true)

	defaultLogger.LogInfo("Search user with given query")
}

func (s SearchUserControllerTest) Test_SearchUserWithGivenMultiFieldAndEmailParamAndSuperadminRole() {
	UserAuth(s.Suite)
	for i := 0; i < 200; i++ {
		user := model.NewUser(nil)
		if i < 10 {
			user.Email = fmt.Sprintf("special-%d@mail.com", i)
			user.Username = fmt.Sprintf("special-%d", i)
			user.EmailHidden = false
		} else {
			user.Email = fmt.Sprintf("other-new-%d@mail.com", i)
			user.Username = fmt.Sprintf("other-new-%d", i)
		}
		err := s.API.App.Database.Insert(new(model.User), user, "id")
		s.Nil(err)
		s.API.App.ElasticClient.Index().Index("users").
			Id(strconv.Itoa(i)).
			BodyString(user.ToJSON()).
			Do(context.TODO())
	}

	ch := make(chan bool)

	time.AfterFunc(time.Second * 5, func() {
		ch<-true
	})

	<-ch

	response := s.JSON(Get, "/api/v1/search/user?email=special&multi=true&order_by=asc", nil)

	s.Equal(response.Status, fasthttp.StatusOK)
	s.Equal(response.Success.TotalCount, int64(10))
	s.Equal(len(response.Success.Data.([]interface{})), 10)

	user := response.Success.Data.([]interface{})[0].(map[string]interface{})

	s.Equal(user["username"], "special-0")
	s.Equal(user["email"], "special-0@mail.com")
	s.Equal(user["email_hidden"], false)

	defaultLogger.LogInfo("Search user with multi field and email and superadmin role")
}

func (s SearchUserControllerTest) Test_SearchUserWithGivenMultiFieldAndUsernameParamAndSuperadminRole() {
	UserAuth(s.Suite)
	for i := 0; i < 200; i++ {
		user := model.NewUser(nil)
		if i < 10 {
			user.Email = fmt.Sprintf("username-%d@mail.com", i)
			user.Username = fmt.Sprintf("username-%d", i)
			user.EmailHidden = false
		} else {
			user.Email = fmt.Sprintf("other-other-%d@mail.com", i)
			user.Username = fmt.Sprintf("other-other-%d", i)
		}
		err := s.API.App.Database.Insert(new(model.User), user, "id")
		s.Nil(err)
		s.API.App.ElasticClient.Index().Index("users").
			Id(strconv.Itoa(i)).
			BodyString(user.ToJSON()).
			Do(context.TODO())
	}

	ch := make(chan bool)

	time.AfterFunc(time.Second * 5, func() {
		ch<-true
	})

	<-ch

	response := s.JSON(Get, "/api/v1/search/user?username=username&multi=true&order_by=asc", nil)

	s.Equal(response.Status, fasthttp.StatusOK)
	s.Equal(response.Success.TotalCount, int64(10))
	s.Equal(len(response.Success.Data.([]interface{})), 10)

	user := response.Success.Data.([]interface{})[0].(map[string]interface{})

	s.Equal(user["username"], "username-0")
	s.Equal(user["email"], "username-0@mail.com")
	s.Equal(user["email_hidden"], false)

	defaultLogger.LogInfo("Search user with multi field and username and superadmin role")
}

func (s SearchUserControllerTest) Test_SearchUserWithGivenMultiFieldAndIsActiveParamAndSuperadminRole() {
	UserAuth(s.Suite)
	for i := 0; i < 200; i++ {
		user := model.NewUser(nil)
		if i < 150 {
			user.Email = fmt.Sprintf("isactive-%d@mail.com", i)
			user.Username = fmt.Sprintf("isactive-%d", i)
			user.EmailHidden = false
		} else {
			user.Email = fmt.Sprintf("other-isactive-%d@mail.com", i)
			user.Username = fmt.Sprintf("other-isactive-%d", i)
			user.IsActive = true
		}
		err := s.API.App.Database.Insert(new(model.User), user, "id")
		s.Nil(err)
		s.API.App.ElasticClient.Index().Index("users").
			Id(strconv.Itoa(i)).
			BodyString(user.ToJSON()).
			Do(context.TODO())
	}

	ch := make(chan bool)

	time.AfterFunc(time.Second * 5, func() {
		ch<-true
	})

	<-ch

	response := s.JSON(Get, "/api/v1/search/user?is_active=true&multi=true&order_by=asc", nil)

	s.Equal(response.Status, fasthttp.StatusOK)
	s.Equal(response.Success.TotalCount, int64(50))
	s.Equal(len(response.Success.Data.([]interface{})), 40)

	user := response.Success.Data.([]interface{})[0].(map[string]interface{})

	s.Equal(user["is_active"], true)

	defaultLogger.LogInfo("Search user with multi field and is active and superadmin role")
}

func (s SearchUserControllerTest) TearDownSuite() {
	TearDownSuite(s.Suite)
}

func Test_SearchUserController(t *testing.T) {
	s := SearchUserControllerTest{NewSuite()}
	Run(t, s)
}