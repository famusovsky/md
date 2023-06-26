# MD pastebin

### Для запуска проекта необходимо:

1. Создать базу данных PostgreSQL
2. Создать таблицы [(см. ниже)](#postgresql-query-для-создания-таблиц)
3. Запустить проект

### PostgreSQL Query для создания таблиц
```SQL
CREATE TABLE notes (
    id SERIAL NOT NULL PRIMARY KEY, 
    content TEXT NOT NULL,
    created TIMESTAMP NOT NULL,
    expires TIMESTAMP NOT NULL
);
```
