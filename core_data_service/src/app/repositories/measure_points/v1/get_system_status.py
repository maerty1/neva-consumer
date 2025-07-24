import logging
from datetime import datetime
from typing import Optional

from pydantic import BaseModel


class SystemStatusResponse(BaseModel):
    fresh_count: int
    outdated_count: int
    success_percentage: float
    last_sync: Optional[datetime] = None


async def get_system_status(self, user_id: int, fresh_interval_days: int = 1, success_interval_days: int = 7):
    # Запрос для подсчёта свежих и устаревших точек измерения
    mp_status_query = f"""
    WITH last_data AS (
      SELECT
        mp.id,
        MAX(mpd.datetime) AS last_datetime
      FROM public.measure_points mp
      LEFT JOIN public.measure_points_data_day mpd
        ON mp.id = mpd.measure_point_id
      WHERE mp.account_id = $1  
      GROUP BY mp.id
    )
    SELECT
      SUM(CASE WHEN last_datetime >= NOW() - INTERVAL '{fresh_interval_days} day' THEN 1 ELSE 0 END) AS fresh_count,
      SUM(CASE WHEN last_datetime < NOW() - INTERVAL '{fresh_interval_days} day' OR last_datetime IS NULL THEN 1 ELSE 0 END) AS outdated_count
    FROM last_data;
    """

    # Запрос для вычисления процента успешных синхронизаций за последние 7 дней
    success_percentage_query = f"""
    SELECT
      ROUND(
        100.0 * SUM(CASE WHEN message = 'Успешная синхронизация' THEN 1 ELSE 0 END)
        / COUNT(*),
      2) AS success_percentage
    FROM public.accounts_sync_log
    WHERE created_at >= NOW() - INTERVAL '{success_interval_days} days'
      AND account_id = $1;
    """

    # Запрос для получения времени последней синхронизации (UTC)
    last_sync_query = """
    SELECT 
      MAX(created_at) AT TIME ZONE 'UTC' AS last_sync
    FROM public.accounts_sync_log
    WHERE account_id = $1;
    """

    status = SystemStatusResponse(success_percentage=0, fresh_count=0, outdated_count=0, last_sync=None)

    async with self.pool.acquire() as connection:
        try:
            mp_row = await connection.fetchrow(mp_status_query, user_id)
            status.fresh_count = mp_row["fresh_count"]
            status.outdated_count = mp_row["outdated_count"]
        except Exception as e:
            logging.error(f"Ошибка выполнения запроса статуса измерений: {e}")
            return status

        try:
            sp_row = await connection.fetchrow(success_percentage_query, user_id)
            status.success_percentage = float(sp_row["success_percentage"])
        except Exception as e:
            logging.error(f"Ошибка выполнения запроса процента успешных синхронизаций: {e}")
            return status

        try:
            ls_row = await connection.fetchrow(last_sync_query, user_id)
            status.last_sync = ls_row["last_sync"]
        except Exception as e:
            logging.error(f"Ошибка выполнения запроса последней синхронизации: {e}")
            return status

    return status
