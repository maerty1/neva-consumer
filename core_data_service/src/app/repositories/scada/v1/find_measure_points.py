import json
import logging

from app.models.v1.scada import LastData
from pydantic import ValidationError


async def find_measure_points(self) -> dict[int, str]:
    query = """
        SELECT 
            smp.elem_id, smp.name
        FROM 
            scada_measure_points smp
        WHERE
            elem_id IS NOT NULL;
        """

    async with self.pool.acquire() as connection:
        try:
            records = await connection.fetch(query)
        except Exception as e:
            logging.error(f"Ошибка выполнения запроса: {e}")
            return {}

    last_data = {}
    for elem_id, name in records:
        try:
            last_data[elem_id] = name
        except ValidationError as e:
            logging.error(f"Ошибка проверки записи {records}: {e}")
    return last_data
