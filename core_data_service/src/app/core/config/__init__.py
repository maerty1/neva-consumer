import os

from app.core.config.postgres import PostgresConfig

reports_postgres_config = PostgresConfig(
    host=os.environ["CORE__POSTGRES_DB_HOST"],
    port=int(os.environ["CORE__POSTGRES_DB_PORT"]),
    username=os.environ["CORE__POSTGRES_DB_LOGIN"],
    password=os.environ["CORE__POSTGRES_DB_PASSWORD"],
    dbname=os.environ["CORE__POSTGRES_DB_NAME"],
    pool_size=int(os.environ["CORE__SQLALCHEMY_POOL_SIZE"]),
)
