# Messenger Project

Микросервисный мессенджер с ML-анализом сообщений, построенный на Go с использованием gRPC для межсервисной коммуникации, Kafka для асинхронной обработки и WebSocket для real-time коммуникации.

## Основные возможности

- **Микросервисная архитектура** - разделение на независимые сервисы (API Gateway, User Service, ML Analyzer)
- **Аутентификация** - регистрация, вход, JWT токены с blacklist в Redis
- **Сообщения в реальном времени** - WebSocket с автоматической retry логикой для доставки
- **ML-анализ сообщений** - автоматическая категоризация и оценка токсичности через OpenRouter API
- **Посты и комментарии** - полнофункциональная система публикаций
- **Безопасность** - bcrypt хеширование, JWT blacklist, валидация данных
- **Отказоустойчивость** - retry механизм для pending сообщений, статусы доставки

## Технологии

- **Backend**: Go 1.25, Gorilla Mux, gRPC
- **Базы данных**: 
  - PostgreSQL - пользователи
  - MongoDB - сообщения, посты, комментарии
  - Redis - blacklist JWT токенов
- **Контейнеризация**: Docker, Docker Compose
- **Аутентификация**: JWT токены
- **Мессенджинг**: Apache Kafka (2 топика)
- **Реальное время**: WebSocket
- **ML**: OpenRouter API (DeepSeek V3.1)

## Архитектура

```
┌─────────────────┐
│     Клиент      │
│   (WebSocket)   │
└────────┬────────┘
         │
         ▼
┌─────────────────────────────────────────┐
│          API Gateway (:8080)            │
│  • REST API                             │
│  • WebSocket Hub                        │
│  • Retry Service (30s)                  │
└──┬────────┬─────────┬──────────────┬───┘
   │        │         │              │
   │ gRPC   │ Kafka   │ MongoDB      │
   │        │         │              │
   ▼        ▼         ▼              │
┌──────┐ ┌─────┐ ┌─────────┐        │
│User  │ │Kafka│ │ MongoDB │        │
│Service│ │     │ │         │        │
│:8082 │ │:9092│ │ :27017  │        │
│gRPC  │ │     │ │         │        │
│:8083 │ │     │ │         │        │
└──┬───┘ └──┬──┘ └─────────┘        │
   │        │                        │
   │ PostgreSQL  messages →          │
   │        │    ml_analyze_messages │
   ▼        ▼                        ▼
┌──────┐ ┌────────────────────────────┐
│Postgres   ML Analyzer (:8081)      │
│:5432 │ │  • Categorization          │
└──────┘ │  • Toxicity Detection      │
         │  • OpenRouter Integration  │
         └────────────────────────────┘
```

## Микросервисы

### 1. API Gateway (:8080)
**Назначение**: Главная точка входа для клиентов

**Функции**:
- REST API для сообщений, постов, комментариев
- WebSocket Hub для real-time коммуникации
- Retry Service для повторной отправки pending сообщений
- gRPC клиент для проверки существования пользователей

**Зависимости**:
- User Service (gRPC)
- MongoDB
- Kafka Producer
- Redis

### 2. User Service (:8082 HTTP, :8083 gRPC)
**Назначение**: Управление пользователями и аутентификация

**Функции**:
- Регистрация и аутентификация пользователей
- JWT токены и blacklist
- gRPC сервер для проверки пользователей
- REST API для auth операций

**Зависимости**:
- PostgreSQL
- Redis

### 3. ML Analyzer (:8081)
**Назначение**: Анализ содержимого сообщений

**Функции**:
- Категоризация сообщений (question, statement, command, greeting, farewell, complaint, compliment, request, general, swear, spam, other)
- Оценка токсичности (0.0-1.0)
- Интеграция с OpenRouter API (DeepSeek V3.1)

**Зависимости**:
- Kafka Consumer (ml_analyze_messages topic)
- MongoDB
- OpenRouter API

## Структура проекта

