import asyncpg
from sqlalchemy.engine import URL
from sqlalchemy.ext.asyncio import AsyncEngine, AsyncSession, create_async_engine
from sqlalchemy.orm import sessionmaker

from app.core.config import reports_postgres_config


# TEST
def construct_url(username: str, password: str, host: str, port: int, dbname: str) -> URL:
    return URL.create(
        "postgresql+asyncpg",
        username=username,
        password=password,
        host=host,
        port=port,
        database=dbname,
    )


def provide_engine(url: URL | str, echo: bool, pool_size: int) -> AsyncEngine:
    return create_async_engine(url=url, echo=echo, pool_size=pool_size)


def provide_session(engine: AsyncEngine) -> sessionmaker:
    return sessionmaker(engine=engine, class_=AsyncSession, expire_on_commit=False)


async def provide_pool() -> asyncpg.pool.Pool:
    return await asyncpg.create_pool(
        user=reports_postgres_config.username,
        password=reports_postgres_config.password,
        database=reports_postgres_config.dbname,
        host=reports_postgres_config.host,
        port=reports_postgres_config.port,
        # max_size=reports_postgres_config.pool_size,
    )


def configure_engines() -> AsyncEngine:
    reports_connection_uri = construct_url(
        reports_postgres_config.username,
        reports_postgres_config.password,
        reports_postgres_config.host,
        reports_postgres_config.port,
        reports_postgres_config.dbname,
    )
    reports_engine = provide_engine(
        sim_card_connection_uri,
        echo=False,
        pool_size=reports_postgres_config.pool_size,
    )
    return reports_engine
