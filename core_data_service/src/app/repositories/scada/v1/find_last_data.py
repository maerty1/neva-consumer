import json
import logging

from app.models.v1.scada import GetLastData, LastData
from pydantic import ValidationError


async def find_last_data(self, elem_id: int, limit: int, skip: int) -> GetLastData:
    query = """
            SELECT 
                grouped.name,
                json_agg(
                    json_build_object(
                        'timestamp', grouped."timestamp",
                        'value', grouped.value,
                        'measurement_type_id', grouped.measurement_type_id,
                        'measurement_name', grouped.measurement_name,
                        'measurement_type', grouped.measurement_type
                    )
                ) AS records
            FROM (
                SELECT 
                    smp.name,
                    sr."timestamp",
                    sr.value,
                    sr.measurement_type_id,
                    mt.title AS measurement_name,
                    mt."type" AS measurement_type
                FROM scada_measure_points smp
                JOIN scada_rawdata sr ON smp.id = sr.scada_measure_point_id
                JOIN measurement_types mt ON sr.measurement_type_id = mt.id
                WHERE smp.elem_id = $1
                ORDER BY sr."timestamp" DESC
                LIMIT $2 OFFSET $3
            ) grouped
            GROUP BY grouped.name
            ORDER BY MAX(grouped."timestamp") DESC;
        """

    async with self.pool.acquire() as connection:
        try:
            records = await connection.fetchrow(query, elem_id, limit, skip)
        except Exception as e:
            logging.error(f"Ошибка выполнения запроса: {e}")
            return {}

    try:
        last_data_name = records["name"]
    except ValidationError as e:
        logging.error(f"Ошибка проверки записи {record}: {e}")
        last_data_name = "-"

    last_data = []
    for record in json.loads(records["records"]):
        try:
            data = LastData(
                timestamp=record["timestamp"],
                value=record["value"],
                measurement_type_id=record["measurement_type_id"],
                measurement_name=record.get("measurement_name", "-"),
                measurement_type=record.get("measurement_type", "-"),
            )
            last_data.append(data)
        except ValidationError as e:
            logging.error(f"Ошибка проверки записи {record}: {e}")

    return GetLastData(name=last_data_name, packets=last_data)
