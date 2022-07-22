package entity

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewBatchURLItem(t *testing.T) {
	item := NewBatchURLItem("correlationID", "originalURL", "localhost:8080")
	assert.IsType(t, &BatchURLItem{}, item)
	assert.IsType(t, &URL{}, item.URL)
	assert.Equal(t, "correlationID", item.CorrelationID)
	assert.Equal(t, "originalURL", item.URL.Original)
}
