# Описание работы сервиса

## Переменные окружения:
- `ELEM_ID` - ID элемента, с которым работает сервис (например, `3580` - котельная 16).
- `ROOT_CONNECTION_IDLE_TIME_SEC` - максимальное время бездействия в БД root.
- `ROOT_POSTGRES_MAX_CONNECTIONS` - максимальное количество подключений к БД root.
- `ROOT_POSTGRESQL_URL` - DSN для БД root.
- `WEATHER_BASE_URL` - ссылка на API для получения данных о погоде.
- `ZULU_BASE_URL` - ссылка на API Zulu.
- `ZULU_CONNECTION_IDLE_TIME_SEC` - максимальное время бездействия в БД Zulu.
- `ZULU_LAYER` - слой Zulu.
- `ZULU_POSTGRES_MAX_CONNECTIONS` - максимальное количество подключений к БД Zulu.
- `ZULU_POSTGRESQL_URL` - DSN для БД Zulu.
- `ZULU_TOKEN` - токен для работы с API Zulu.

## Описание работы:
1. Сервис получает данные о погоде через API, используя переменную `WEATHER_BASE_URL`.
2. Выполняет расчёты на сервере Zulu, взаимодействуя с API Zulu, указанным в `ZULU_BASE_URL`.
3. После расчётов данные сохраняются в таблицы:
    - `zulu.object_records`
    - `zulu.object_records_fromjson` (в виде JSON).
4. Сервис должен запускаться раз в час.
