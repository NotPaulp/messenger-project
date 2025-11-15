Messenger Project
Проект мессенджера с API-шлюзом на Go, поддерживающий сообщения, посты и комментарии.

Основные возможности
Аутентификация - регистрация, вход, JWT токены

Сообщения - отправка и получение личных сообщений

Посты - создание и управление постами

Комментарии - комментирование постов

Безопасность - хеширование паролей, черный список JWT

Технологии
Backend: Go 1.25, Gorilla Mux

Базы данных:

PostgreSQL - пользователи

MongoDB - сообщения, посты, комментарии

Redis - черный список JWT токенов

Контейнеризация: Docker, Docker Compose

Аутентификация: JWT токены

Структура проекта
text
messenger-project/
├── cmd/api-gateway/          # Точка входа приложения
├── internal/
│   ├── auth/                 # Аутентификация (JWT, пароли)
│   ├── database/            # Подключение к БД (PostgreSQL, MongoDB)
│   ├── handlers/            # HTTP обработчики
│   │   ├── auth/           # Регистрация, вход, выход
│   │   ├── messages/       # Сообщения
│   │   ├── posts/          # Посты
│   │   └── comments/       # Комментарии
│   ├── models/             # Структуры данных
│   ├── repository/         # Слой работы с данными
│   └── redis/              # Redis клиент
├── pkg/
│   ├── config/             # Конфигурация
│   ├── logger/             # Логирование
│   └── utils/              # Вспомогательные функции
└── api/                    # OpenAPI спецификация

Быстрый старт
Требования
Docker и Docker Compose

Go 1.25+ (для локальной разработки)

Запуск
bash
# Клонирование репозитория
git clone <repository-url>
cd messenger-project

# Запуск всех сервисов
docker-compose up -d

# Приложение доступно на http://localhost:8080
Переменные окружения
Создайте файл .env:

env
SERVER_PORT=8080
JWT_SECRET=your-secret-key
DEBUG=true

# PostgreSQL
DATABASE_URL=postgres://user:password@postgres:5432/db?sslmode=disable
POSTGRES_DB=db
POSTGRES_USER=user
POSTGRES_PASSWORD=password

# MongoDB
MONGO_URL=mongodb://mongo:27017
MONGO_DB=messenger

# Redis
REDIS_URL=redis://redis:6379
REDIS_PASSWORD=password

API Endpoints
Аутентификация
POST /api/register - Регистрация

POST /api/login - Вход

POST /api/logout - Выход

Сообщения
POST /api/messages - Отправить сообщение

GET /api/messages?sender=username - Получить сообщения

DELETE /api/messages - Удалить сообщение

Посты
POST /api/posts - Создать пост

GET /api/posts - Получить посты

DELETE /api/posts - Удалить пост

Комментарии
POST /api/posts/comments - Добавить комментарий

GET /api/posts/comments?post_id=123 - Получить комментарии

DELETE /api/posts/comments - Удалить комментарий

Системные
GET /health - Проверка здоровья сервиса

Безопасность
Пароли хешируются с помощью bcrypt

JWT токены с сроком действия 24 часа

Черный список токенов при выходе

Валидация email и сложности пароля

SQL/Mongo injection protection

Базы данных
PostgreSQL
Пользователи - учетные записи и аутентификация

MongoDB
Сообщения - личные сообщения между пользователями

Посты - публикации пользователей

Комментарии - вложенные в посты

Redis
Черный список - отозванные JWT токены

Разработка
Локальный запуск
bash
# Установка зависимостей
go mod download

# Запуск БД
docker-compose up postgres mongo redis -d

# Запуск приложения
go run ./cmd/api-gateway
Сборка
bash
# Сборка Docker образа
docker build -t messenger-api .

# Или сборка Go приложения
go build -o main ./cmd/api-gateway

Мониторинг
Приложение логирует:

Информационные сообщения (запуск, запросы)

Ошибки (базы данных, валидация)

Debug сообщения (при включенном DEBUG режиме)

Миграции
Таблицы создаются автоматически при запуске:

Таблица пользователей в PostgreSQL

Коллекции в MongoDB

Тестирование
bash
# Пример тестового запроса
curl -X GET http://localhost:8080/health

Примечания
ID сообщений, постов и комментариев - Unix timestamp

Все защищенные endpoints требуют JWT токен в заголовке Authorization

Пользователь может удалять только свои сообщения, посты и комментарии
