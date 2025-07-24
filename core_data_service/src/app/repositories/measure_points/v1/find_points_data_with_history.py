import json
import logging
from datetime import datetime

from app.core import exceptions
from app.models.measure_points import GetPointsDataRequest, GetPointsDataResponseWithHistory


async def find_points_data_with_history(
    self, points: list[GetPointsDataRequest], last_n: int | None, timestamp: datetime | None
) -> list[GetPointsDataResponseWithHistory]:
    if not points:
        return []

    elem_ids = [point.elem_id for point in points]

    if timestamp:
        # Запрос для получения данных по конкретному timestamp
        query = """
        SELECT
            mp.elem_id,
            mpd.values,
            mpd."datetime"
        FROM
            measure_points mp
        JOIN
            measure_points_data_day mpd ON mp.id = mpd.measure_point_id
        WHERE
            mp.elem_id = ANY($1)
            AND mpd."datetime" = $2
        ORDER BY
            mp.elem_id, mpd."datetime" DESC
        """
        query_params = (elem_ids, timestamp)
    else:
        # Запрос для получения последних n записей по каждому elem_id
        query = """
        SELECT
            sub.elem_id,
            sub.values,
            sub."datetime"
        FROM (
            SELECT
                mp.elem_id,
                mpd.values,
                mpd."datetime",
                ROW_NUMBER() OVER (PARTITION BY mp.elem_id ORDER BY mpd."datetime" DESC) as rn
            FROM
                measure_points mp
            JOIN
                measure_points_data_day mpd ON mp.id = mpd.measure_point_id
            WHERE
                mp.elem_id = ANY($1)
        ) sub
        WHERE
            sub.rn <= $2
        ORDER BY
            sub.elem_id, sub."datetime" DESC
        """
        query_params = (elem_ids, last_n)

    response_data: dict[str, dict[str, dict[str, float]]] = {}

    async with self.pool.acquire() as connection:
        try:
            records = await connection.fetch(query, *query_params)
        except Exception as e:
            logging.error(f"Ошибка при выполнении запроса: {e}")
            raise exceptions.DatabaseError("Не удалось получить данные из базы данных.")

    points_dict = {point.elem_id: point.measurements for point in points}

    for record in records:
        elem_id = str(record["elem_id"])
        datetime_obj = record["datetime"]
        datetime_str = datetime_obj.strftime("%Y-%m-%d")
        values_raw = record["values"]

        if isinstance(values_raw, str):
            try:
                values = json.loads(values_raw)
            except json.JSONDecodeError as e:
                logging.error(f"Ошибка при парсинге JSON для elem_id {elem_id} и datetime {datetime_str}: {e}")
                continue
        elif isinstance(values_raw, list):
            values = values_raw
        else:
            logging.error(f"Неподдерживаемый тип данных для values: {type(values_raw)}")
            continue

        if elem_id not in response_data:
            response_data[elem_id] = {}

        if datetime_str not in response_data[elem_id]:
            response_data[elem_id][datetime_str] = {}

        measurements = points_dict.get(int(elem_id), {})

        for measurement_key, measurement in measurements.items():
            group_id = measurement_key

            if group_id not in response_data[elem_id][datetime_str]:
                response_data[elem_id][datetime_str][group_id] = {}

            i_param = measurement.i
            o_param = measurement.o

            # Ищем соответствующие значения в списке values
            i_value = next((item["value"] for item in values if item.get("dataParameter") == i_param), None)
            o_value = next((item["value"] for item in values if item.get("dataParameter") == o_param), None)

            if i_value is not None:
                response_data[elem_id][datetime_str][group_id]["i"] = round(i_value, 2)

            if o_value is not None:
                response_data[elem_id][datetime_str][group_id]["o"] = round(o_value, 2)

    return response_data
