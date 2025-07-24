from app.models.v1.scada import GetLastData
from app.repositories.providers import provide_scada_repository_stub
from app.repositories.scada import ScadaRepository
from fastapi import APIRouter, Depends, Header, Query

scada_router = APIRouter(tags=["SCADA"], prefix="/core/api/v1/scada")


@scada_router.get(
    "/last_data/{elem_id}",
    response_model=GetLastData,
    summary="Получение имени с последними данными.",
)
async def find_last_data(
    elem_id: int,
    x_user_id: int = Header(..., ge=1),
    limit: int = Query(default=100),
    skip: int = Query(default=0),
    scada_repository: ScadaRepository = Depends(provide_scada_repository_stub),
):
    return await scada_repository.find_last_data(elem_id, limit, skip)


@scada_router.get(
    "/measure_points",
    response_model=dict[int, str],
    summary="Получение имени с id.",
)
async def find_measure_points(
    x_user_id: int = Header(..., ge=1), scada_repository: ScadaRepository = Depends(provide_scada_repository_stub)
):
    return await scada_repository.find_measure_points()
