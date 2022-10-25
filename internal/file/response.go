package file

import "github.com/gin-gonic/gin"

type messageResponse struct {
	Message string `json:"message"`
}

func NewErrorResponse(c *gin.Context, statusCode int, message string) {
	c.AbortWithStatusJSON(statusCode, messageResponse{message})
}

func NewResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, messageResponse{message})
}
