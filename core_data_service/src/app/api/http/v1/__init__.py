from app.api.http.v1.measure_points import measure_points_router
from app.api.http.v1.reports import reports_router
from app.api.http.v1.scada import scada_router
from app.api.http.v1.status import status_router
from app.api.http.v1.weather import weather_router

routes = [measure_points_router, reports_router, scada_router, weather_router, status_router]
