# golang-advert-app

Бэкенд сервиса объявлений на Go. Один процесс поднимает **два транспорта одновременно** — REST (HTTP) и gRPC — поверх общего слоя бизнес-логики. Хранилище — PostgreSQL, изображения объявлений лежат в локальной файловой системе.

## Возможности

- **Объявления** — создание, редактирование, просмотр, список с фильтрами, счётчик, архивация и удаление.
- **Модерация** — объявления проходят подтверждение (`Approve`) / отклонение (`Reject`) администратором.
- **Изображения** — загрузка и отдача картинок объявления через gRPC-стримы.
- **Категории** — справочник категорий (управление доступно только администратору).
- **Избранное** — добавление объявлений в избранное, список и счётчик.
- **Пользователи и аутентификация** — регистрация, логин, logout, обновление токенов на JWT (access/refresh), роли (`user` / `admin`).
- **Векторный поиск** — поиск объявлений по векторному представлению (миграция `000005_vector_search`).

## Технологии

- **Go 1.26**
- **gRPC** / **protobuf** — `google.golang.org/grpc`, `protoc`
- **net/http** — собственный роутер с версионированием API (`/api/v1`)
- **PostgreSQL** через `pgx/v5` (пул соединений)
- **JWT** — `golang-jwt/v5`
- **Валидация** — `go-playground/validator/v10`
- **Конфигурация** — `kelseyhightower/envconfig` (из переменных окружения)
- **Логирование** — `go.uber.org/zap`
- **Тесты** — `stretchr/testify`, моки через `go.uber.org/mock` (mockgen)
- **Миграции** — `migrate/migrate`
- Инфраструктура — Docker Compose

## Архитектура

Код организован по фичам (`internal/features/<feature>`), каждая фича разбита на слои:

```
internal/features/advert/
├── entity/       # доменные сущности и правила валидации
├── usecase/      # бизнес-логика (не знает про транспорт)
├── repo/
│   ├── postgres/ # доступ к БД
│   └── local/    # хранилище изображений на диске
└── controller/
    ├── rest/     # HTTP-хендлеры
    └── grpc/     # gRPC-хендлеры
```

Оба контроллера (`rest` и `grpc`) вызывают один и тот же `usecase` — транспорт отделён от логики.

Общий инфраструктурный код — в `internal/core` (config, logger, postgres pool, http/grpc серверы и middleware/интерсепторы, jwt, validator, roles, errors). Protobuf-схемы — в `internal/protos`, сгенерированный код — в `internal/gen`.

## Требования

- Go 1.26+
- Docker и Docker Compose
- `protoc` с плагинами `protoc-gen-go` и `protoc-gen-go-grpc` — только если нужно перегенерировать gRPC-код

## Быстрый старт

1. Скопируйте пример конфигурации и заполните значения:

   ```bash
   cp .env.example .env
   ```

   Переменные окружения:

   | Переменная | Назначение |
   |---|---|
   | `TIME_ZONE` | Часовой пояс приложения |
   | `HTTP_ADDR` / `HTTP_SHUTDOWN_TIMEOUT` | Адрес и таймаут graceful-shutdown HTTP-сервера |
   | `GRPC_ADDR` / `GRPC_SHUTDOWN_TIMEOUT` / `GRPC_REFLECTION` | Адрес, таймаут и reflection gRPC-сервера |
   | `POSTGRES_USER` / `POSTGRES_PASSWORD` / `POSTGRES_DB` / `POSTGRES_TIMEOUT` | Параметры подключения к PostgreSQL |
   | `JWT_ACCESS_SECRET` | Секрет для подписи JWT |

2. Поднимите базу данных:

   ```bash
   make env-up
   ```

3. Пробросьте порт PostgreSQL на `localhost` — приложение запускается на хосте и
   ходит в БД по `localhost:5432` (контейнер БД сам порт наружу не публикует):

   ```bash
   make env-port-forward
   ```

4. Примените миграции:

   ```bash
   make migrate-up
   ```

5. Запустите приложение:

   ```bash
   make app-run
   ```

   Стартуют оба сервера — HTTP и gRPC. Логи пишутся в `out/logs`.

   > Миграции (`make migrate-up`) выполняются внутри docker-сети и проброс порта
   > не требуют — он нужен именно приложению (`make app-run`).

## Команды Makefile

| Команда | Что делает |
|---|---|
| `make env-up` | Запустить контейнер PostgreSQL |
| `make env-down` | Остановить PostgreSQL |
| `make env-cleanup` | Удалить контейнеры и **volume с данными** (спросит подтверждение) |
| `make env-port-forward` / `make env-port-close` | Проброс порта 5432 на localhost |
| `make migrate-create seq=<name>` | Создать новую миграцию |
| `make migrate-up` / `make migrate-down` | Применить / откатить миграции |
| `make grpc-gen-advert` (`-user`, `-category`, `-favourite`) | Перегенерировать gRPC-код из `.proto` |
| `make app-run` | Запустить приложение |

