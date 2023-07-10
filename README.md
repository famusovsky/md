# MD pastebin

Сервис, позволяющий сохранять и делиться заметками в формате [`markdown`](https://www.markdownguide.org/getting-started/)

### Для запуска проекта необходимо:

1. Создать базу данных PostgreSQL
2. Запустить проект:

```bash
go run ./cmd/web -addr=<TCP адрес проекта> -dsn=<строка подключения к бд>

# addr по умолчанию = ":8080"
# dsn по умолчанию = "port=5432 user=postgres password=qwerty dbname=MD sslmode=disable"
```
