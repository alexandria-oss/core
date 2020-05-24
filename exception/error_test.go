package exception

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetErrorDescription(t *testing.T) {
	err := NewErrorDescription(RequiredField, fmt.Sprintf(RequiredFieldString, "test"))
	assert.Equal(t, "missing required request field test", GetErrorDescription(err))
}
