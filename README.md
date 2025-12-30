# task-flow
Сервис задач

## Описание
task-flow - сервис для управдение задачами пользователя. Данные сохраняются в PostgreSQL. RabbitMQ - используется для асинхронных уведомлений

## Стэк
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
```