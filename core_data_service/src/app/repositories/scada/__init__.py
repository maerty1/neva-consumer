import asyncpg

from .v1.find_last_data import find_last_data
from .v1.find_measure_points import find_measure_points


class ScadaRepository:
    def __init__(self, pool: asyncpg.pool.Pool):
        self.pool = pool

    find_last_data = find_last_data
    find_measure_points = find_measure_points
