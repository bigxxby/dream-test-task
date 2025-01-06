# dream-test-task

url shortener МЕЧТЫ

## Структура проекта

```
├───cmd
├───docs
└───internal
    ├───api
    │   ├───middleware
    │   ├───repo
    │   │   ├───auth
    │   │   ├───shortener
    │   │   └───user
    │   ├───service
    │   │   ├───auth
    │   │   └───shortener
    │   └───transport
    │       ├───auth
    │       └───shortener
    ├───app
    ├───config
    ├───database
    │   ├───connection
    │   └───migration
    ├───models
    ├───router
    └───utils
```

## Установка и запуск

### 1. Клонируйте репозиторий

```bash
git clone https://github.com/bigxxby/dream-test-task.git
cd dream-test-task
```

### 2. Настройка окружения (есть пример .env)

```
DB_HOST=localhost
DB_PORT=5432
DB_USER=your-db-user
DB_PASSWORD=your-db-password
DB_NAME=your-db-name
JWT_SECRET=your-jwt-secret
APP_PORT=8081
```

### 3. Сборка и запуск с использованием Docker

Проект использует Docker и Docker Compose для упрощения развертывания.

## 3.1. Запуск через контейнер

docker-compose up --build

### 4.1. Установите зависимости

go mod tidy

## 4.2. Скомпилируйте приложение

go build -o app ./cmd

## 4.3. Запустите приложение

```
go run ./cmd
```

### 5. Доступ к Swagger UI

Swagger UI доступен по следующему маршруту:

```
http://localhost:8081/swagger/index.html
```

### 6. Маршруты API

```
/auth
POST /register — Регистрация нового пользователя.
POST /login — Вход в систему.
GET /whoami — Получение информации о текущем пользователе (необходима аутентификация).
```

```
/shortener
GET / — Получение всех сокращенных ссылок пользователя (необходима аутентификация).
GET /:shortID — Редирект на оригинальную ссылку по сокращенному идентификатору.
GET /stats/:shortID — Получение статистики по сокращенной ссылке.
POST / — Создание новой сокращенной ссылки (необходима аутентификация).
DELETE /:shortID — Удаление сокращенной ссылки (необходима аутентификация).
```
