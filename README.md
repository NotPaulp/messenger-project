# Messenger Project

Проект мессенджера с API-шлюзом на Go, поддерживающий сообщения, посты и комментарии с реальным временем через WebSocket.

## Основные возможности

- **Аутентификация** - регистрация, вход, JWT токены
- **Сообщения в реальном времени** - отправка и получение личных сообщений через WebSocket
- **Посты** - создание и управление постами
- **Комментарии** - комментирование постов
- **Безопасность** - хеширование паролей, черный список JWT
- **Асинхронная обработка** - Kafka для обработки сообщений

## Технологии

- **Backend**: Go 1.25, Gorilla Mux
- **Базы данных**: 
  - PostgreSQL - пользователи
  - MongoDB - сообщения, посты, комментарии
  - Redis - черный список JWT токенов
- **Контейнеризация**: Docker, Docker Compose
- **Аутентификация**: JWT токены
- **Мессенджинг**: Apache Kafka
- **Реальное время**: WebSocket

## Архитектура
```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Клиент    │◄──►│ API Gateway │◄──►│  PostgreSQL │
│  (WebSocket)│    │   (Go)      │    │   (Users)   │
└─────────────┘    └─────────────┘    └─────────────┘
                           │                 │
                           ▼                 │
                    ┌─────────────┐          │
                    │   Kafka     │          │
                    │  (Messages) │          │
                    └─────────────┘          │
                           │                 │
                           ▼                 │
                    ┌─────────────┐          │
                    │  Consumer   │──────────┘
                    │    (Go)     │
                    └─────────────┘
                           │
                           ▼
                    ┌─────────────┐
                    │   MongoDB   │
                    │ (Messages,  │
                    │ Posts,      │
                    │ Comments)   │
                    └─────────────┘
```
## Структура проекта

```
messenger-project/
├── cmd/api-gateway/          # Точка входа приложения
├── internal/
│   ├── auth/                 # Аутентификация (JWT, пароли)
│   ├── database/            # Подключение к БД (PostgreSQL, MongoDB)
│   ├── handlers/            # HTTP обработчики
│   │   ├── auth/           # Регистрация, вход, выход
│   │   ├── messages/       # Сообщения (REST + WebSocket)
│   │   ├── posts/          # Посты
│   │   └── comments/       # Комментарии
│   ├── models/             # Структуры данных
│   ├── repository/         # Слой работы с данными
│   ├── kafka/              # Kafka producer/consumer
│   ├── websocket/          # WebSocket hub и клиенты
│   └── redis/              # Redis клиент
└── pkg/
    ├── config/             # Конфигурация
    ├── logger/             # Логирование
    └── utils/              # Вспомогательные функции
```

## Быстрый старт

### Требования
- Docker и Docker Compose
- Go 1.25+ (для локальной разработки)

### Запуск
```bash
# Клонирование репозитория
git clone <repository-url>
cd messenger-project

# Запуск всех сервисов
docker-compose up -d

# Приложение доступно на http://localhost:8080
```

### Переменные окружения
Создайте файл `.env`:
```env
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

# Kafka
KAFKA_BROKERS=kafka:9092
KAFKA_TOPIC_MESSAGES=messages
```

## API Endpoints

### Аутентификация
- `POST /api/register` - Регистрация
- `POST /api/login` - Вход
- `POST /api/logout` - Выход

### Сообщения
- `POST /api/messages` - Отправить сообщение (через Kafka)
- `GET /api/messages?sender=username` - Получить сообщения
- `DELETE /api/messages` - Удалить сообщение
- `GET /websocket` - WebSocket соединение для реального времени

### Посты
- `POST /api/posts` - Создать пост
- `GET /api/posts` - Получить посты
- `DELETE /api/posts` - Удалить пост

### Комментарии
- `POST /api/posts/comments` - Добавить комментарий
- `GET /api/posts/comments?post_id=123` - Получить комментарии
- `DELETE /api/posts/comments` - Удалить комментарий

### Системные
- `GET /health` - Проверка здоровья сервиса

## Работа с WebSocket

Для получения сообщений в реальном времени подключитесь к WebSocket:

```javascript
// JavaScript пример
const ws = new WebSocket('ws://localhost:8080/websocket?token=YOUR_JWT_TOKEN');

ws.onmessage = function(event) {
    const message = JSON.parse(event.data);
    console.log('Новое сообщение:', message);
};
```

Или используйте wscat для тестирования:
```bash
wscat -c "ws://localhost:8080/websocket?token=YOUR_JWT_TOKEN"
```

## Безопасность

- Пароли хешируются с помощью bcrypt
- JWT токены с сроком действия 24 часа
- Черный список токенов при выходе через Redis
- Валидация email и сложности пароля
- SQL/Mongo injection protection
- WebSocket аутентификация через JWT

## Поток данных

### Отправка сообщения:
1. Клиент → REST API → Kafka Producer
2. Kafka → Consumer → MongoDB
3. Consumer → WebSocket Hub → Получатель

### Получение сообщений:
1. WebSocket подключение с JWT токеном
2. Сообщения доставляются в реальном времени
3. История доступна через REST API

## Базы данных

### PostgreSQL
- **Пользователи** - учетные записи и аутентификация

### MongoDB
- **Сообщения** - личные сообщения между пользователями
- **Посты** - публикации пользователей
- **Комментарии** - вложенные в посты

### Redis
- **Черный список** - отозванные JWT токены

### Kafka
- **Очередь сообщений** - асинхронная обработка входящих сообщений

## Разработка

### Локальный запуск
```bash
# Установка зависимостей
go mod download

# Запуск БД и Kafka
docker-compose up postgres mongo redis zookeeper kafka -d

# Запуск приложения
go run ./cmd/api-gateway
```

### Сборка
```bash
# Сборка Docker образа
docker build -t messenger-api .

# Или сборка Go приложения
go build -o main ./cmd/api-gateway
```

## Мониторинг

Приложение логирует:
- Информационные сообщения (запуск, запросы)
- Ошибки (базы данных, валидация)
- Debug сообщения (при включенном DEBUG режиме)
- Kafka события (отправка/получение сообщений)
- WebSocket подключения

## Миграции

Таблицы создаются автоматически при запуске:
- Таблица пользователей в PostgreSQL
- Коллекции в MongoDB
- Топик сообщений в Kafka

## Тестирование

```bash
# Пример тестового запроса
curl -X GET http://localhost:8080/health

# Отправка тестового сообщения
curl -X POST http://localhost:8080/api/messages \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"receiver_username":"username","body":"Hello!"}'

# Подключение к WebSocket
wscat -c "ws://localhost:8080/websocket?token=YOUR_JWT_TOKEN"
```

## Примечания

- ID сообщений, постов и комментариев - Unix timestamp
- Все защищенные endpoints требуют JWT токен в заголовке Authorization
- WebSocket использует тот же JWT токен для аутентификации
- Сообщения обрабатываются асинхронно через Kafka
- Пользователь может удалять только свои сообщения, посты и комментарии
- Для реального времени используется WebSocket с автоматической доставкой сообщений
```

