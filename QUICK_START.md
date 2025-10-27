# Быстрый старт GophKeeper

## Шаг 1: Установите Docker Desktop

1. Скачайте: https://www.docker.com/products/docker-desktop/
2. Установите Docker Desktop
3. Запустите Docker Desktop
4. Дождитесь полной загрузки (зеленая лампочка)

## Шаг 2: Запустите проект

```bash
# 1. Настройте окружение (запустит PostgreSQL)
make dev-setup

# 2. Запустите сервер (в одном терминале)
make run-server
```

Вы увидите:
```
2024/01/XX XX:XX:XX Starting server on :8080
```

## Шаг 3: Протестируйте API

В **другом терминале**:

```bash
# Запустите тесты
make quick-test
```

Или тестируйте вручную:

```bash
# 1. Регистрация
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}'

# 2. Скопируйте токен из ответа и используйте его:
TOKEN="ваш_токен_здесь"

# 3. Создайте запись
curl -X POST http://localhost:8080/api/entries \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "login",
    "title": "GitHub",
    "data": "encrypted_data",
    "metadata": "{\"site\":\"github.com\"}"
  }'

# 4. Получите список записей
curl -X GET http://localhost:8080/api/entries \
  -H "Authorization: Bearer $TOKEN"
```

## Полезные команды

```bash
# Проверить статус PostgreSQL
docker-compose ps

# Посмотреть логи PostgreSQL
docker-compose logs postgres

# Остановить все контейнеры
docker-compose down

# Очистить базу данных и перезапустить
docker-compose down -v
docker-compose up -d
```

## Что дальше?

- Откройте http://localhost:8080 в браузере (если есть health endpoint)
- Используйте Swagger UI для интерактивной документации
- Изучите полную документацию в `docs/`
