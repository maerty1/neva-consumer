import logging

from app.core import exceptions
from app.models.measure_points import GetPointsDataKeyvalueRequest, GetPointsDataKeyvalueResponse


async def find_points_data_keyvalue(self, points: list[GetPointsDataKeyvalueRequest]) -> GetPointsDataKeyvalueResponse:
    if not points:
        return GetPointsDataKeyvalueResponse()

    response = GetPointsDataKeyvalueResponse(measurements={})

    elem_ids = []
    keys = []
    for elem in points:
        response.measurements[elem.elem_id] = {}
        elem_ids.append(elem.elem_id)
        for key in elem.measurements:
            keys.append(key)

    query = """
    SELECT DISTINCT ON (sr.varname) 
        sr.varname,
        sr.value,
        smp.elem_id,
        sr."timestamp"
    FROM
        public.scada_rawdata sr
    JOIN
        public.scada_measure_points smp ON smp.id = sr.scada_measure_point_id
    WHERE
        smp.elem_id = ANY ($1)
        AND sr.varname = ANY ($2)
    ORDER BY
        sr.varname,
        sr."timestamp" DESC;
    """

    async with self.pool.acquire() as connection:
        try:
            records = await connection.fetch(query, elem_ids, keys)
            for record in records:
                response.measurements[record["elem_id"]][record["varname"]] = record["value"]

        except Exception as e:
            logging.error(f"Ошибка выполения запроса: {e}")
            raise exceptions.DatabaseError("Не удалось получить данные из базы данных.")
    return response
