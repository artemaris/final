# Инструкции по сборке

## Требования

- Go 1.21 или выше
- PostgreSQL 13+ (или используйте Docker)
- Make (необязательно, но рекомендуется)

## Быстрый старт

### Использование Make

```bash
# Собрать всё
make build

# Запустить тесты
make test

# Запустить PostgreSQL
make docker-up

# Запустить сервер
make run-server
```

### Ручная сборка

#### Сборка сервера

```bash
cd cmd/server
go build -ldflags "-X main.version=dev -X main.buildDate=$(date -u '+%Y-%m-%d_%H:%M:%S')" -o gophkeeper-server .
```

#### Сборка клиента

```bash
cd cmd/client
go build -ldflags "-X main.version=dev -X main.buildDate=$(date -u '+%Y-%m-%d_%H:%M:%S')" -o gophkeeper-client .
```

## Кроссплатформенная сборка

### Windows

```bash
# Сервер
GOOS=windows GOARCH=amd64 go build -ldflags "-X main.version=dev -X main.buildDate=$(date -u '+%Y-%m-%d_%H:%M:%S')" -o gophkeeper-server.exe ./cmd/server

# Клиент
GOOS=windows GOARCH=amd64 go build -ldflags "-X main.version=dev -X main.buildDate=$(date -u '+%Y-%m-%d_%H:%M:%S')" -o gophkeeper-client.exe ./cmd/client
```

### Linux

```bash
# Сервер
GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=dev -X main.buildDate=$(date -u '+%Y-%m-%d_%H:%M:%S')" -o gophkeeper-server ./cmd/server

# Клиент
GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=dev -X main.buildDate=$(date -u '+%Y-%m-%d_%H:%M:%S')" -o gophkeeper-client ./cmd/client
```

### macOS

```bash
# Сервер
GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.version=dev -X main.buildDate=$(date -u '+%Y-%m-%d_%H:%M:%S')" -o gophkeeper-server ./cmd/server

# Клиент
GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.version=dev -X main.buildDate=$(date -u '+%Y-%m-%d_%H:%M:%S')" -o gophkeeper-client ./cmd/client

# Для Apple Silicon (ARM64)
GOOS=darwin GOARCH=arm64 go build -ldflags "-X main.version=dev -X main.buildDate=$(date -u '+%Y-%m-%d_%H:%M:%S')" -o gophkeeper-client-arm64 ./cmd/client
```

## Тестирование

### Запустить все тесты

```bash
go test ./...
```

### Запустить с покрытием

```bash
go test -cover ./...
```

### Сгенерировать отчёт о покрытии

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

## Зависимости

### Установить зависимости

```bash
go mod download
go mod tidy
```

### Обновить зависимости

```bash
go get -u ./...
go mod tidy
```

## Docker

### Сборка Docker образов

```bash
# Собрать образ сервера
docker build -t gophkeeper-server:latest .

# Запустить PostgreSQL
docker-compose up -d postgres
```

### Docker Compose

```bash
# Запустить все сервисы
docker-compose up -d

# Остановить все сервисы
docker-compose down

# Просмотреть логи
docker-compose logs -f
```

## Настройка окружения

Создайте файл `.env`:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=gophkeeper
JWT_SECRET=your-secret-key-here-change-this
```

## Решение проблем

### Ошибки сборки

1. Убедитесь, что Go установлен: `go version`
2. Очистите кеш модулей: `go clean -modcache`
3. Переустановите зависимости: `go mod download`

### Ошибки тестов

1. Проверьте, что PostgreSQL запущен
2. Проверьте учётные данные базы данных
3. Проверьте переменные окружения

### Слишком большой размер бинарного файла

Используйте флаги `-ldflags="-s -w"` для уменьшения размера:

```bash
go build -ldflags "-s -w -X main.version=dev -X main.buildDate=$(date -u '+%Y-%m-%d_%H:%M:%S')" -o gophkeeper-server ./cmd/server
```
