# get_order

Этот проект представляет собой небольшой микросервис, написанный на языке Go, обрабатывающий присланные запросы по HTTP (REST) и сообщения в топик (Kafka).

## Архитектура

```text
Kafka (producer) → Consumer → PostgreSQL ↔ Cache (LRU)
                                   ↓
                             HTTP Server (REST + UI)
```

## Структура проекта

```text
├───cmd
│   ├───get_order      # основной бинарь сервиса: запуск API-сервера и консьюмера Kafka
│   └───producer       # утилита-продюсер: скрипт для отправки заказов в Kafka
│
├───internal           # внутренняя бизнес-логика и реализация (по Clean Architecture)
│   ├───adapter        # слой адаптеров: внешние интерфейсы, приводящие данные к usecase
│   │   ├───cache
│   │   │   └───inmemory   # реализация in-memory кэша (LRU, read & write aside)
│   │   ├───consumer
│   │   │   └───kafka      # адаптер Kafka-консьюмера
│   │   ├───controller
│   │   │   └───rest       # REST-контроллеры (HTTP endpoints)
│   │   ├───mapper         # маппинг DTO <-> domain модели
│   │   └───middleware     # HTTP-middleware (логирование)
│   │
│   ├───config             # работа с конфигурацией (переменные окружения)
│   ├───domain             # описание доменных сущностей
│   ├───repository
│   │   └───postgres       # реализация репозитория заказов для PostgreSQL
│   ├───storage            # фунции для подлючения к внешним источникам
│   └───usecase            # бизнес-логика приложения (сервисы, сценарии работы, интерфейсы)
│
├───migrations             # миграции базы данных (SQL-файлы для goose)
│
├───templates              # HTML-шаблон для UI
│
├───docs                   # документация проекта (скрины)
│
└───tools
    └───goose              # инструмент для миграций БД (goose binary/config)
```

## Возможности

- Посылать заказы в топик.
- Обрабатывать сообщения в консьюмере и сохранять в БД.
- Запрашивать информацию о заказе по uuid.
- Отображать информацию о заказе в простом HTML-интерфейсе.

## Особенности

- Чистая архитектура.
- Чтение заказов из топика (Kafka), гарантия доставки - at-least-once.
- Пайплайн работы консьюмера: FetchMessage -> Unmarshal -> Validate -> DTO -> ProcessIncomingOrder -> CommitMessages.
- Заказы хранятся в БД (PostgreSQL) и дополнительно кэшируются in-memory (list, map, mutex).
- LRU cache, работающий по принципам read & write aside.
- Кэш при запуске сервиса "прогревается" заданным количеством последних заказов из БД.
- usecase-слой и кэш покрыты тестами.
- DTO-mapping при чтении сообщений из консьюмера и при отдаче по запросу.
- Конфиг из переменных окружения.
- Все поднимается в контейнерах через `docker-compose`, приложение запускается после доступности БД и брокера сообщений.
- Миграции БД с помощью `goose`, для которого поднимается отдельный контейнер со скриптом.
- Запросы обернуты в retry-функцию, запрос сохранения заказа в несколько таблиц обернут в транзакцию.
- В данные о заказе добавлен атрибут `quantity` для нормализации схемы.

## Библиотеки

- Логгирование - `log/slog`.
- Роутинг - `gorilla/mux`.
- Драйвер PostgreSQL - `jackc/pgx/v5/stdlib`.
- UUID - `google/uuid`
- Генератор данных - `brianvoe/gofakeit/v7`
- Драйвер Kafka - `segmentio/kafka-go`.
- Подгрузка конфига из .env - `caarlos0/env`.
- Тесты и моки - `stretchr/testify`.

## Сборка и тестирование

1. Конфигурация .env файла

```shell
vim .env.example
mv .env.example .env
```

```shell
KAFKA_BROKER=kafka:9092                 # адрес брокера Kafka
KAFKA_GET_ORDER_TOPIC=get_orders        # имя топика для заказов
KAFKA_CONSUMER_GROUP=default            # группа консьюмеров
KAFKA_BACKOFF=500ms                     # задержка перед повторной попыткой
KAFKA_MAX_WAIT=500ms                    # макс. ожидание сообщений

PG_DSN=postgres://app:app@postgres:5432/?sslmode=disable   # строка подключения
PG_SAVE_TIMEOUT=5s                  # таймаут сохранения заказа
PG_EXISTS_TIMEOUT=2s                # таймаут проверки существования
PG_GET_TIMEOUT=2s                   # таймаут получения заказа
PG_POOL_MAX_OPEN_CONNS=10           # макс. открытых соединений
PG_POOL_MAX_IDLE_CONNS=5            # макс. неактивных соединений
PG_POOL_CONN_MAX_LIFETIME=1h        # срок жизни соединения
PG_PING_TIMEOUT=500ms               # таймаут пинга БД
PG_MAX_RETRIES=3                    # число повторных попыток
PG_BACKOFF=500ms                    # пауза между ретраями

SERVER_HTTP_PORT=8080               # порт сервера
SERVER_HTTP_SHUTDOWN_TIMEOUT=5s     # время на корректное завершение
SERVER_HTTP_READ_HEADER_TIMEOUT=5s  # таймаут чтения заголовков
SERVER_HTTP_READ_TIMEOUT=5s         # таймаут чтения запроса
SERVER_HTTP_WRITE_TIMEOUT=10s       # таймаут записи ответа
SERVER_HTTP_IDLE_TIMEOUT=120s       # таймаут idle-соединений

CACHE_CAPACITY=30                   # вместимость кэша
CACHE_WARM_SIZE=15                  # предзагрузка заказов при старте
```

2. Тестирование (unit, integration)

```shell
go test ./... -v
```

3. Сборка образов и локальная развертка

```shell
docker compose up --build -d
```

4. Отправка заказов в топик (30 заказов)

```shell
go run cmd/producer/main.go
```

5. Тестирование (UI)

```text
http://localhost:8080/templates/
```

![Скрин](docs/screen1.png)

6. Тестирование (Postman)

`GET /order/{uid}`

_request_
```text
empty
```

_response_

`200` - заказ найден

```json
{
  "order_uid": "f7cc03e5-d018-4164-9057-99763d2b6622",
  "track_number": "99cf7d68-7613-4385-b6c3-76043cc5540e",
  "delivery": {
    "name": "Isabella Johnston",
    "phone": "5408680102",
    "zip": "33501",
    "city": "Omaha",
    "address": "jarvisbotsford@ullrich.name",
    "region": "Oklahoma",
    "email": "mariannabergnaum@klocko.com"
  },
  "payment": {
    "currency": "SOS",
    "amount": 16,
    "delivery_cost": 96
  },
  "items": [
    {
      "price": 54,
      "name": "Sharp Biometric Thermometer",
      "sale": 77,
      "size": "PublicEngines",
      "total_price": 66,
      "brand": "74",
      "quantity": 2
    },
    {
      "price": 8,
      "name": "Zen Smart Home Device X",
      "sale": 37,
      "size": "Healthline",
      "total_price": 18,
      "brand": "85",
      "quantity": 4
    },
    {
      "price": 83,
      "name": "Purple Bike Turbo",
      "sale": 47,
      "size": "Porch",
      "total_price": 81,
      "brand": "38",
      "quantity": 1
    }
  ],
  "delivery_service": "nYANH",
  "date_created": "2025-08-16T22:02:39Z"
}
```

`400` - некорректный uuid заказа

```text
invalid UUID length: 35
```

`404` - заказ не найден

```text
order does not exists
```

## Схема данных

![Схема](docs/screen2.png)