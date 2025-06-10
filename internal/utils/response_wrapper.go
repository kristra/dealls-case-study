package utils

import "github.com/gin-gonic/gin"

func WrapSuccessResponse(data interface{}) gin.H {

	return gin.H{
		"message": "success",
		"data":    data,
	}
}
