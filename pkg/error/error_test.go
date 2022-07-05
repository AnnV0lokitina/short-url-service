package error

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestError(t *testing.T) {
	err := NewLabelError(TypeConflict, errors.New("test"))
	assert.Equal(t, "["+TypeConflict+"] test", err.Error())
}
