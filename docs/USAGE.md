# Руководство пользователя GophKeeper

## Начало работы

### Установка

1. Клонируйте репозиторий
2. Установите зависимости: `go mod download`
3. Соберите проект: `make build`

### Запуск сервера

1. Запустите PostgreSQL (используя Docker):
   ```bash
   make docker-up
   ```

2. Запустите сервер:
   ```bash
   make run-server
   ```

Сервер запустится на `http://localhost:8080`

## Команды клиента

### Информация о версии

Показать версию и дату сборки:
```bash
./bin/gophkeeper-client version
```

### Регистрация пользователя

Зарегистрировать нового пользователя:
```bash
./bin/gophkeeper-client register
```

Вам будет предложено ввести:
- Имя пользователя
- Пароль

### Вход в систему

Войти в свой аккаунт:
```bash
./bin/gophkeeper-client login
```

Вам будет предложено ввести:
- Имя пользователя
- Пароль

### Добавление данных

#### Добавить запись логина
```bash
./bin/gophkeeper-client add login
```

Вам будет предложено ввести:
- Название сайта
- Логин/Email
- Пароль
- Необязательные заметки

#### Добавить текстовую запись
```bash
./bin/gophkeeper-client add text
```

Вам будет предложено ввести:
- Заголовок
- Содержимое
- Необязательные заметки

#### Добавить запись карты
```bash
./bin/gophkeeper-client add card
```

Вам будет предложено ввести:
- Номер карты
- Владелец карты
- Месяц истечения
- Год истечения
- CVV
- Необязательные заметки

### Список записей

Показать все ваши записи:
```bash
./bin/gophkeeper-client list
```

### Получение записи

Получить конкретную запись по ID:
```bash
./bin/gophkeeper-client get <entry-id>
```

### Синхронизация данных

Синхронизировать локальные данные с сервером:
```bash
./bin/gophkeeper-client sync
```

## Примеры использования API

### Использование cURL

#### Регистрация пользователя
```bash
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}'
```

#### Вход
```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}'
```

#### Создание записи (с токеном)
```bash
curl -X POST http://localhost:8080/api/entries \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "type": "login",
    "title": "Example Site",
    "data": "encrypted_data",
    "metadata": "{\"site\":\"example.com\"}"
  }'
```

#### Список записей
```bash
curl -X GET http://localhost:8080/api/entries \
  -H "Authorization: Bearer YOUR_TOKEN"
```

#### Синхронизация данных
```bash
curl -X GET "http://localhost:8080/api/sync?since=1234567890" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Переменные окружения

### Конфигурация сервера

Создайте файл `.env` в корне проекта:

```env
# База данных
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=gophkeeper

# JWT
JWT_SECRET=your-secret-key-here

# Сервер
SERVER_ADDR=:8080
```

## Решение проблем

### Проблемы с подключением к базе данных

Если вы видите ошибки подключения к базе данных:

1. Убедитесь, что PostgreSQL запущен:
   ```bash
   docker-compose ps
   ```

2. Проверьте учётные данные базы данных в файле `.env`

3. Проверьте доступность PostgreSQL:
   ```bash
   psql -h localhost -U postgres -d gophkeeper
   ```

### Проблемы с аутентификацией

Если у вас проблемы со входом:

1. Убедитесь, что используете правильное имя пользователя/пароль
2. Проверьте, что сервер запущен
3. Проверьте, что JWT_SECRET установлен в `.env`

### Проблемы со сборкой

Если сборка не удалась:

1. Обновите Go до последней версии: `go version`
2. Установите зависимости: `go mod tidy`
3. Очистите и пересоберите: `make clean && make build`
