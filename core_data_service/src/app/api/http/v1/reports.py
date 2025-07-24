from fastapi import APIRouter, Depends, Header, HTTPException, Query

from app.file_converters.base_strategy import ReportFormatStrategy
from app.file_converters.boiler_room_engineer import create_boiler_room_engineer_report_filename
from app.file_converters.boiler_room_engineer.json import JSONStrategy
from app.file_converters.boiler_room_engineer.xlsx import XLSXStrategy
from app.models.report import BoilerRoomEngineerReport
from app.repositories.providers import provide_reports_repository_stub
from app.repositories.reports import ReportsRepository

reports_router = APIRouter(tags=["Отчеты"], prefix="/core/api/v1")


def provide_boiler_room_engineer_report_strategy(format: str) -> ReportFormatStrategy:
    if format.lower() == "xlsx":
        return XLSXStrategy()
    elif format.lower() == "json":
        return JSONStrategy()
    else:
        raise HTTPException(status_code=400, detail="Поддерживаемые форматы: xlsx, json")


@reports_router.get(
    "/boiler_room_engineer_report/format/{format}",
    response_model=list[BoilerRoomEngineerReport],
    name="Отчет инженера котельной",
    description="Доступные форматы: `json`, `xslx`",
    status_code=200,
)
async def get_report(
    format: str,
    x_user_id: int = Header(..., ge=1),
    year: int = Query(..., ge=1900, le=2100, description="Год отчета"),
    month: int = Query(..., ge=1, le=12, description="Месяц отчета"),
    reports_repository: ReportsRepository = Depends(provide_reports_repository_stub),
    strategy: ReportFormatStrategy = Depends(provide_boiler_room_engineer_report_strategy),
):
    data = await reports_repository.boiler_room_engineer_report(year, month)
    filename = create_boiler_room_engineer_report_filename(year=year, month=month, format=format)
    return strategy.convert_and_response(data, filename=filename, year=year)