```
messenger-project/
├── cmd/
│   ├── api-gateway/          # API Gateway сервис
│   │   ├── main.go
│   │   └── Dockerfile
│   ├── user-service/         # User Service
│   │   ├── main.go
│   │   └── Dockerfile
│   └── ml-analyzer/          # ML Analyzer
│       ├── main.go
│       └── Dockerfile
├── internal/
│   ├── auth/                 # JWT, пароли, middleware
│   ├── database/            # PostgreSQL, MongoDB подключения
│   ├── grpc/                # gRPC клиент и сервер
│   ├── handlers/
│   │   ├── auth/           # Регистрация, login, logout
│   │   ├── messages/       # Сообщения (REST + WebSocket)
│   │   ├── posts/          # Посты
│   │   ├── comments/       # Комментарии
│   │   └── openrouter/     # OpenRouter API интеграция
│   ├── kafka/              # Producer и Consumers
│   ├── ml-analyze/         # ML анализ логика
│   ├── models/             # Структуры данных
│   ├── redis/              # Redis клиент
│   ├── repository/         # Слой данных
│   │   ├── api-gateway/    # Репозитории для API Gateway
│   │   └── user-service/   # Репозитории для User Service
│   ├── retry/              # Retry сервис
│   └── websocket/          # WebSocket Hub
├── proto/user/             # Protocol Buffers для gRPC
├── pkg/
│   ├── config/             # Конфигурация
│   ├── logger/             # Логирование
│   └── utils/              # Утилиты
└── docker-compose.yml
```

## Быстрый старт

### Требования
- Docker и Docker Compose
- OpenRouter API key (для ML анализа)
- Go 1.25+ (для локальной разработки)

### Запуск

```bash
# Клонирование репозитория
git clone <repository-url>
cd messenger-project

# Создать .env файл (см. ниже)
# Добавить OPENROUTER_API_KEY

# Запуск всех сервисов
docker-compose up -d

# Проверка здоровья сервисов
curl http://localhost:8080/health  # API Gateway
curl http://localhost:8082/health  # User Service
```

### Переменные окружения

Создайте файл `.env`:

```env
# Порты сервисов
SERVER_PORT=8080                    # API Gateway HTTP
USER_SERVICE_PORT=8082              # User Service HTTP
GRPC_PORT=8083                      # User Service gRPC
ML_ANALYZER_PORT=8081               # ML Analyzer
USER_GRPC_ADDR=user-service:8083    # gRPC адрес

# Безопасность
JWT_SECRET=your-secret-key-change-in-production
DEBUG=true

# PostgreSQL
DATABASE_URL=postgres://user:password@postgres:5432/db?sslmode=disable
POSTGRES_DB=db
POSTGRES_USER=user
POSTGRES_PASSWORD=password

# MongoDB
MONGO_URL=mongodb://mongo:27017
MONGO_DB=messenger
MONGO_MESSAGES_COLLECTION=messages
MONGO_POSTS_COLLECTION=posts

# Redis
REDIS_URL=redis://redis:6379
REDIS_PASSWORD=password

# Kafka
KAFKA_BROKERS=kafka:9092
KAFKA_TOPIC_MESSAGES=messages
KAFKA_TOPIC_ML_ANALYZE_MESSAGES=ml_analyze_messages

# OpenRouter (для ML анализа)
OPENROUTER_API_KEY=your-openrouter-api-key
```

## API Endpoints

### Аутентификация (User Service :8082)
- `POST /api/register` - Регистрация нового пользователя
- `POST /api/login` - Вход (возвращает JWT токен)
- `POST /api/logout` - Выход (blacklist токена)

### Сообщения (API Gateway :8080)
- `POST /api/messages` - Отправить сообщение (→ Kafka → ML анализ)
- `GET /api/messages?sender=username&all=true` - Получить все сообщения
- `GET /api/messages?sender=username` - Получить последнее сообщение
- `DELETE /api/messages` - Удалить сообщение (только свое)
- `GET /websocket?token=JWT_TOKEN` - WebSocket для real-time

### Посты (API Gateway :8080)
- `POST /api/posts` - Создать пост
- `GET /api/posts?author_username=user&all=true` - Получить все посты
- `GET /api/posts?author_username=user` - Получить последний пост
- `DELETE /api/posts` - Удалить пост (только свой)

### Комментарии (API Gateway :8080)
- `POST /api/posts/comments` - Добавить комментарий к посту
- `GET /api/posts/comments?post_id=123&all=true` - Получить все комментарии
- `GET /api/posts/comments?post_id=123` - Получить последний комментарий
- `DELETE /api/posts/comments` - Удалить комментарий (только свой)