## API

Оба транспорта предоставляют одинаковый набор операций.

- **REST** — под префиксом версии, например `/api/v1/adverts`, `/api/v1/adverts/categories` и т.д.
- **gRPC** — сервисы `Advert`, `User`, `Category`, `Favourite` (см. `internal/protos`). При `GRPC_REFLECTION=true` доступна reflection для `grpcurl`/Postman.

### Аутентификация и роли

Большинство методов требуют валидный JWT access-токен. Публичные операции (регистрация, логин, обновление токенов, публичный просмотр категорий) доступны без токена. Часть действий — управление категориями, `Approve`/`Reject` объявлений — требуют роль `admin`. Авторизация реализована в gRPC-интерсепторах (`JWToken`, `Role`) и аналогичном HTTP-middleware.

- **REST** передаёт токен заголовком `Authorization: Bearer <access_token>`.
- **gRPC** передаёт токен в метаданных `authorization: bearer <access_token>`.

## Примеры запросов (happy path)

Данные в примерах — обезличенные, замените на свои. Значения таймаутов/адресов берутся из `.env`; ниже HTTP считается поднятым на `localhost:8080`.

```bash
# Базовый URL REST API и переменные для примеров
export BASE="http://localhost:8080/api/v1"
export TOKEN="<access_token>"   # заполнится после логина
```

### REST

#### Пользователи и аутентификация

```bash
# Регистрация (публично)
curl -X POST "$BASE/users" \
  -H 'Content-Type: application/json' \
  -d '{
    "email": "john.doe@example.com",
    "full_name": "John Doe",
    "phone_number": "+10000000000",
    "password": "SuperSecret123"
  }'

# Логин (публично) — возвращает user_id, access_token, refresh_token
curl -X POST "$BASE/auth/login" \
  -H 'Content-Type: application/json' \
  -d '{
    "email": "john.doe@example.com",
    "password": "SuperSecret123"
  }'

# Обновление пары токенов (публично)
curl -X POST "$BASE/auth/refresh" \
  -H 'Content-Type: application/json' \
  -d '{"refresh_token": "<refresh_token>"}'

# Logout (нужен access-токен)
curl -X POST "$BASE/auth/logout" \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"refresh_token": "<refresh_token>"}'

# Список пользователей с пагинацией
curl "$BASE/users?limit=20&offset=0" \
  -H "Authorization: Bearer $TOKEN"

# Пользователь по id
curl "$BASE/users/1" \
  -H "Authorization: Bearer $TOKEN"
```

#### Категории

```bash
# Создать категорию (admin). parent_id опционален (null для корневой)
curl -X POST "$BASE/adverts/categories" \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"parent_id": null, "name": "Electronics"}'

# Список категорий (публично)
curl "$BASE/adverts/categories"

# Изменить категорию (admin) — все поля опциональны
curl -X PATCH "$BASE/adverts/categories/1" \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"name": "Home Electronics"}'

# Удалить категорию (admin)
curl -X DELETE "$BASE/adverts/categories/1" \
  -H "Authorization: Bearer $TOKEN"
```

#### Объявления

```bash
# Создать объявление — multipart/form-data (можно с картинками)
curl -X POST "$BASE/adverts" \
  -H "Authorization: Bearer $TOKEN" \
  -F 'title=Bicycle' \
  -F 'description=Almost new city bike' \
  -F 'price=15000' \
  -F 'category_id=1' \
  -F 'images=@/path/to/photo1.jpg' \
  -F 'images=@/path/to/photo2.jpg'

# Объявление по id
curl "$BASE/adverts/1" \
  -H "Authorization: Bearer $TOKEN"

# Список объявлений с фильтрами и пагинацией.
# Доступные query-параметры: limit, offset, search_query, min_price, max_price,
# category_id, status, from_date, to_date, sort, order.
# Формат дат: "YYYY-MM-DD HH:MM:SS".
curl -G "$BASE/adverts" \
  -H "Authorization: Bearer $TOKEN" \
  --data-urlencode 'limit=20' \
  --data-urlencode 'offset=0' \
  --data-urlencode 'search_query=bike' \
  --data-urlencode 'min_price=1000' \
  --data-urlencode 'max_price=50000' \
  --data-urlencode 'category_id=1' \
  --data-urlencode 'sort=price' \
  --data-urlencode 'order=asc'

# Мои объявления
curl "$BASE/adverts/my?limit=20&offset=0" \
  -H "Authorization: Bearer $TOKEN"

# Количество объявлений (принимает те же фильтры, что и список)
curl -G "$BASE/adverts/count" \
  -H "Authorization: Bearer $TOKEN" \
  --data-urlencode 'category_id=1'

# Частичное обновление. version обязателен (оптимистичная блокировка),
# остальные поля опциональны.
curl -X PATCH "$BASE/adverts/1" \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{
    "version": 1,
    "title": "Bicycle (red)",
    "price": 14000
  }'

# Подтвердить объявление (admin)
curl -X POST "$BASE/adverts/1/approve" \
  -H "Authorization: Bearer $TOKEN"

# Отклонить объявление (admin)
curl -X POST "$BASE/adverts/1/reject" \
  -H "Authorization: Bearer $TOKEN"

# Архивировать объявление (владелец)
curl -X POST "$BASE/adverts/1/archive" \
  -H "Authorization: Bearer $TOKEN"

# Удалить объявление (владелец)
curl -X DELETE "$BASE/adverts/1" \
  -H "Authorization: Bearer $TOKEN"
```

