import json
import logging

from pydantic import ValidationError

from app.models.measure_points import (
    GetMeasurePointsWithLastData,
    MeasurePointType,
    MeasurePointTypeParameters,
    PeriodType,
)


async def find_all_with_last_data(self, period_type: PeriodType):
    # Определяем имя таблицы и параметры даты в зависимости от периода
    if period_type == PeriodType.DAY:
        table = "measure_points_data_day"
        date_trunc_unit = "day"
        interval = "365 days"
    elif period_type == PeriodType.HOUR:
        table = "measure_points_data"
        date_trunc_unit = "hour"
        interval = "30 days"  # Пример: для часов можно брать за последний месяц
    else:
        raise ValueError(f"Неизвестный тип периода: {period_type}")

    query = f"""
            SELECT DISTINCT ON (mpd.measure_point_id)
                mpd.measure_point_id,
                mp.title,
                mp.address,
                mpd."datetime",
                mpd."values" as packet,
                mp.full_title,
                mpm.lat,
                mpm.lon,
                mpt.title as type_title,
                mptp.type_key,
                mptp.min_zoom,
                mptp.max_zoom
            FROM public.{table} mpd
            LEFT JOIN public.measure_points mp ON mp.id = mpd.measure_point_id
            LEFT JOIN public.measure_points_metadata mpm ON mpm.measure_point_id = mpd.measure_point_id
            LEFT JOIN public.measure_point_types mpt ON mpt.id = mp.type_id
            LEFT JOIN public.measure_point_type_parameters mptp ON mptp.type_id = mpt.id
            WHERE mpd."datetime" > (date_trunc('{date_trunc_unit}', NOW() AT TIME ZONE 'utc') - INTERVAL '{interval}')
            AND mp.account_id = 1
            ORDER BY mpd.measure_point_id, mpd."datetime" DESC;
        """

    async with self.pool.acquire() as connection:
        try:
            records = await connection.fetch(query)
        except Exception as e:
            logging.error(f"Ошибка выполнения запроса: {e}")
            return []

    measure_points = []
    for record in records:
        try:
            point_type = (
                MeasurePointType(
                    title=record["type_title"],
                    parameters=MeasurePointTypeParameters(
                        type_key=record["type_key"], minzoom=record["min_zoom"], maxzoom=record["max_zoom"]
                    ),
                )
                if record["type_title"]
                else None
            )

            measure_point = GetMeasurePointsWithLastData(
                packet_datetime=record["datetime"],
                title=record["title"],
                address=record["address"],
                measure_point_id=record["measure_point_id"],
                packet=json.loads(record["packet"]),
                full_title=record["full_title"],
                lat=record["lat"],
                lon=record["lon"],
                type=point_type,
            )
            measure_points.append(measure_point)
        except ValidationError as e:
            logging.error(f"Ошибка проверки записи {record}: {e}")
        except json.JSONDecodeError as e:
            logging.error(f"Ошибка декодирования JSON для записи {record}: {e}")

    return measure_points
