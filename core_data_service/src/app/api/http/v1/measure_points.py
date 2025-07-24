from datetime import datetime

from fastapi import APIRouter, Body, Depends, Header, HTTPException, Path, Query

from app.models.measure_points import (
    GetMeasurePointsWithLastData,
    GetMeasurePointWithLastData,
    GetPointsDataKeyvalueRequest,
    GetPointsDataKeyvalueResponse,
    GetPointsDataRequest,
    GetPointsDataResponse,
    MeasureTypes,
    PeriodType,
)
from app.repositories.measure_points import MeasurePointsRepository
from app.repositories.providers import provide_measure_points_repository_stub

measure_points_router = APIRouter(tags=["Точки измерения"], prefix="/core/api/v1")


@measure_points_router.get(
    "/measure_points/with_last_data",
    response_model=list[GetMeasurePointsWithLastData],
    summary="Получение списка точек измерения с последними данными",
)
async def find_measure_points(
    x_user_id: int = Header(..., ge=1),
    period_type: PeriodType = Query(default=PeriodType.DAY),
    measure_points_repository: MeasurePointsRepository = Depends(provide_measure_points_repository_stub),
):
    return await measure_points_repository.find_all_with_last_data(period_type)


@measure_points_router.get(
    "/measure_points/{measure_point_id}/with_last_data",
    summary="Получение точки измерения по id с конфигурируемым количеством последних данных",
    response_model=GetMeasurePointWithLastData,
)
async def find_measure_point_by_id(
    measure_point_id: int = Path(..., ge=1),
    packets_limit: int = Query(ge=1, le=100, default=12, description="Максимальное количество пакетов в ответе"),
    x_user_id: int = Header(..., ge=1),
    period_type: PeriodType = Query(default=PeriodType.DAY),
    measure_points_repository: MeasurePointsRepository = Depends(provide_measure_points_repository_stub),
):
    return await measure_points_repository.find_by_id_with_data(measure_point_id, period_type, packets_limit)


@measure_points_router.get(
    "/measure_points/measure_types",
    summary="Получение маппинга для сопопставления переменных LERS и версии на русском языке",
    response_model=list[MeasureTypes],
)
async def find_measure_type(
    x_user_id: int = Header(..., ge=1),
    measure_points_repository: MeasurePointsRepository = Depends(provide_measure_points_repository_stub),
):
    return await measure_points_repository.find_measure_types()


@measure_points_router.post(
    "/points/last_data",
    response_model=list[GetPointsDataResponse],
)
async def get_points_data(
    x_user_id: int = Header(..., ge=1),
    req: list[GetPointsDataRequest] = Body(),
    timestamp: datetime = Query(default=None),
    measure_points_repository: MeasurePointsRepository = Depends(provide_measure_points_repository_stub),
):
    return await measure_points_repository.find_points_last_data(req, timestamp)


@measure_points_router.post(
    "/points/data",
    # response_model=list[GetPointsDataResponseWithHistory]
)
async def get_points_data(
    x_user_id: int = Header(..., ge=1),
    req: list[GetPointsDataRequest] = Body(),
    n_days: int = Query(None),
    timestamp: datetime = Query(None),
    measure_points_repository: MeasurePointsRepository = Depends(provide_measure_points_repository_stub),
):
    if n_days is not None and timestamp is not None:
        raise HTTPException(
            status_code=400,
            detail="Должно быть указано только одно из значений «n_days» или «timestamp», а не оба одновременно.",
        )
    if n_days is None and timestamp is None:
        raise HTTPException(
            status_code=400,
            detail="Необходимо указать одно из значений «n_days» или «timestamp».",
        )
    return await measure_points_repository.find_points_data_with_history(req, n_days, timestamp)


@measure_points_router.post("/points/data/keyvalue", response_model=GetPointsDataKeyvalueResponse)
async def get_points_data(
    x_user_id: int = Header(..., ge=1),
    req: list[GetPointsDataKeyvalueRequest] = Body(),
    measure_points_repository: MeasurePointsRepository = Depends(provide_measure_points_repository_stub),
):
    return await measure_points_repository.find_points_data_keyvalue(req)


@measure_points_router.get("/points/element_ids", response_model=list)
async def get_element_ids(
    x_user_id: int = Header(..., ge=1),
    measure_points_repository: MeasurePointsRepository = Depends(provide_measure_points_repository_stub),
):
    return await measure_points_repository.find_points_element_ids()
