# Price Tracker Service 

Микросервис для асинхронного отслеживания цен на товары в интернет-магазинах 

## 🛠 Технологический стек

* **Язык:** Go (Golang)
* **Архитектура:** Clean Architecture (Domain, Service, Repository, Delivery)
* **API:** gRPC
* **База данных:** PostgreSQL
* **Брокер сообщений:** Apache Kafka
* **Кэширование:** Redis 
* **Инфраструктура:** Docker, Docker Compose
* **Парсинг:** goquery

## ⚙️ Функционал

1. Добавление товара для отслеживания (URL, желаемая цена).
2. Фоновый Воркер (Cron-like), который с заданной периодичностью обходит ссылки и парсит актуальные цены.
3. Сохранение истории цен в PostgreSQL.
4. Отправка асинхронного уведомления в Kafka, если текущая цена опустилась ниже целевой (Target Price).

## 📂 Структура проекта

Проект спроектирован с учетом разделения зон ответственности:
- `cmd/server/` — точка входа в приложение, инициализация зависимостей (DI).
- `internal/domain/` — основные бизнес-сущности (Product, Subscription, User).
- `internal/delivery/` — транспортный слой (gRPC хендлеры).
- `internal/service/` — бизнес-логика и оркестрация.
- `internal/repository/` — слой работы с данными (Postgres SQL).
- `internal/infrastructure/` — внешние зависимости (Парсер сайтов, Kafka Producer).
- `internal/worker/` — фоновый процесс (Price Fetcher) для актуализации цен.

## 🚀 Как запустить

1. Клонируйте репозиторий.
2. Убедитесь, что у вас установлен Docker и Docker Compose.
3. Запустите инфраструктуру (PostgreSQL, Zookeeper, Kafka, Redis):
   ```bash
   docker-compose up -d
4. Примените SQL-миграции из папки /migrations к базе данных Price
5. Запустите gRPC сервер и фоновый воркер: go run ./cmd/server/main.go
   


