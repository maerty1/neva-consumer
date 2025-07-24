import asyncpg

from .v1.boiler_room_engineer_report import boiler_room_engineer_report
from .v2.boiler_room_engineer_report_v2 import boiler_room_engineer_report_v2

import logging


class ReportsRepository:
    def __init__(
        self,
        pool: asyncpg.pool.Pool,
    ):
        self.pool = pool
        self.logger = logging.getLogger(self.__class__.__name__)

    boiler_room_engineer_report = boiler_room_engineer_report
    boiler_room_engineer_report_v2 = boiler_room_engineer_report_v2