### Системные
- `GET /health` - Проверка здоровья (на каждом сервисе)

## Поток данных

### 1. Отправка сообщения с ML-анализом

```
Клиент → POST /api/messages
    ↓
API Gateway: проверка JWT
    ↓
Kafka Producer → messages topic
    ↓
┌─────────────────────────────────┐
│  Kafka Consumer (API Gateway)   │
│  1. Сохранение в MongoDB        │
│  2. Отправка в WebSocket        │
│  3. Публикация в ml_analyze     │
└─────────────────────────────────┘
    ↓
ML Analyzer Consumer
    ↓
OpenRouter API (DeepSeek V3.1)
    ↓
Обновление документа в MongoDB:
  • category
  • toxicity_score
  • toxic (boolean)
  • analyzed_at
```

### 2. Retry механизм

```
Retry Service (каждые 30 секунд)
    ↓
Поиск сообщений со status=0 (pending)
    ↓
Попытка отправить через WebSocket
    ↓
Success: status→1, sent_at обновляется
Failed: retries++
    ↓
retries >= 3 → status→-1 (failed)
```

### 3. WebSocket real-time доставка

```
Клиент → WebSocket connect с JWT
    ↓
Hub регистрирует клиента
    ↓
Новое сообщение → BroadcastToUser()
    ↓
Если пользователь online:
    → сообщение доставлено
    → status→1 (sent)
Если offline:
    → status→0 (pending)
    → Retry Service попытается позже
```

## Работа с WebSocket

### JavaScript пример

```javascript
const token = 'YOUR_JWT_TOKEN';
const ws = new WebSocket(`ws://localhost:8080/websocket?token=${token}`);

ws.onopen = () => {
    console.log('Connected to messenger');
};

ws.onmessage = (event) => {
    const message = JSON.parse(event.data);
    console.log('New message:', message);
    // message содержит:
    // - id, sender_username, receiver_username
    // - body, sent_at, status
    // - category, toxic, toxicity_score (после ML анализа)
};

ws.onerror = (error) => {
    console.error('WebSocket error:', error);
};

ws.onclose = () => {
    console.log('Disconnected');
};
```

### wscat для тестирования

```bash
# Установка wscat
npm install -g wscat

# Подключение
wscat -c "ws://localhost:8080/websocket?token=YOUR_JWT_TOKEN"

# Вы будете получать сообщения в реальном времени
```

## ML-анализ сообщений

### Категории сообщений

ML Analyzer определяет одну из категорий:
- `question` - вопросы
- `statement` - утверждения
- `command` - команды/приказы
- `greeting` - приветствия
- `farewell` - прощания
- `complaint` - жалобы
- `compliment` - комплименты
- `request` - просьбы
- `general` - общие сообщения
- `swear` - ругательства
- `spam` - спам
- `other` - прочее

### Оценка токсичности

- **Шкала**: 0.0 (не токсично) → 1.0 (экстремально токсично)
- **Порог**: сообщение считается токсичным при score > 0.7
- **Использование**: фильтрация, модерация, аналитика

### Пример анализа

```json
{
  "id": 1703001234,
  "sender_username": "user1",
  "receiver_username": "user2",
  "body": "Hello! How are you?",
  "category": "greeting",
  "toxic": false,
  "toxicity_score": 0.05,
  "sent_at": "2024-12-22T10:30:00Z",
  "status": 1,
  "analyzed_at": "2024-12-22T10:30:02Z"
}
```

## Статусы сообщений

| Статус | Значение | Описание |
|--------|----------|----------|
| -1 | Failed | Не удалось доставить после 3 попыток |
| 0 | Pending | Ожидает доставки (пользователь offline) |
| 1 | Sent | Успешно доставлено |

**Retry логика**:
- Максимум 3 попытки (`MaxRetries = 3`)
- Проверка каждые 30 секунд
- После 3 неудачных попыток → статус Failed

## Безопасность

### Аутентификация
- **Пароли**: bcrypt хеширование (cost=10)
- **JWT**: 24-часовой срок действия
- **Blacklist**: отозванные токены в Redis с TTL
- **Валидация**:
  - Email формат (regex)
  - Пароль: минимум 8 символов, буквы + цифры

### gRPC коммуникация
- User Service gRPC сервер: проверка существования пользователей
- API Gateway gRPC клиент: валидация перед созданием сообщений/постов
- Insecure credentials (для dev окружения)

### Защита данных
- SQL injection protection через prepared statements
- MongoDB injection protection через typed queries
- Валидация имен таблиц/коллекций (regex)
- Context timeout для всех DB операций

## Базы данных

### PostgreSQL (Users)
```sql
CREATE TABLE users (
    id TEXT PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL
);
```

### MongoDB Collections

#### messages
```javascript
{
  id: Int64,                    // Unix timestamp
  sender_username: String,
  receiver_username: String,
  body: String,
  category: String,             // ML анализ
  toxic: Boolean,               // ML анализ
  toxicity_score: Float32,      // ML анализ (0.0-1.0)
  sent_at: ISODate,
  status: Int32,                // -1, 0, 1
  status_updated_at: ISODate,
  retries: Int32,               // Количество попыток
  analyzed_at: ISODate          // Время ML анализа
}
```

#### posts
```javascript
{
  id: Int64,
  author_username: String,
  body: String,
  comments: [Comment],
  published_at: ISODate
}
```

#### Comment (embedded)
```javascript
{
  id: Int64,
  author_username: String,
  body: String,
  published_at: ISODate
}
```

### Redis
- **Key**: JWT token string
- **Value**: "blacklisted"
- **TTL**: время до истечения токена

### Kafka Topics

#### messages
- **Назначение**: Новые сообщения от клиентов
- **Consumers**: API Gateway Consumer
- **Partitions**: 1
- **Replication**: 1

#### ml_analyze_messages
- **Назначение**: Сообщения для ML анализа
- **Consumers**: ML Analyzer
- **Partitions**: 1
- **Replication**: 1

## Разработка

### Локальный запуск без Docker

```bash
# Установка зависимостей
go mod download

