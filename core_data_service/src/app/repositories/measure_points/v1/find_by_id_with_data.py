import json

from pydantic import ValidationError

from app.core.decorators import asyncpg_exc_handler
from app.models.measure_points import (
    GetMeasurePointWithLastData,
    MeasurePointType,
    MeasurePointTypeParameters,
    PacketData,
    PeriodType,
)


@asyncpg_exc_handler
async def find_by_id_with_data(self, measure_point_id: int, period_type: PeriodType, packets_limit: int = 12):
    if period_type == PeriodType.DAY:
        table = "measure_points_data_day"
        date_trunc_unit = "day"
        interval = "365 days"
    elif period_type == PeriodType.HOUR:
        table = "measure_points_data"
        date_trunc_unit = "hour"
        interval = "30 days"
    else:
        raise ValueError(f"Неизвестный тип периода: {period_type}")

    query = f"""
            SELECT
                mpd.measure_point_id,
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
            WHERE mpd.measure_point_id = $1
            AND mpd."datetime" > (date_trunc('{date_trunc_unit}', NOW() AT TIME ZONE 'utc') - INTERVAL '{interval}')
            AND mp.account_id = 1
            ORDER BY mpd."datetime" DESC
            LIMIT {packets_limit};
        """

    async with self.pool.acquire() as connection:
        try:
            records = await connection.fetch(query, measure_point_id)
        except Exception as e:
            logging.error(f"Ошибка выполнения запроса: {e}")
            return []

    if records == []:
        raise exceptions.NotFoundException(detail=f"Точка измерения с id {measure_point_id} не найдена")

    try:
        first_record = records[0]
        point_type = (
            MeasurePointType(
                title=first_record["type_title"],
                parameters=MeasurePointTypeParameters(
                    type_key=first_record["type_key"],
                    minzoom=first_record["min_zoom"],
                    maxzoom=first_record["max_zoom"],
                ),
            )
            if first_record["type_title"]
            else None
        )
        measure_point_data = GetMeasurePointWithLastData(
            measure_point_id=first_record["measure_point_id"],
            lat=first_record.get("lat"),
            lon=first_record.get("lon"),
            full_title=first_record["full_title"],
            packets=[],
            type=point_type,
        )
    except ValidationError as e:
        logging.error(f"Ошибка инициализации модели для записи {first_record}: {e}")
        return None

    for record in records:
        try:
            packet_data = PacketData(packet_datetime=record["datetime"], packet=json.loads(record["packet"]))
            measure_point_data.packets.append(packet_data)
        except ValidationError as e:
            logging.error(f"Ошибка проверки пакета данных {record}: {e}")
        except json.JSONDecodeError as e:
            logging.error(f"Ошибка декодирования JSON для пакета данных {record}: {e}")

    return measure_point_data
