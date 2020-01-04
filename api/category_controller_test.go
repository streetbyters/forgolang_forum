package api

import "testing"

type CategoryControllerTest struct {
	*Suite
}

func (s CategoryControllerTest) SetupSuite() {
	SetupSuite(s.Suite)
	UserAuth(s.Suite)
}

func (s CategoryControllerTest) TearDownSuite() {
	TearDownSuite(s.Suite)
}

func Test_CategoryController(t *testing.T) {
	s := CategoryControllerTest{NewSuite()}
	Run(t, s)
}
