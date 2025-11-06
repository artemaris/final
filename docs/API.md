# GophKeeper API Specification

OpenAPI 3.0 Specification для API GophKeeper

## Swagger/OpenAPI Specification

```yaml
openapi: 3.0.0
info:
  title: GophKeeper API
  description: |
    REST API для системы управления паролями GophKeeper.
    
    API обеспечивает:
    - Регистрацию и аутентификацию пользователей
    - CRUD операции с данными (логины, тексты, бинарные данные, карты)
    - Синхронизацию данных между клиентами
    - Безопасное хранение зашифрованных данных
  
    Все защищенные endpoints требуют JWT токен в заголовке Authorization.
  version: 1.0.0
  contact:
    name: GophKeeper Support
    email: support@gophkeeper.com

servers:
  - url: http://localhost:8080
    description: Local development server
  - url: https://api.gophkeeper.com
    description: Production server

tags:
  - name: Authentication
    description: Операции регистрации и аутентификации
  - name: Entries
    description: Управление записями (логины, тексты, бинарные данные, карты)

paths:
  /api/register:
    post:
      tags:
        - Authentication
      summary: Регистрация нового пользователя
      description: Создает нового пользователя и возвращает JWT токен
      operationId: registerUser
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/RegisterRequest'
      responses:
        '201':
          description: Пользователь успешно зарегистрирован
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/AuthResponse'
        '400':
          $ref: '#/components/responses/BadRequest'
        '409':
          description: Пользователь с таким именем уже существует
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
        '500':
          $ref: '#/components/responses/InternalServerError'

  /api/login:
    post:
      tags:
        - Authentication
      summary: Аутентификация пользователя
      description: Проверяет учетные данные и возвращает JWT токен
      operationId: loginUser
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/LoginRequest'
      responses:
        '200':
          description: Успешная аутентификация
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/AuthResponse'
        '400':
          $ref: '#/components/responses/BadRequest'
        '401':
          description: Неверные учетные данные
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'

  /api/entries:
    get:
      tags:
        - Entries
      summary: Список всех записей
      description: Возвращает список всех записей текущего пользователя
      operationId: listEntries
      security:
        - bearerAuth: []
      responses:
        '200':
          description: Список записей
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/Entry'
        '401':
          $ref: '#/components/responses/Unauthorized'
        '500':
          $ref: '#/components/responses/InternalServerError'

    post:
      tags:
        - Entries
      summary: Создать новую запись
      description: Создает новую запись (логин, текст, бинарные данные или карта)
      operationId: createEntry
      security:
        - bearerAuth: []
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/CreateEntryRequest'
      responses:
        '201':
          description: Запись успешно создана
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Entry'
        '400':
          $ref: '#/components/responses/BadRequest'
        '401':
          $ref: '#/components/responses/Unauthorized'
        '500':
          $ref: '#/components/responses/InternalServerError'

  /api/entries/{id}:
    get:
      tags:
        - Entries
      summary: Получить запись по ID
      description: Возвращает конкретную запись пользователя
      operationId: getEntry
      security:
        - bearerAuth: []
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: integer
            format: int64
          description: ID записи
      responses:
        '200':
          description: Данные записи
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Entry'
        '401':
          $ref: '#/components/responses/Unauthorized'
        '404':
          description: Запись не найдена
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'

    put:
      tags:
        - Entries
      summary: Обновить запись
      description: Обновляет существующую запись (оптимистичная блокировка по версии)
      operationId: updateEntry
      security:
        - bearerAuth: []
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: integer
            format: int64
          description: ID записи
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/UpdateEntryRequest'
      responses:
        '200':
          description: Запись успешно обновлена
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Entry'
        '400':
          description: Неверные данные или конфликт версий
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
        '401':
          $ref: '#/components/responses/Unauthorized'
        '404':
          description: Запись не найдена
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'

    delete:
      tags:
        - Entries
      summary: Удалить запись
      description: Удаляет запись пользователя
      operationId: deleteEntry
      security:
        - bearerAuth: []
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: integer
            format: int64
          description: ID записи
      responses:
        '200':
          description: Запись успешно удалена
          content:
            application/json:
              schema:
                type: object
                properties:
                  message:
                    type: string
                    example: Entry deleted successfully
        '401':
          $ref: '#/components/responses/Unauthorized'
        '404':
          description: Запись не найдена
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'

  /api/sync:
    get:
      tags:
        - Entries
      summary: Синхронизация данных
      description: Возвращает записи, обновленные после указанного времени
      operationId: syncEntries
      security:
        - bearerAuth: []
      parameters:
        - name: since
          in: query
          schema:
            type: integer
            format: int64
          description: Unix timestamp начала периода синхронизации
      responses:
        '200':
          description: Обновленные записи
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/SyncResponse'
        '401':
          $ref: '#/components/responses/Unauthorized'
        '500':
          $ref: '#/components/responses/InternalServerError'

components:
  securitySchemes:
    bearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT
      description: JWT токен полученный при регистрации или входе

  schemas:
    RegisterRequest:
      type: object
      required:
        - username
        - password
      properties:
        username:
          type: string
          description: Имя пользователя
          minLength: 3
          maxLength: 50
          example: john_doe
        password:
          type: string
          description: Пароль
          minLength: 8
          format: password
          example: mySecurePass123

    LoginRequest:
      type: object
      required:
        - username
        - password
      properties:
        username:
          type: string
          example: john_doe
        password:
          type: string
          format: password
          example: mySecurePass123

    AuthResponse:
      type: object
      properties:
        token:
          type: string
          description: JWT токен для последующих запросов
          example: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
        user:
          $ref: '#/components/schemas/User'

    User:
      type: object
      properties:
        id:
          type: integer
          format: int64
          description: ID пользователя
        username:
          type: string
          description: Имя пользователя
        created_at:
          type: string
          format: date-time
          description: Дата регистрации

    Entry:
      type: object
      properties:
        id:
          type: integer
          format: int64
          description: ID записи
        user_id:
          type: integer
          format: int64
          description: ID пользователя
        type:
          type: string
          enum: [login, text, binary, card]
          description: Тип записи
        title:
          type: string
          description: Заголовок записи
        data:
          type: string
          description: Зашифрованные данные (base64)
        metadata:
          type: string
          description: JSON метаданные
        version:
          type: integer
          description: Версия для оптимистичной блокировки
        updated_at:
          type: string
          format: date-time
        created_at:
          type: string
          format: date-time

    CreateEntryRequest:
      type: object
      required:
        - type
        - title
        - data
      properties:
        type:
          type: string
          enum: [login, text, binary, card]
          description: Тип записи
        title:
          type: string
          description: Заголовок записи
          maxLength: 255
        data:
          type: string
          description: Зашифрованные данные
        metadata:
          type: string
          description: Дополнительные метаданные в формате JSON

    UpdateEntryRequest:
      type: object
      required:
        - version
      properties:
        title:
          type: string
          maxLength: 255
        data:
          type: string
        metadata:
          type: string
        version:
          type: integer
          description: Текущая версия записи

    SyncResponse:
      type: object
      properties:
        entries:
          type: array
          items:
            $ref: '#/components/schemas/Entry'
        version:
          type: integer
          description: Текущая версия на сервере

    Error:
      type: object
      properties:
        error:
          type: string
          description: Описание ошибки

  responses:
    BadRequest:
      description: Неверный запрос
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/Error'
          example:
            error: Invalid request

    Unauthorized:
      description: Требуется аутентификация
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/Error'
          example:
            error: Missing authorization header

    InternalServerError:
      description: Внутренняя ошибка сервера
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/Error'
          example:
            error: Internal server error
```

