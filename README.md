# Authentication Service Test Task

**Описание проекта:**  
Микросервис на Go для аутентификации пользователей через JWT-токены (Access + Refresh).  
Реализует безопасную выдачу, обновление и валидацию токенов с привязкой к IP-адресу и защитой от повторного использования.  

---

## Основной функционал

- Генерация пары Access (JWT) + Refresh (Base64 + Bcrypt) токенов.
- Валидация токенов с проверкой связки Access ↔ Refresh.
- Защита от подмены IP (отправка email-уведомления при изменении).
- Поддержка PostgreSQL для хранения данных.
- Конфигурация через `config.yaml`.

---

## Конфигурация

Настройки сервиса в `config.yaml`:

```yaml
server:
  host: 0.0.0.0
  port: 8888
  env: local

auth-secret: supersecretkey  # Секрет для подписи JWT

postgres:
  user: postgres
  password: postgres
  name: mydb
  port: 5432
  host: postgres
  migrate: true
  sslmode: disable

notifyer:
  smtp-host: smtp.gmail.com
  smtp-port: 587
  email: "___email"     # Моковые данные для тестов
  password: "___password"
```

## Запуск

```bash
docker-compose up --build
```

## REST API
- (POST /login)
Пример запроса:
```yaml
{
"id": "0f5ae05b-4d6e-0c0e-43f6-73deb928c0a3",
"email": "user@example.com"
}
```

Пример ответа:
```yaml
{
"Access": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9...",
"Refresh": "pKnh+HoarQhh3ItA6fSAU3KaZUQY2djIktY7Egi9g9Y="
}
```

- (POST /refresh)
Запрос и ответ имеют одинаковую структуру:
```yaml
{
"Access": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9...",
"Refresh": "pKnh+HoarQhh3ItA6fSAU3KaZUQY2djIktY7Egi9g9Y="
}
```

## Технологии

- Go

- JWT

- PostgreSQL

- Docker

## Структура проекта

```yaml
📄 README.md
🐳 docker-compose.yaml
🐳 dockerfile
📦 go.mod
🔐 go.sum
💻 cmd/
└── main.go
⚙️ config.yaml
🛠️ internal/
├── api/ *HTTP-обработчики*
│ ├── handlers.go
│ └── mocks/
├── auth/ Логика аутентификации
│ ├── auth.go
│ └── mocks/
├── config/ Конфигурация
│ └── config.go
├── domain/ Модели данных
│ └── models.go
├── logger/ Логирование
│ └── logger.go
├── notify/ *Email-уведомления*
│ ├── notify.go
│ └── mocks/
├── repository/ Работа с БД
│ ├── postgres/
│ └── mocks/
└── token/ Генерация токенов
├── token.go
└── mocks/
```

## Тестирование

Запуск тестов:
```bash
go test ./...
```