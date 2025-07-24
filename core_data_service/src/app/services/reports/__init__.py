from app.repositories.reports import ReportsRepository

from .v1.boiler_room_engineer_report import boiler_room_engineer_report


class ReportsService:
    def __init__(self, reports_repository: ReportsRepository):
        self.reports_repository = reports_repository

    boiler_room_engineer_report = boiler_room_engineer_report
