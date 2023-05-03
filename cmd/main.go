package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/mserebryaakov/file-service/config"
	_ "github.com/mserebryaakov/file-service/docs"
	"github.com/mserebryaakov/file-service/pkg/httpserver"
	"github.com/mserebryaakov/file-service/pkg/logger"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/mserebryaakov/file-service/internal/file"
	"github.com/mserebryaakov/file-service/internal/file/minio"

	"os"
	"os/signal"
	"syscall"
)

// Swagger документация
// @title           Upload file
// @version         1.0
// @description     API Server for upload file

// @host      localhost:8000
// @BasePath  /
func main() {
	// Получение логера
	log := logger.NewLogger("debug", &logger.MainLogHook{})

	// Получение конфига
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load configs: %v", err)
	}

	// Чтение файла переменных окружения
	env, err := config.GetEnvironment()
	if err != nil {
		log.Fatalf(err.Error())
	}

	// Создание логера для fileService
	serviceLog := logger.NewLogger(env.LogLvl, &file.FileServiceLogHook{})

	// Создание слоёв по чистой архитектуре
	repository, err := minio.NewStorage(cfg.Minio.Host+cfg.Minio.Port, env.MinioRootUser, env.MinioRootPassword, false, serviceLog)
	if err != nil {
		log.Fatalf("minio start error: %v", err)
	}

	service := file.NewService(repository, serviceLog)

	handler := file.NewFilesHandler(service, serviceLog)

	router := gin.New()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	handler.Register(router)

	server := new(httpserver.Server)

	go func() {
		if err := server.Run(cfg.Server.Port, router); err != nil {
			log.Fatal("Failed running server %v", err)
		}
	}()

	// Завершение сервера
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	oscall := <-interrupt
	log.Infof("Shutdown server, %s", oscall)

	if err := server.Shutdown(context.Background()); err != nil {
		log.Errorf("Error occured on server shutting down: %v", err)
	}
}
