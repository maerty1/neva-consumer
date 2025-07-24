import json
import logging
from datetime import datetime

from app.models.measure_points import GetPointsDataRequest, GetPointsDataResponse, GetPointsDataResponseMeasurement

async def find_points_last_data(
    self, points: list[GetPointsDataRequest], timestamp: datetime | None
) -> list[GetPointsDataResponse]:
    if not points:
        return []

    elem_ids = [point.elem_id for point in points]

    latest_data = {}

    if timestamp is None:
        query = """
            SELECT
                mp.elem_id,
                mpd.values,
                mpd.datetime
            FROM
                measure_points mp
            JOIN
                measure_points_data_day mpd ON mp.id = mpd.measure_point_id
            WHERE
                mp.elem_id = ANY($1)
            ORDER BY
                mp.elem_id,
                mpd."datetime" DESC
        """

        async with self.pool.acquire() as connection:
            records = await connection.fetch(query, elem_ids)
            for record in records:
                elem_id = record["elem_id"]
                if elem_id not in latest_data:
                    latest_data[elem_id] = {
                        "values": record["values"],
                        "datetime": record["datetime"]
                    }
        # Для случая без timestamp, iscopied всегда False
        iscopied_map = {elem_id: False for elem_id in latest_data.keys()}

    else:
        query = """
    WITH ranked_data AS (
        SELECT
            mp.elem_id,
            mpd.values,
            mpd.datetime,
            ROW_NUMBER() OVER (
                PARTITION BY mp.elem_id
                ORDER BY mpd.datetime DESC
            ) AS rn
        FROM
            measure_points mp
        JOIN
            measure_points_data_day mpd ON mp.id = mpd.measure_point_id
        WHERE
            mp.elem_id = ANY($1)
            AND mpd.datetime <= $2
    )
    SELECT
        rd.elem_id,
        rd.values,
        rd.datetime
    FROM
        ranked_data rd
    WHERE
        rd.rn = 1;
        """

        async with self.pool.acquire() as connection:
            records = await connection.fetch(query, elem_ids, timestamp)
            for record in records:
                elem_id = record["elem_id"]
                if elem_id not in latest_data:
                    latest_data[elem_id] = {
                        "values": record["values"],
                        "datetime": record["datetime"]
                    }

        iscopied_map = {}
        for elem_id, data in latest_data.items():
            # Проверяем, совпадает ли datetime с запрошенным timestamp
            iscopied_map[elem_id] = data["datetime"] < timestamp

    responses = []

    for point in points:
        elem_id = point.elem_id
        measurements_request = point.measurements

        if elem_id not in latest_data:
            logging.warning(f"Данные для elem_id не найдены: {elem_id}")
            continue

        values_json = latest_data[elem_id]["values"]

        try:
            values = json.loads(values_json)
        except json.JSONDecodeError as e:
            logging.error(f"Ошибка декодирования JSON для elem_id {elem_id}: {e}")
            continue

        data_param_map = {
            item["dataParameter"]: item["value"] for item in values if "dataParameter" in item and "value" in item
        }

        measurements_response = {}

        for key, measurement in measurements_request.items():
            data_param_i = measurement.i
            data_param_o = measurement.o

            value_i = data_param_map.get(data_param_i)
            value_o = data_param_map.get(data_param_o)

            if value_i is None:
                logging.warning(f"DataParameter '{data_param_i}' не найден для elem_id {elem_id}")
                value_i = None

            if value_o is None:
                logging.warning(f"DataParameter '{data_param_o}' не найден для elem_id {elem_id}")
                value_o = None

            measurements_response[key] = GetPointsDataResponseMeasurement(i=value_i, o=value_o)

        # Определяем значение iscopied
        iscopied = iscopied_map.get(elem_id, False)

        response = GetPointsDataResponse(
            elem_id=elem_id,
            measurements=measurements_response,
            iscopied=iscopied  # Устанавливаем флаг
        )
        responses.append(response)

    return responses
