# 📋 TodoApp

Простое ToDo API-приложение на Go для создания задач на день

## ✨ Возможности

- CRUD для задач
- Фильтрация по дате и статусу
- Пагинация по статусу
- Swagger-документация
- Локальный и Docker-режимы
- Покрытие тестами
- Makefile

---

## 📦 Стек

- **Go**
- **Fiber**
- **GORM**
- **zap (logger)**
- **PostgreSQL**
- **Docker + docker-compose**
- **Swagger (swaggo)**
- **golangci-lint**
- **Makefile** 

---

## ⚖️ Установка и запуск

### 🛠 Требования

- Go 1.24.2
- Docker + Docker Compose ( для контейнерного запуска )
- PostgreSQL ( локально )
- `golangci-lint`, `swag`

### 💾 Локальный запуск

```bash
git clone git@github.com:volchok96/todoapp.git
cd todoapp
go mod tidy
make db-create db-migrate db-seed
make run
```

Приложение: [http://localhost:8080](http://localhost:8080)

Swagger: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)


### 🐳 Docker

```bash
make docker-up
```

Приложение: [http://localhost:8081](http://localhost:8081)

Swagger: [http://localhost:8081/swagger/index.html](http://localhost:8081/swagger/index.html)

Остановка:

```bash
make docker-down
```

---

## 🧪 Тесты

```bash
make test
```

---

## 📚 Swagger-документация

Генерация:

```bash
make go-generate
```

---

## 🔍 Линт

```bash
make lint
```

---

## 🛠 Makefile команды

| Команда             | Описание                                   |
|---------------------|--------------------------------------------|
| `make run`          | Запустить приложение локально              |
| `make build`        | Собрать бинарник                           |
| `make test`         | Запустить тесты                            |
| `make lint`         | Запустить линтер                           |
| `make go-generate`  | Сгенерировать Swagger-доки в `docs/`       |
| `make docker-up`    | Поднять Docker-контейнеры                  |
| `make docker-down`  | Остановить Docker-контейнеры                      |
| `make docker-logs-app` | Логи приложения из Docker               |
| `make docker-psql`  | Подключиться к базе PostgreSQL в Docker     |
| `make db-create`    | Создать базу данных                        |
| `make db-migrate`   | Выполнить миграции                         |
| `make db-seed`      | Заполнить БД тестовыми данными                    |
| `make db-reset`     | Сбросить и заново инициализировать БД|
| `make local-start`   | Создать БД с тестовыми данными и запустить локально|
---

