package httputil

import (
	"errors"
	"fmt"
	"github.com/alexandria-oss/core/exception"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestErrorToCode(t *testing.T) {
	err := exception.EntityNotFound
	assert.Equal(t, 404, ErrorToCode(err))

	err = exception.NewErrorDescription(exception.RequiredField,
		fmt.Sprintf(exception.RequiredFieldString, "test"))
	assert.Equal(t, 400, ErrorToCode(err))

	err = errors.New("custom error")
	assert.Equal(t, 500, ErrorToCode(err))
}
