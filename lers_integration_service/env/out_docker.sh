export POSTGRES_DB_LOGIN='admin'
export POSTGRES_DB_PASSWORD='ZpVurRuj5AHX4C4W@uKt'
export POSTGRES_DB_HOST='localhost'
export POSTGRES_DB_PORT=5432
export POSTGRES_DB_NAME='postgres_db'
export POSTGRESQL_URL="postgres://${POSTGRES_DB_LOGIN}:${POSTGRES_DB_PASSWORD}@${POSTGRES_DB_HOST}:${POSTGRES_DB_PORT}/${POSTGRES_DB_NAME}?sslmode=disable"
export POSTGRES_MAX_CONNECTIONS=1