#### Изображения объявлений

```bash
# Догрузить изображения к объявлению — multipart/form-data
curl -X POST "$BASE/adverts/images" \
  -H "Authorization: Bearer $TOKEN" \
  -F 'advert_id=1' \
  -F 'images=@/path/to/photo.jpg'

# Получить изображение по id (бинарный поток) — сохранить в файл
curl "$BASE/adverts/images/1" \
  -H "Authorization: Bearer $TOKEN" \
  --output image.jpg

# Удалить изображение по id
curl -X DELETE "$BASE/adverts/images/1" \
  -H "Authorization: Bearer $TOKEN"
```

#### Избранное

```bash
# Добавить объявление в избранное (id объявления в пути)
curl -X POST "$BASE/adverts/favourites/1" \
  -H "Authorization: Bearer $TOKEN"

# Список избранных объявлений с пагинацией
curl "$BASE/adverts/favourites?limit=20&offset=0" \
  -H "Authorization: Bearer $TOKEN"

# Количество избранных
curl "$BASE/adverts/favourites/count" \
  -H "Authorization: Bearer $TOKEN"

# Только id избранных объявлений
curl "$BASE/adverts/favourites/ids" \
  -H "Authorization: Bearer $TOKEN"

# Убрать объявление из избранного (id объявления в пути)
curl -X DELETE "$BASE/adverts/favourites/1" \
  -H "Authorization: Bearer $TOKEN"
```

### gRPC

Примеры на [`grpcurl`](https://github.com/fullstorydev/grpcurl). При `GRPC_REFLECTION=true` схему указывать не нужно — она берётся через reflection. Токен передаётся метаданными `authorization: bearer <token>`.

```bash
export GRPC_ADDR="localhost:9090"

# Список сервисов и методов (reflection)
grpcurl -plaintext "$GRPC_ADDR" list
grpcurl -plaintext "$GRPC_ADDR" describe user.User

# Регистрация (публично)
grpcurl -plaintext -d '{
  "email": "john.doe@example.com",
  "full_name": "John Doe",
  "phone_number": "+10000000000",
  "password": "SuperSecret123"
}' "$GRPC_ADDR" user.User/Create

# Логин (публично)
grpcurl -plaintext -d '{
  "email": "john.doe@example.com",
  "password": "SuperSecret123"
}' "$GRPC_ADDR" user.User/Login

# Защищённый вызов — с токеном в метаданных
grpcurl -plaintext \
  -H "authorization: bearer $TOKEN" \
  -d '{"limit": 20, "offset": 0}' \
  "$GRPC_ADDR" advert.Advert/List
```

> Остальные gRPC-методы (`Advert`, `Category`, `Favourite`, `User`) вызываются по той же схеме — точные имена полей запроса смотрите через `grpcurl ... describe <Service>` или в `internal/protos`.

## Тесты

Юнит-тесты покрывают usecase- и controller-слои (табличные тесты, моки зависимостей через mockgen):

```bash
go test ./...
```

## Структура проекта

```
cmd/app/            # точка входа, сборка зависимостей и запуск серверов
internal/
├── core/           # инфраструктура: config, logger, postgres, http, grpc, jwt, ...
├── features/       # фичи: advert, category, favourite, image, user
├── protos/         # protobuf-схемы
└── gen/            # сгенерированный gRPC-код
migrations/         # SQL-миграции
docker-compose.yaml # PostgreSQL, migrate, port-forwarder
Makefile            # команды для окружения, миграций, генерации и запуска
```