## Примеры использования

### Регистрация пользователя

```bash
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "password": "securePass123"
  }'
```

**Ответ:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "john_doe",
    "created_at": "2024-01-15T10:30:00Z"
  }
}
```

### Создание записи с логином

```bash
curl -X POST http://localhost:8080/api/entries \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "type": "login",
    "title": "GitHub",
    "data": "encrypted_data_here",
    "metadata": "{\"site\":\"github.com\"}"
  }'
```

### Синхронизация данных

```bash
curl -X GET "http://localhost:8080/api/sync?since=1640995200" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Типы данных и примеры

### Login Entry
```json
{
  "type": "login",
  "title": "GitHub",
  "data": "encrypted_login_data",
  "metadata": "{\"site\":\"github.com\",\"login\":\"user@example.com\"}"
}
```

### Text Entry
```json
{
  "type": "text",
  "title": "Secret Note",
  "data": "encrypted_text_data",
  "metadata": "{\"notes\":\"Important information\"}"
}
```

### Card Entry
```json
{
  "type": "card",
  "title": "Visa Card",
  "data": "encrypted_card_data",
  "metadata": "{\"number\":\"****\",\"holder\":\"John Doe\",\"expiry\":\"12/25\"}"
}
```

## Коды ошибок

- `400` - Неверный запрос (неверные данные)
- `401` - Не авторизован (отсутствует или неверный токен)
- `404` - Запись не найдена
- `409` - Конфликт (пользователь уже существует)
- `500` - Внутренняя ошибка сервера

## Безопасность

1. **JWT токены** - Срок действия 24 часа
2. **HTTPS** - Рекомендуется использовать в production
3. **Валидация данных** - Все входящие данные проверяются
4. **Изоляция пользователей** - Данные строго изолированы по user_id
5. **Шифрование** - Все чувствительные данные шифруются на клиенте

## Версионирование

API использует оптимистичную блокировку через поле `version`. 
При обновлении записи обязательно указывайте текущую версию.
Если версия не совпадает, будет возвращена ошибка 400.
