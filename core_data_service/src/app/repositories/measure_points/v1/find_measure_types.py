import logging

from pydantic import ValidationError

from app.models.measure_points import MeasureTypes


async def find_measure_types(self):
    query = """
            SELECT id, measure_type, title, unit
            FROM public.measure_types;
        """

    async with self.pool.acquire() as connection:
        try:
            records = await connection.fetch(query)
        except Exception as e:
            logging.error(f"Ошибка выполнения запроса: {e}")
            return []

    measure_types = []
    for record in records:
        try:
            # Создаем экземпляр модели MeasureTypes из данных записи
            measure_type = MeasureTypes(
                id=record["id"], measure_type=record["measure_type"], title=record["title"], unit=record["unit"]
            )
            measure_types.append(measure_type)
        except ValidationError as e:
            logging.error(f"Ошибка проверки записи {record}: {e}")

    return measure_types
