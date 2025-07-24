export POSTGRES_DB_LOGIN='scada_admin'
export POSTGRES_DB_PASSWORD='scada_pass'
export POSTGRES_DB_HOST='scada_db'
export POSTGRES_DB_PORT=5432
export POSTGRES_DB_NAME='scada_db'
export POSTGRESQL_URL="postgres://${POSTGRES_DB_LOGIN}:${POSTGRES_DB_PASSWORD}@${POSTGRES_DB_HOST}:${POSTGRES_DB_PORT}/${POSTGRES_DB_NAME}?sslmode=disable"
export POSTGRES_MAX_CONNECTIONS=1
# export TESTS_BASE_DIR='/Users/macbook/Desktop/Data_Project/scada_consumer/tests'
export TESTS_BASE_DIR='/code/tests'

export RABBITMQ_URL=amqp://admin:pYEDBqnMLoWWE@rabbitmq:5672/
export RABBITMQ_CONSUMER_QUEUE='test'
export RABBITMQ_CONSUMER_TAG='tag'

export RABBITMQ_DAILY_QUEUE='DAILY'
export RABBITMQ_ASAP_QUEUE='ASAP'
export RABBITMQ_10MIN_QUEUE='10MIN'
