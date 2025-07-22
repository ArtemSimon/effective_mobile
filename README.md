# REST-сервис для агрегации данных об онлайн-подписках пользователей.

# Структура проекта

/effective_mobile
├── app                    # Точка входа (main.go)
├── /internal
│   ├── /api               # HTTP-обработчики (роуты, контроллеры)
│   ├── /service           # Бизнес-логика (use cases)
│   ├── /objects           # Структуры 
│   ├── /repository        # Работа с БД (PostgreSQL)
│   └── /config            # Конфигурация 
├── /migrations            # SQL-миграции 
├── /pkg/logger_module     # Логгер
├── /docs                  # Swagger-документация
├── docker-compose.yml     
├── Dockerfile
├── .env                   # Добавил чтобы проще было запустить, а так нельзя из за безопасности
├── .gitignore
├── app.log 
├── go.sum
├── go.mod
└── README.md 

## ⚙️ Технологии
- Go 1.23
- PostgreSQL
- Docker
- Swagger (документация API)
- GORM (ORM)

## 🚀 Запуск

# Предварительные требования

    Установленный Docker 

    Установленный Docker Compose 

    Файлы окружения (если не созданы, см. раздел "Настройка окружения")

# Настройка окружения

Я остовил .env файл в репозитории, но так делать нельзя из за соображений безопасности,обычно можно использовать docker secrets или же Vault,так же можете создать после клонирования репозитория собственный .env файл и подцепить к проекту,так как тут в коде читаются переменные окружения переданные из docker compose.

Вот формат env файла: 

```
POSTGRES_USER=ваш_пользователь
POSTGRES_PASSWORD=ваш_пароль
POSTGRES_DB=ваша_база_данных

DB_HOST=хост по которому идет connect (в docker это имя контейнера)
DB_PORT=порт
DB_USER=имя пользователя
DB_PASSWORD=пароль от бд
DB_NAME=название бд
DB_SSLMODE=disable  # для разработки

# HTTP-сервер
HTTP_PORT=порт для приложения
```

# Клонируйте репозиторий

git clone https://github.com/ArtemSimon/effective_mobile.git

# Перейдите в директорию проекта

```
cd effective_mobile
```
# Выполните команду:

Тут можно использовать готовый образ который я запулил в открытый доступ, он указан в docker compose файле  или же собирать как указано в docker compose
Если хотите использовать готовый образ:

1. Скачиваем образ к себе
```
docker pull artemsim/effective_mobile_hub:latest
```
2.Запускаем с готовым образом docker compose  
```
docker-compose up 
```

Если просто хотите собрать образ сами то:

```
docker-compose up --build
```
Соберет образ на один запуск

Для остановки всех сервисов выполните:

```
docker-compose down
```

# Так же у меня написаны тесты к приложению и чтобы их запустить нужно:
1. Провалится в контейнер 

```
docker exec -it имя_контейнера(rest_service) sh
```
2. Запустить бинарник с тестами

```
/app/effective-mobile-test -test.v 
``` 

Будут видны логи плюсом