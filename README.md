# postgres-sqlpad-setup
Docker compose окружение для развёртывания PostgreSQL и SQLPad с пользователем только на чтение.

## Копируем себе
```bash
git clone https://github.com/nskondratev/postgres-sqlpad-setup.git
```

## Требования
* [Docker](https://docs.docker.com/engine/install/)
* [Docker Compose](https://docs.docker.com/compose/install/)

## Структура проекта
* `data` - директория, где лежат данные PostgreSQL и базы данных SQLPad. Нужно для того, чтобы сохранять состояние между перезапусками.
* `docker-entrypoint-initdb.d` - директория со скриптами, которые выполняются при пустом старте PostgreSQL. **Если в директории `data/postgres` есть какие-нибудь файлы, то скрипты из данной директории выполняться не будут.**
* `pg_init_data` - директория, куда нужно положить CSV-файлы с данными для PostgreSQL.
* `reset.sh` - вспомогательный скрипт для остановки и полной очистки данных PostgreSQL и SQLPad (полезно при отладке скриптов из `docker-entrypoint-initdb.d`)

## Как этим пользоваться
### Запуск
```bash
docker-compose up -d
```
После этого можно войти в веб-интерфейс SQLPad по адресу http://localhost:3000.
Логин (по умолчанию): `admin@sqlpad.com`
Пароль (по умолчанию): `admin`

### Как подключиться к PostgreSQL через терминал
```bash
docker exec -it postgres-sqlpad-setup_postgres_1 psql -U sqlpad sqlpad
```
Откроется сессия в psql. У пользователя полные права (создание/изменение таблиц, запись данных и т.п.)

### Как отладить скрипты создания таблиц и импорта данных
* Запустить PostgreSQL не в фоновом режиме: `docker-compose up postgres`
* Смотреть на логи:
  * Если будет ошибка, то контейнер завершится, сообщение об ошибке будет последним в логах.
  * Если ошибки не будет, то можно увидеть в логах, что скрипты инициализации выполнились и PostgreSQL перезапустился:
```
postgres_1  | /usr/local/bin/docker-entrypoint.sh: running /docker-entrypoint-initdb.d/02_create_tables.sql
postgres_1  | CREATE TABLE
postgres_1  | CREATE TABLE
postgres_1  | CREATE TABLE
postgres_1  | 
postgres_1  | 
postgres_1  | /usr/local/bin/docker-entrypoint.sh: running /docker-entrypoint-initdb.d/03_fill_tables.sql
```
* Чтобы заново отладить эти скрипты, нужно остановить контейнер, если он не завершился: Ctrl + C
* Выполнить очистку: `./reset.sh`
* Вернуться к первому шагу.
