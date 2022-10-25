package handler

import "github.com/gin-gonic/gin"

// Форма возврата
type UploadResponse struct {
	Msg string `json:"message"`
}

// Функция обработки ошибок
func newErrorResponse(c *gin.Context, statusCode int, message string) {
	// метод блокирует выполнение следующих обработчик и записывает в ответ статус и сообщение
	c.AbortWithStatusJSON(statusCode, UploadResponse{message})
}
