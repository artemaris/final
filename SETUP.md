# Установка и настройка GophKeeper

## Установка Docker Desktop

### macOS (Apple Silicon)

1. **Скачайте Docker Desktop:**
   ```bash
   brew install --cask docker
   ```

2. **Запустите Docker Desktop** из Applications
3. **Проверьте установку:**
   ```bash
   docker --version
   docker-compose --version
   ```

## Запуск проекта

### Вариант 1: С Docker (Рекомендуется)

```bash
# 1. Запустите PostgreSQL через Docker Compose
docker-compose up -d

# 2. Проверьте, что PostgreSQL запущен
docker-compose ps

# 3. Соберите и запустите сервер
make build-server
./bin/gophkeeper-server

# Или сразу:
make run-server
```

### Вариант 2: Без Docker (требует установки PostgreSQL)

Если у вас установлен PostgreSQL локально:

```bash
# Создайте базу данных
createdb gophkeeper

# Установите переменные окружения
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=your_password
export DB_NAME=gophkeeper

# Запустите сервер
make run-server
```

## Тестирование API

### Автоматическое тестирование

```bash
# Запустите сервер в одном терминале
make run-server

# В другом терминале запустите тесты
./test_api.sh
```

### Ручное тестирование

#### 1. Регистрация пользователя
```bash
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}'
```

#### 2. Вход
```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}'
```

Скопируйте токен из ответа и используйте его в следующих запросах.

#### 3. Создание записи
```bash
curl -X POST http://localhost:8080/api/entries \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "login",
    "title": "GitHub",
    "data": "encrypted_data",
    "metadata": "{\"site\":\"github.com\"}"
  }'
```

#### 4. Список записей
```bash
curl -X GET http://localhost:8080/api/entries \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Проблемы и решения

### Docker не запускается

**macOS:** 
- Проверьте, что Docker Desktop запущен
- Запустите его из Applications
- Дождитесь полной загрузки (иконка в трее)

**Linux:**
```bash
sudo systemctl start docker
```

### Не могу подключиться к PostgreSQL

```bash
# Проверьте статус контейнера
docker-compose ps

# Посмотрите логи
docker-compose logs postgres

# Перезапустите
docker-compose restart postgres
```

### Порт 8080 занят

Измените порт в `cmd/server/main.go` или используйте:
```bash
./bin/gophkeeper-server -addr :8081
```

### Ошибки миграции

```bash
# Очистите базу данных
docker-compose down -v
docker-compose up -d
```

## Структура проекта

```
final/
├── cmd/
│   ├── server/        # Серверное приложение
│   └── client/        # Клиентское приложение
├── internal/
│   ├── crypto/        # Криптография
│   ├── models/        # Модели данных
│   ├── server/        # Серверная логика
│   └── storage/       # Хранилище данных
├── docs/              # Документация
├── Makefile           # Команды сборки
├── docker-compose.yml # Docker конфигурация
└── test_api.sh        # Скрипт тестирования API
```

## Дальнейшие шаги

1. Установите Docker Desktop
2. Запустите `docker-compose up -d`
3. Запустите сервер: `make run-server`
4. Протестируйте API: `./test_api.sh`
5. Откройте Swagger UI для просмотра документации
