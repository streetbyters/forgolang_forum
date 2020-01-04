package api

import (
	"github.com/valyala/fasthttp"
	"path/filepath"
	"testing"
)

type UploadControllerTest struct {
	*Suite
}

func (s UploadControllerTest) SetupSuite() {
	SetupSuite(s.Suite)
	UserAuth(s.Suite)
}

func (s UploadControllerTest) Test_PostUploadFile() {
	file1 := filepath.Join(s.API.App.Config.Path, "assets", "user.png")

	body := make(map[string]interface{})
	body["file"] = file1
	body["dir"] = filepath.Join("test", "upload")

	response := s.File(Post, "/api/v1/upload", body, "file")

	s.Equal(response.Status, fasthttp.StatusCreated)
	data, _ := response.Success.Data.(map[string]interface{})
	s.Equal(data["filename"], "user.png")

	defaultLogger.LogInfo("Post upload file")
}

func (s UploadControllerTest) Test_Should_422Err_PostUploadFileIFDirIsNil() {
	file1 := filepath.Join(s.API.App.Config.Path, "assets", "user.png")

	body := make(map[string]interface{})
	body["file"] = file1

	response := s.File(Post, "/api/v1/upload", body, "file")

	s.Equal(response.Status, fasthttp.StatusUnprocessableEntity)

	defaultLogger.LogInfo("Should be 422 error post upload file if dir is nil")
}

func (s UploadControllerTest) Test_Should_400Err_PostUploadFileIfFileNotValid() {
	file1 := filepath.Join(s.API.App.Config.Path, "files", "tests", "notfound.png")

	body := make(map[string]interface{})
	body["file"] = file1

	response := s.File(Post, "/api/v1/upload", body, "file")

	s.Equal(response.Status, fasthttp.StatusBadRequest)

	defaultLogger.LogInfo("Should 400 error post upload file if file is not valid")
}

func (s UploadControllerTest) TearDownSuite() {
	TearDownSuite(s.Suite)
}

func Test_UploadController(t *testing.T) {
	s := UploadControllerTest{NewSuite()}
	Run(t, s)
}
