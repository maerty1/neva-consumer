import asyncpg
from app.repositories.measure_points import MeasurePointsRepository
from app.repositories.reports import ReportsRepository
from app.repositories.scada import ScadaRepository


def provide_reports_repository_stub():
    raise NotImplementedError


def provide_measure_points_repository_stub():
    raise NotImplementedError


def provide_scada_repository_stub():
    raise NotImplementedError


def provide_reports_repository(core_pool: asyncpg.pool.Pool) -> ReportsRepository:
    return ReportsRepository(core_pool)


def provide_measure_points_repository(core_pool: asyncpg.pool.Pool) -> MeasurePointsRepository:
    return MeasurePointsRepository(core_pool)


def provide_scada_repository(core_pool: asyncpg.pool.Pool) -> ScadaRepository:
    return ScadaRepository(core_pool)
