//go_dnd

Backend для проведения пошаговых боёв по правилам D&D 5e. Мастер работает через REST API, игроки - через Telegram-бота.

Пет-проект на Go. Пишу, чтобы разобраться на практике с бэкендом: слоями, PostgreSQL, JWT, тестами.

//Что сделано

Спринт 1 закончен:

- Docker Compose с PostgreSQL 16 и Redis 7, healthcheck, volumes
- Конфигурация из переменных окружения через cleanenv
- JSON-логи на slog, у каждого запроса свой request_id, возвращается в заголовке X-Request-ID
- Graceful shutdown: по Ctrl+C сервер дожидается активных запросов и завершается
- Миграции лежат в бинарнике через embed.FS, применяются при старте через golang-migrate
- Регистрация Мастера: bcrypt для пароля, JWT для сессии
- GET /api/v1/me - возвращает данные текущего пользователя

Спринт 2 в работе:

- Telegram-бот подключён, работает в отдельной горутине, останавливается вместе с приложением
- /start регистрирует игрока: сохраняет telegram_id и telegram_username
- Если username занят другим аккаунтом, игрок создаётся без username

//Стек

- Go 1.26
- chi - HTTP-роутер
- pgx v5 - драйвер PostgreSQL
- PostgreSQL 16, Redis 7
- golang-migrate - миграции
- golang-jwt - JWT (HS256)
- go-telegram/bot - Telegram-бот
- slog - логирование
- Docker Compose

//Как запустить

Нужны Go 1.26+, Docker с Compose и токен бота от @BotFather.

Создать .env из шаблона:

Copy-Item .env.example .env

Заполнить в .env:

- POSTGRES_PASSWORD
- REDIS_PASSWORD
- JWT_SECRET
- TELEGRAM_BOT_TOKEN

Поднять PostgreSQL и Redis:

docker compose -f deploy/docker-compose.yml --env-file .env up -d

Запустить приложение:

go run ./cmd/app

Миграции применятся автоматически при старте.

Проверить API. В PowerShell используйте curl.exe, а не curl - curl это алиас на Invoke-WebRequest.

Регистрация Мастера:

curl.exe -X POST http://localhost:8080/api/v1/auth/register -H "Content-Type: application/json" -d '{"username":"master","password":"password123"}'

Логин, вернёт токен:

curl.exe -X POST http://localhost:8080/api/v1/auth/login -H "Content-Type: application/json" -d '{"username":"master","password":"password123"}'

Текущий пользователь, подставить токен из ответа логина:

curl.exe http://localhost:8080/api/v1/me -H "Authorization: Bearer <токен>"

Тесты:

go test ./...

//API

| Метод | Путь                  | Описание |
| GET   | /healthz              | проверка работы |
| POST  | /api/v1/auth/register | регистрация Мастера |
| POST  | /api/v1/auth/login    | вход, возвращает JWT |
| GET   | /api/v1/me            | текущий пользователь, требует Authorization: Bearer ... |

//Telegram-бот

Пока работает только /start — регистрирует игрока. Если username занят другим аккаунтом, регистрирует без него. И выводит сообщение с просьбой создать username.

//Структура

cmd/app/                 main-функция
internal/config/         конфигурация из окружения
internal/domain/         сущности, ошибки, интерфейсы
internal/usecase/        бизнес-логика
internal/repository/     работа с PostgreSQL
internal/transport/      REST и Telegram
internal/migrations/     SQL-миграции (вшиты в бинарник)
internal/pkg/            hasher, token, logger
deploy/                  docker-compose
docs/                    планы спринтов, техдолг

//Тесты

- hasher, token - unit-тесты
- AuthUseCase, PlayerUseCase - табличные тесты с mock-репозиториями

//Что дальше

- Спринт 2: столы Мастера, инвайты игрокам по @username, кнопки «Принять» и «Отклонить»
- Спринт 3: правила игры в JSON, персонажи, бестиарий, MinIO
- Спринт 4: боевой движок, Event Bus, polling событий, XP
- Спринт 5: Swagger, Dockerfile, CI, документация