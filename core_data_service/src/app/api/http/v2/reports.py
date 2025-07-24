from fastapi import APIRouter, Depends, Header, Query

from app.models.v2.report import MergedBoilerRoomEngineerReportV2
from app.services.providers import provide_reports_service_stub
from app.services.reports import ReportsService

reports_router_v2 = APIRouter(tags=["Отчеты"], prefix="/core/api/v2")


@reports_router_v2.get(
    "/boiler_room_engineer_report/format/{format}",
    response_model=MergedBoilerRoomEngineerReportV2,
    name="Отчет инженера котельной второй версии",
    description="""Доступный формат: только json. \n
    В параметрах передается температура в формате float.""",
    status_code=200,
)
async def get_report(
    x_user_id: int = Header(..., ge=1),
    reports_service: ReportsService = Depends(provide_reports_service_stub),
    temperature: float = Query(...),
):
    return await reports_service.boiler_room_engineer_report(temperature)
