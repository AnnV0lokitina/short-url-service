package userid

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateUserID(t *testing.T) {
	userID, err := GenerateUserID()
	assert.Nil(t, err)
	var i uint32
	assert.IsType(t, i, userID)
}
