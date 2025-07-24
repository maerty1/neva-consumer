export USER__POSTGRES_DB_LOGIN='scada_admin'
export USER__POSTGRES_DB_PASSWORD='scada_pass'
export USER__POSTGRES_DB_HOST='localhost'
export USER__POSTGRES_DB_PORT=5434
export USER__POSTGRES_DB_NAME='user_db'
export USER__MAX_CONNECTIONS=1
export USER__CONNECTION_IDLE_TIME_SEC=10

export USER__POSTGRESQL_URL="postgres://${USER__POSTGRES_DB_LOGIN}:${USER__POSTGRES_DB_PASSWORD}@${USER__POSTGRES_DB_HOST}:${USER__POSTGRES_DB_PORT}/${USER__POSTGRES_DB_NAME}?sslmode=disable"

export USER__HTTP_HOST='0.0.0.0'
export USER__HTTP_PORT=8001
# Заменить на GIN_MODE=release во время деплоя
export GIN_MODE=release


