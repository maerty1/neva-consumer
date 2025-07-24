import logging


async def find_points_element_ids(self):
    query = """
    SELECT elem_id FROM public.measure_points WHERE elem_id IS NOT NULL
    UNION
    SELECT elem_id FROM public.scada_measure_points WHERE elem_id IS NOT NULL;
    """

    async with self.pool.acquire() as connection:
        try:
            records = await connection.fetch(query)
        except Exception as e:
            logging.error(f"Ошибка выполнения запроса: {e}")
            return []

    element_ids = [record["elem_id"] for record in records]

    return element_ids
