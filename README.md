# Sensory Navigator: Users and Reviews Module

## Описание проекта

Система анализа торговых точек по сенсорным характеристикам. Данный модуль отвечает за работу с пользователями и отзывами.

### Основные функции:
- Регистрация и авторизация пользователей (JWT)
- Управление профилем пользователя
- Создание, редактирование и удаление отзывов
- Работа с избранным
- Оценка мест по сенсорным критериям:
  - Сенсорность (1-5)
  - Освещение (1-5)
  - Громкость звука (1-5)
  - Плотность людей (1-5)
  - Доступность (1-5)

## Технологический стек

### Backend
- **Go** - язык программирования
- **Gin** - web-фреймворк
- **PostgreSQL** - база данных
- **JWT** - авторизация

### Frontend
- **Tauri** - desktop-приложение
- **HTML/CSS/JavaScript** - веб-интерфейс
- **Vite** - сборщик

## Структура проекта

```
sensory-navigator/
├── backend/
│   ├── config/         # Конфигурация приложения
│   ├── database/       # Подключение к БД и миграции
│   ├── handlers/       # HTTP обработчики
│   ├── middleware/     # Middleware (авторизация)
│   ├── models/         # Модели данных
│   ├── repository/     # Работа с БД
│   ├── services/       # Бизнес-логика
│   └── main.go         # Точка входа
├── frontend/
│   ├── src/            # Исходный код (HTML, CSS, JS)
│   ├── src-tauri/      # Tauri конфигурация
│   ├── package.json    # NPM зависимости
│   └── vite.config.js  # Конфигурация Vite
└── README.md
```

## API Endpoints

### Авторизация
- `POST /api/auth/register` - Регистрация
- `POST /api/auth/login` - Вход
- `POST /api/auth/refresh` - Обновление токена
- `POST /api/auth/forgot-password` - Запрос восстановления пароля
- `POST /api/auth/reset-password` - Сброс пароля

### Пользователи
- `GET /api/users/me` - Получить профиль
- `PUT /api/users/me` - Обновить профиль
- `GET /api/users/me/reviews` - Мои отзывы
- `GET /api/users/me/favorites` - Моё избранное

### Отзывы
- `GET /api/places/:id/reviews` - Отзывы места
- `POST /api/places/:id/reviews` - Создать отзыв
- `PUT /api/reviews/:id` - Редактировать отзыв
- `DELETE /api/reviews/:id` - Удалить отзыв

### Избранное
- `POST /api/favorites/:placeId` - Добавить в избранное
- `DELETE /api/favorites/:placeId` - Удалить из избранного

## Установка и запуск

### Backend

1. Установите Go (версия 1.21+)
2. Установите PostgreSQL
3. Создайте базу данных:
   ```sql
   CREATE DATABASE sensory_navigator;
   ```
4. Выполните миграции:
   ```bash
   psql -d sensory_navigator -f backend/database/migrations/001_init.sql
   ```
5. Скопируйте env.example.txt в .env и настройте параметры
6. Запустите backend:
   ```bash
   cd backend
   go mod download
   go run main.go
   ```

### Frontend

1. Установите Node.js (версия 18+)
2. Установите Rust (для Tauri)
3. Установите зависимости:
   ```bash
   cd frontend
   npm install
   ```
4. Запуск в режиме разработки:
   ```bash
   npm run dev
   ```
5. Сборка Tauri приложения:
   ```bash
   npm run tauri build
   ```

## Авторы

Насрулаев Ш.М.  
Атаханов Н.Р.
Направление: 09.03.04 Программная инженерия  
Высшая школа экономики

## Лицензия

Этот проект является учебной работой.