# Запуск инфраструктуры
docker-compose up postgres mongo redis zookeeper kafka -d

# Запуск сервисов (в разных терминалах)

# Terminal 1: User Service
go run ./cmd/user-service

# Terminal 2: API Gateway
go run ./cmd/api-gateway

# Terminal 3: ML Analyzer
export OPENROUTER_API_KEY=your-key
go run ./cmd/ml-analyzer
```

### Генерация gRPC кода

```bash
# Установка protoc и плагинов
# macOS: brew install protobuf
# Linux: apt install protobuf-compiler

go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Генерация
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/user/user.proto
```

### Сборка образов

```bash
# API Gateway
docker build -f cmd/api-gateway/Dockerfile -t messenger-api-gateway .

# User Service
docker build -f cmd/user-service/Dockerfile -t messenger-user-service .

# ML Analyzer
docker build -f cmd/ml-analyzer/Dockerfile -t messenger-ml-analyzer .
```

## Мониторинг и логирование

### Логи сервисов

**API Gateway**:
- HTTP запросы (метод, путь, время)
- WebSocket подключения/отключения
- Kafka события
- Retry попытки

**User Service**:
- HTTP запросы
- gRPC вызовы
- Аутентификация события

**ML Analyzer**:
- Обработка сообщений
- OpenRouter API вызовы
- Результаты анализа

### Health checks

Каждый сервис предоставляет `/health` endpoint:

```bash
# Проверка всех сервисов
curl http://localhost:8080/health
curl http://localhost:8082/health

# Ответ
{
  "status": "healthy",
  "service": "api-gateway",
  "timestamp": "2024-12-22T10:30:00Z",
  "version": "1.0.0"
}
```

## Тестирование

### Полный цикл тестирования

```bash
# 1. Регистрация пользователя
curl -X POST http://localhost:8082/api/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "Test1234"
  }'

# 2. Вход
curl -X POST http://localhost:8082/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "Test1234"
  }'
# Сохранить token из ответа

# 3. Отправка сообщения
curl -X POST http://localhost:8080/api/messages \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "receiver_username": "otheruser",
    "body": "Hello! This is a test message."
  }'

# 4. Получение сообщений
curl -X GET "http://localhost:8080/api/messages?sender=testuser&all=true" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN"

# 5. WebSocket подключение
wscat -c "ws://localhost:8080/websocket?token=YOUR_TOKEN"
```
