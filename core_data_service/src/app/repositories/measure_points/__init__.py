import asyncpg

from .v1.find_all_with_last_data import find_all_with_last_data
from .v1.find_by_id_with_data import find_by_id_with_data
from .v1.find_measure_types import find_measure_types
from .v1.find_points_data_keyvalue import find_points_data_keyvalue
from .v1.find_points_data_with_history import find_points_data_with_history
from .v1.find_points_element_ids import find_points_element_ids
from .v1.find_points_last_data import find_points_last_data
from .v1.find_weather_forecast import find_weather_forecast
from .v1.get_system_status import get_system_status


class MeasurePointsRepository:
    def __init__(
        self,
        pool: asyncpg.pool.Pool,
    ):
        self.pool = pool

    find_by_id_with_data = find_by_id_with_data
    find_all_with_last_data = find_all_with_last_data
    find_measure_types = find_measure_types
    find_points_last_data = find_points_last_data
    find_points_data_with_history = find_points_data_with_history
    find_points_data_keyvalue = find_points_data_keyvalue
    find_points_element_ids = find_points_element_ids
    find_weather_forecast = find_weather_forecast

    get_system_status = get_system_status
