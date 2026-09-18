# todo

CLI и HTTP API для списка задач. Данные хранятся в JSON-файле (по умолчанию `todos.json`).

## Структура репо

```
.
├── cmd/
│   ├── todo/          # CLI
│   └── server/        # HTTP-сервер
├── internal/
│   ├── todo/          # домен: Store, Load/Save
│   └── api/           # HTTP handlers
├── go.mod
└── README.md
```

## Запуск

Требуется Go 1.22+ (в проекте используются маршруты вида `GET /tasks/{id}`).

### CLI

```bash
go run ./cmd/todo add "купить молоко"
go run ./cmd/todo list
go run ./cmd/todo done 1
go run ./cmd/todo delete 1
```

Флаг `-file` — путь к файлу (по умолчанию `todos.json`):

```bash
go run ./cmd/todo -file my.json list
```

### Server

```bash
go run ./cmd/server
```

По умолчанию слушает `:8080`, файл — `todos.json`:

```bash
go run ./cmd/server -addr :8080 -file todos.json
```

## HTTP API: `/tasks`

| Метод  | Путь              | Описание              | Успех |
|--------|-------------------|-----------------------|-------|
| GET    | `/tasks`          | Список задач          | 200   |
| POST   | `/tasks`          | Создать задачу        | 201   |
| GET    | `/tasks/{id}`     | Получить задачу       | 200   |
| POST   | `/tasks/{id}/done`| Отметить выполненной  | 200   |
| DELETE | `/tasks/{id}`     | Удалить задачу        | 204   |

Тело `POST /tasks`:

```json
{"text": "купить молоко"}
```

Ответ задачи:

```json
{"id": 1, "text": "купить молоко", "done": false}
```

Ошибки: `400` (невалидный JSON / пустой текст / неверный id), `404` (задача не найдена).

## Примеры curl

```bash
curl -s http://localhost:8080/tasks

curl -s -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"text":"купить молоко"}'

curl -s http://localhost:8080/tasks/1

curl -s -X POST http://localhost:8080/tasks/1/done

curl -s -X DELETE http://localhost:8080/tasks/1 -w "\n%{http_code}\n"
```

## Примеры CLI

```bash
go run ./cmd/todo add "купить молоко"
# 1: купить молоко, false

go run ./cmd/todo list
# 1: купить молоко false

go run ./cmd/todo done 1
go run ./cmd/todo list
# 1: купить молоко true

go run ./cmd/todo delete 1
```

## Важно

CLI и server используют **один дефолтный файл** `todos.json`. Не пишите в него одновременно из CLI и server — возможна порча данных или потеря изменений. Останавливайте server перед CLI-мутациями (или используйте разные `-file`).
