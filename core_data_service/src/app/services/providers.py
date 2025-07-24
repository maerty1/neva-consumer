from app.repositories.reports import ReportsRepository
from app.services.reports import ReportsService


def provide_reports_service_stub():
    raise NotImplementedError


def provide_reports_service(
    reports_repository: ReportsRepository,
) -> ReportsService:
    return ReportsService(reports_repository)
