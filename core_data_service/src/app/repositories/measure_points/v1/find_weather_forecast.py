import logging
from datetime import datetime, timedelta


async def find_weather_forecast(self):
    today_date = datetime.now().date()
    yesterday_date = today_date - timedelta(days=1)
    tomorrow_date = today_date + timedelta(days=1)

    query = """
    WITH wdr_agg AS (
        SELECT 
            dateutc::date AS date,
            AVG(tempf) AS avg_tempf,
            AVG(baromabsin) AS avg_baromabsin,
            AVG(humidityin) AS avg_humidityin,
            AVG(windspeedmph) AS avg_windspeedmph,
            MAX(dateutc) AS latest_dateutc
        FROM 
            public.weather_data_raw
        WHERE
            dateutc::date IN ($1, $2, $3)
        GROUP BY 
            dateutc::date
    ),
    wd_latest AS (
        SELECT DISTINCT ON (date::date)
            *
        FROM 
            public.weather_data
        WHERE
            date::date IN ($1, $2, $3)
        ORDER BY 
            date::date, date DESC
    )
    SELECT
        wd_latest.outdoor_temperature AS temp_avg,
        wd_latest.pressure AS pressure_avg,
        ROUND((wdr_agg.avg_tempf - 32) / 1.8 * 10) / 10 AS temp,
        ROUND(wdr_agg.avg_baromabsin * 25.3285279248736) AS pressure,
        wdr_agg.avg_humidityin AS humidity,
        wd_latest.humidity AS humidity_avg,
        ROUND(wdr_agg.avg_windspeedmph * 1.60934 * 10) / 10 AS wind_speed,
        (wd_latest.max_wind_speed + wd_latest.min_wind_speed) / 2 AS wind_speed_avg,
        CASE 
            WHEN wd_latest.date::date = $1 THEN wdr_agg.latest_dateutc 
            ELSE wd_latest.date 
        END AS date
    FROM
        wd_latest
    LEFT JOIN
        wdr_agg 
        ON wdr_agg.date = wd_latest.date::date
    WHERE
        wd_latest.date::date IN ($1, $2, $3);
    """

    dates = [today_date, tomorrow_date, yesterday_date]

    async with self.pool.acquire() as connection:
        try:
            records = await connection.fetch(query, *dates)
            result = {"today": None, "tomorrow": None, "yesterday": None}

            for record in records:
                record_datetime = record["date"]
                record_date = record_datetime.date()

                data = {
                    "temp_avg": round(record["temp_avg"]) if record["temp_avg"] is not None else None,
                    "temp": record["temp"],
                    "pressure_avg": round(record["pressure_avg"]) if record["pressure_avg"] is not None else None,
                    "pressure": record["pressure"],
                    "humidity": round(record["humidity"]) if record["humidity"] is not None else None,
                    "humidity_avg": record["humidity_avg"],
                    "wind_speed": record["wind_speed"],
                    "wind_speed_avg": record["wind_speed_avg"],
                    "date": record_datetime.strftime("%Y-%m-%dT%H:%M:%SZ") if record_datetime else None,
                }

                if record_date == today_date:
                    result["today"] = data
                elif record_date == tomorrow_date:
                    result["tomorrow"] = data
                elif record_date == yesterday_date:
                    result["yesterday"] = data

            return result

        except Exception as e:
            logging.error(f"Ошибка выполнения запроса: {e}")
            return {"today": None, "tomorrow": None, "yesterday": None}
