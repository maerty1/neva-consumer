from app.repositories.measure_points import MeasurePointsRepository
from app.repositories.measure_points.v1.get_system_status import SystemStatusResponse
from app.repositories.providers import provide_measure_points_repository_stub
from fastapi import APIRouter, Depends, Header, Query

status_router = APIRouter(tags=["Состояние"], prefix="/core/api/v1")


@status_router.get("/status/current", response_model=SystemStatusResponse)
async def get_system_status(
    x_user_id: int = Header(..., ge=1),
    fresh_interval_days: int = Query(default=1),
    success_interval_days: int = Query(default=7),
    measure_points_repository: MeasurePointsRepository = Depends(provide_measure_points_repository_stub),
):
    return await measure_points_repository.get_system_status(
        user_id=x_user_id, fresh_interval_days=fresh_interval_days, success_interval_days=success_interval_days
    )
