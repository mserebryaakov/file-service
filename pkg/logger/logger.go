package logger

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var log *logrus.Entry

type Logger struct {
	*logrus.Entry
}

func GetLogger() *Logger {
	return &Logger{log}
}

// Инициализация logrus
func init() {
	// Создание логгера
	l := logrus.New()

	// Установка объекта вывода
	l.SetOutput(os.Stdout)

	// Устанавка формата вывода
	l.SetFormatter(&logrus.TextFormatter{})

	// Загрузка файла переменных окружения
	err := godotenv.Load(".env")
	if err != nil {
		l.Fatalf("Failed to load env file: %s", err.Error())
	}

	// Чтение файла переменных окружения
	env, err := godotenv.Read(".env")
	if err != nil {
		l.Fatalf("Error loading env into map[string]string: %s", err.Error())
	}

	// Инициализация LOG_LEVEL из параметров окружения
	logLevel, err := logrus.ParseLevel(env["LOG_LEVEL"])
	if err != nil {
		logLevel = logrus.InfoLevel
	}

	// Установка уровня логирования LOG_LEVEL
	l.SetLevel(logLevel)

	// Возвращение логгера глобальной переменной
	log = logrus.NewEntry(l)
}
