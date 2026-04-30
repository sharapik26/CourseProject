# Сенсорный навигатор — Модуль пользователей и отзывов

Индивидуальная часть командного курсового проекта.

> ОП «Программная инженерия», ФКН НИУ ВШЭ, 2025/2026.
> Исполнитель: **Насрулаев Ш. М.** (БПИ234).
> Научный руководитель: А. К. Бегичева, ст. преп. департамента программной инженерии.

Этот репозиторий содержит **только** модуль пользователей, авторизации,
отзывов и избранного. После прохождения индивидуальной защиты модуль
объединяется с модулем карты и мест (исполнитель: Атаханов Н. Р.) в
объединённый репозиторий `Project/`.

## Технологический стек

### Backend
- Go 1.21+, Gin, GORM, PostgreSQL 16
- JWT (`github.com/golang-jwt/jwt/v5`), bcrypt
- OpenAPI 3.0 спецификация: `backend/docs/openapi.yaml`

### Frontend
- Vue 3 + Pinia + Vue Router
- Vite + TypeScript

## Состав модуля

| Слой | Файл | Назначение |
|---|---|---|
| Модели | `backend/internal/models/models.go` | `User`, `Review`, `Favorite`, `PlaceRef` (стаб) |
| Сервисы | `backend/internal/services/users.go` | Регистрация, аутентификация, профиль |
| Сервисы | `backend/internal/services/reviews.go` | Отзывы, избранное, агрегаты |
| HTTP | `backend/internal/handlers/auth_handler.go` | `/api/auth/register`, `/api/auth/login` |
| HTTP | `backend/internal/handlers/users_handler.go` | `/api/users/me`, `/api/users/me/password` |
| HTTP | `backend/internal/handlers/reviews_handler.go` | `/api/reviews`, `/api/places/:id/reviews`, `/api/favorites/:id` |
| Auth | `backend/internal/auth/jwt.go` | Выпуск/проверка JWT, bcrypt |
| Middleware | `backend/internal/middleware/auth.go` | `RequireAuth` для защиты эндпоинтов |
| Frontend | `frontend/src/views/{Login,Register,Profile,MyReviews,Favorites,NewReview}View.vue` | UI |

## Запуск

### Через Docker Compose

```bash
docker compose up -d --build
```

API доступен на `http://localhost:8081`. PostgreSQL — на `localhost:5433`.

Демо-пользователь после `cmd/seed`: `demo@example.com / demo123`.

```bash
cd backend
go run ./cmd/seed
```

### Локально (с уже работающим PostgreSQL)

```bash
cd backend
cp .env.example .env       # отредактируйте под свою БД
go run ./cmd/server
```

### Frontend

```bash
cd frontend
npm install
npm run dev                # http://localhost:1421
```

## Тестирование

```bash
cd backend
go test ./...
```

Юнит-тесты `internal/auth/` и `internal/services/` запускаются без БД.
Интеграционные тесты `tests/` требуют `TEST_DB_DSN`.

## Интеграция с модулем карты

Модуль предоставляет два публичных HTTP-API, потребляемых модулем карты:

1. **Middleware `RequireAuth`** (Go-импорт: `internal/middleware`) — проверка
   JWT-токена. Используется модулем карты для защиты `POST/PUT/DELETE /api/places`.
2. **GET `/api/places/{id}/aggregate`** — публичный эндпоинт, возвращающий
   средние сенсорные оценки места. Используется модулем карты для отображения
   агрегатов в списке и карточке места.
3. **GET `/api/places/{id}/reviews`** — публичный список отзывов места.

В индивидуальном стенде модуля используется упрощённая таблица `places`
(только `id` + `name`), создаваемая `database.SeedDemoPlaces`. При объединении
с модулем карты эта таблица заменяется полной таблицей мест с координатами
и категориями.

## Лицензия

MIT.