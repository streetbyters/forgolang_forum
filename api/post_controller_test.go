package api

import (
	"forgolang_forum/cmn"
	"forgolang_forum/database/model"
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

func (s PostControllerTest) TearDownSuite() {
	TearDownSuite(s.Suite)
}

func Test_PostController(t *testing.T) {
	s := PostControllerTest{NewSuite()}
	Run(t, s)
}