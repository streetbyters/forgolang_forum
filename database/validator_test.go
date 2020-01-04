package database

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestStruct struct {
	Name   string `validate:"required"`
	Code   string
	Detail string `validate:"gte=5"`
}

func TestValidateStruct(t *testing.T) {
	testStruct := new(TestStruct)
	testStruct.Code = "code"
	testStruct.Detail = "det"

	errs, err := ValidateStruct(testStruct)
	assert.NotNil(t, err)
	assert.Equal(t, errs["detail"], "gte: 5")
	assert.Equal(t, errs["name"], "required")
}
