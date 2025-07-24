from fastapi import APIRouter, Depends, Header

from app.repositories.measure_points import MeasurePointsRepository
from app.repositories.providers import provide_measure_points_repository_stub

weather_router = APIRouter(tags=["Погода"], prefix="/core/api/v1")


@weather_router.get(
    "/weather/with_forecast",
    responses={
        200: {
            "content": {
                "application/json": {
                    "example": {
                        "today": {
                            "temp_avg": -1.36,
                            "temp": -4.3,
                            "pressure_avg": 770.8,
                            "pressure": 765.0,
                            "humidity": 27.3,
                            "humidity_avg": 74,
                            "wind_speed": 0.0,
                            "wind_speed_avg": 2,
                            "date": "2024-12-04 00:00:00",
                        },
                        "tomorrow": {
                            "temp_avg": -3.22,
                            "temp": 2,
                            "pressure_avg": 777.0,
                            "pressure": 3,
                            "humidity": 4,
                            "humidity_avg": 73,
                            "wind_speed": 2,
                            "wind_speed_avg": 2,
                            "date": "2024-12-05 00:00:00",
                        },
                        "yesterday": {
                            "temp_avg": 3.28,
                            "temp": 0.3,
                            "pressure_avg": 756.7,
                            "pressure": 756.0,
                            "humidity": 33.2,
                            "humidity_avg": 91,
                            "wind_speed": 2.7,
                            "wind_speed_avg": 5,
                            "date": "2024-12-03 00:00:00",
                        },
                    }
                }
            }
        }
    },
)
async def get_weather_forecast(
    x_user_id: int = Header(..., ge=1),
    measure_points_repository: MeasurePointsRepository = Depends(provide_measure_points_repository_stub),
):
    return await measure_points_repository.find_weather_forecast()
