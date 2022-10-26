## file-service - golang сервис для работы с файлами

### Краткое описание:

- Gin
- Minio
- Docker (docker-compose)
- Swagger 3.0

### Применены следующие пакеты:

- gin - http framework
- viper - конфигурация
- logrus - логирование
- godotenv - переменные окружения
- swaggo/swag - документация API
- minio/minio-go/v7 - работа с minio

### Локальное развёртывание:

- Убедиться, что в файле конфига `config/config.json` имя хоста в конфигурации Minio - `"host": "localhost"`
- Убедиться, что Minio запущена
- Создать файл окружения `.env` с тремя переменными:`LOG_LEVEL - уровень логирования (например - LOG_LEVEL=debug), MINIO_ROOT_USER - пользователь minio (например - MINIO_ROOT_USER=minioadmin), MINIO_ROOT_PASSWORD - пароль пользователя minio (например MINIO_ROOT_PASSWORD=minioadmin)`
- Выполнить `go run cmd/main.go` из корня проекта

### Развёртывание с помощью docker-compose:

Для запуска сервиса с помощью команды `docker-compose up` необходимо:
- Убедиться, что в файле конфига `config/config.json` имя хоста в конфигурации Minio совпадает с названием сервиса в файле `docker-compose.yml` (по умолчанию - `"host": "miniodb"`).
- Создать файл окружения `.env` с тремя переменными:`LOG_LEVEL - уровень логирования (например - LOG_LEVEL=debug), MINIO_ROOT_USER - пользователь minio (например - MINIO_ROOT_USER=minioadmin), MINIO_ROOT_PASSWORD - пароль пользователя minio (например MINIO_ROOT_PASSWORD=minioadmin)`
- Убедиться, что 8000, 9000 и 9001 порты не заняты.
Выполнить команду `docker-compose up`.
Далее, обратившись к `localhost:8000/swagger`, можно проверить работоспособность API