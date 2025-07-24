from datetime import datetime
import logging
from fastapi import FastAPI

import logging
import asyncio
from datetime import datetime
from apscheduler.schedulers.asyncio import AsyncIOScheduler
from fastapi import FastAPI
from cache import cache
from aiocache import caches

import asyncpg.pool
from app.api.http.v1 import routes as v1_routes
from app.api.http.v2.reports import reports_router_v2
from app.core.exc_handling import configure_exceptions_handlers
from app.repositories.providers import (
    provide_measure_points_repository,
    provide_measure_points_repository_stub,
    provide_reports_repository,
    provide_reports_repository_stub,
    provide_scada_repository,
    provide_scada_repository_stub,
)
from app.services.providers import provide_reports_service, provide_reports_service_stub
from fastapi import FastAPI

from app.repositories.reports import ReportsRepository


class Application:
    def __init__(self, app: FastAPI, pool: asyncpg.pool.Pool):
        self.app = app
        self.core_pool = pool
        self.logger = logging.getLogger(__name__)
        self.scheduler = AsyncIOScheduler()

    def _create_repositories(self):
        self.reports_repository: ReportsRepository = lambda: provide_reports_repository(self.core_pool)
        self.measure_points_repository = lambda: provide_measure_points_repository(self.core_pool)
        self.scada_repository = lambda: provide_scada_repository(self.core_pool)

    def _create_services(self):
        self.reports_service = lambda: provide_reports_service(self.reports_repository())

    def _override_dependencies(self):
        self.app.dependency_overrides[provide_reports_repository_stub] = self.reports_repository
        self.app.dependency_overrides[provide_measure_points_repository_stub] = self.measure_points_repository
        self.app.dependency_overrides[provide_scada_repository_stub] = self.scada_repository
        self.app.dependency_overrides[provide_reports_service_stub] = self.reports_service

    def _add_routes(self):
        for route in v1_routes:
            self.app.include_router(route)

        self.app.include_router(reports_router_v2)

    def _configure_logging(self):
        logging.basicConfig(
            # level=int(os.environ["LOGGING_LEVEL"]),
            format="%(levelname)s %(asctime)s %(filename)s:%(lineno)d %(message)s",
        )

    def _configure_exception_handlers(self):
        configure_exceptions_handlers(self.app)

    async def _update_cache(self):
        current_year = datetime.now().year
        start_year = current_year - 1

        for year in range(start_year, current_year + 1):  # 2 года: прошлый и текущий
            for month in range(1, 13):
                try:
                    self.logger.info(f"🔄 Обновление кеша для {year}-{month}")
                    data = await self.reports_repository().boiler_room_engineer_report(year, month)

                    cache_key = f"boiler_room_engineer_report:{year}-{month}"
                    cached_data = await cache.get(cache_key)

                    if cached_data:
                        self.logger.info(f"Данные успешно записаны в кеш для {year}-{month}")
                    else:
                        self.logger.error(f"Кеш не записался для {year}-{month}")

                except Exception as e:
                    self.logger.error(f"⚠️ Ошибка при обновлении кеша {year}-{month}: {e}")

    async def _start_scheduler(self):
        self.scheduler.add_job(self._update_cache, "interval", hours=12)
        self.scheduler.start()

    async def startup(self):
        asyncio.create_task(self._update_cache())
        await self._start_scheduler()

    def build(self):
        self._create_repositories()
        self._create_services()
        self._override_dependencies()
        self._add_routes()
        self._configure_logging()
        self._configure_exception_handlers()
        return self
