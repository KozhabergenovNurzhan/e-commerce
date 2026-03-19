# E-Commerce REST API

REST API для интернет-магазина на Go.

## Стек

- **Go** + **Gin** — HTTP сервер
- **PostgreSQL** + **sqlx** + **pgx** — база данных
- **JWT** — аутентификация
- **Docker** + **docker-compose** — контейнеризация
- **goose** — миграции

## Структура проекта

```
├── cmd/api/          — точка входа
├── internal/
│   ├── auth/         — JWT генерация и парсинг
│   ├── config/       — конфигурация из .env
│   ├── handler/      — HTTP handlers
│   ├── middleware/   — JWT, CORS, Logger
│   ├── models/       — модели и ошибки
│   ├── pkg/          — shared утилиты
│   ├── repository/   — слой БД
│   ├── server/       — роуты и gin setup
│   └── service/      — бизнес логика
└── migrations/       — SQL миграции
```

## Запуск

### С Docker

```bash
docker-compose up --build
```

### Локально

```bash
# 1. Запустить postgres
docker-compose up postgres

# 2. Создать .env
cp .env.example .env

# 3. Применить миграции
goose -dir migrations postgres "host=localhost port=5432 user=postgres password=postgres dbname=ecommerce sslmode=disable" up

# 4. Запустить сервер
go run ./cmd/api
```

## API Endpoints

### Auth
| Method | Path | Описание |
|--------|------|----------|
| POST | `/api/v1/auth/register` | Регистрация |
| POST | `/api/v1/auth/login` | Логин → JWT |

### Products
| Method | Path | Описание |
|--------|------|----------|
| GET | `/api/v1/products` | Список товаров (поиск, фильтр) |
| GET | `/api/v1/products/:id` | Один товар |
| POST | `/api/v1/admin/products` | Создать товар (admin) |
| PUT | `/api/v1/admin/products/:id` | Обновить товар (admin) |
| DELETE | `/api/v1/admin/products/:id` | Удалить товар (admin) |

### Categories
| Method | Path | Описание |
|--------|------|----------|
| GET | `/api/v1/categories` | Список категорий |
| POST | `/api/v1/admin/categories` | Создать категорию (admin) |
| PUT | `/api/v1/admin/categories/:id` | Обновить категорию (admin) |
| DELETE | `/api/v1/admin/categories/:id` | Удалить категорию (admin) |

### Cart
| Method | Path | Описание |
|--------|------|----------|
| GET | `/api/v1/cart` | Моя корзина |
| POST | `/api/v1/cart` | Добавить товар |
| DELETE | `/api/v1/cart/:product_id` | Убрать товар |

### Orders
| Method | Path | Описание |
|--------|------|----------|
| POST | `/api/v1/orders` | Оформить заказ из корзины |
| GET | `/api/v1/orders` | Мои заказы |
| GET | `/api/v1/orders/:id` | Один заказ |

## Примеры запросов

### Регистрация
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"123456"}'
```

### Логин
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"123456"}'
```

### Список товаров с поиском
```bash
curl 'http://localhost:8080/api/v1/products?search=iphone&category_id=1&limit=10&offset=0'
```

### Добавить в корзину
```bash
curl -X POST http://localhost:8080/api/v1/cart \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"product_id":1,"quantity":2}'
```

### Оформить заказ
```bash
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Authorization: Bearer <token>"
```


