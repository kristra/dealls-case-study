package utils

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestWrapSuccessResponse(t *testing.T) {
	mockData := gin.H{
		"id":   1,
		"name": "John Doe",
	}

	response := WrapSuccessResponse(mockData)

	assert.Equal(t, "success", response["message"])
	assert.Equal(t, mockData, response["data"])
}
