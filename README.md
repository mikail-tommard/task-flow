# task-flow
Сервис задач

## Описание
task-flow - сервис для управление задачами пользователя. Данные сохраняются в PostgreSQL. RabbitMQ - используется для асинхронных уведомлений

## Стек
1. Go(net/htto)
2. PostgreSQL
3. RabbitMQ
4. Docker Compose

```text
task-flow/
  cmd/
    api/
      main.go
    worker/
      main.go
  internal/
    domain/
      task.go
    usecase/
      service.go
      ports.go
    config/
      config.go
    adapters/
      http/
        json.go
    repository/
    broker/
  migrations/
  docker-compose.yml
  Makefile
  .env
```

## Как запускать сервис
Если используешь Docker-compose и Makefile:
```bash
docker-compose up -d
make migrate-up
make run
```

make run - поднимает API

Если Makefile нет и команды другие, используй прямые команды:
```bash
docker-compose up -d
migrate -path $MIGRATION_DIR -database $DB_DSN up
go run ./cmd/api
```

## Roadmap

- CRUD задач
- Фильтрация и сортировка задач
- Авторизация
- Уведомления через RabbitMQ

Данный проект делается чисто в учебных целях.