package handler

import (
	"github.com/gin-gonic/gin"
	_ "github.com/mserebryaakov/file-service/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Роутинг
func InitRoutes() *gin.Engine {
	// Инициализация роутера
	router := gin.New()

	// Ограничение для multipart memory 3GB
	router.MaxMultipartMemory = 3000 << 20

	// Инициализация swagger docs
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	v1_API := router.Group("/v1")
	{
		// Получение файла
		v1_API.POST("/file", uploadFileHandler)
	}

	return router
}
