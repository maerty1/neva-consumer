import logging

import asyncpg
import json
from pydantic import ValidationError

from app.models.report import BoilerRoomEngineerReport

# from aiocache import Cache, cached
import pickle
from cache import cache

from aiocache.serializers import PickleSerializer


async def boiler_room_engineer_report(self, year: int, month: int):
    self.logger.warning(f"🔄 Выполняем запрос в БД за {year}-{month}")
    cache_key = f"boiler_room_engineer_report:{year}-{month}"

    # Проверяем кеш перед запросом к БД
    cached_data = await cache.get(cache_key)
    if cached_data:
        self.logger.info(f"Использую кеш для {year}-{month}")
        try:
            return pickle.loads(cached_data)  # 🚀 Распаковываем через pickle
        except Exception as e:
            self.logger.error(f"Ошибка при десериализации данных из кеша: {e}")
            return []

    async with self.pool.acquire() as connection:
        query = """
WITH params AS (
    SELECT 
        $1::integer AS target_year,   
        $2::integer AS target_month      
),
date_series AS (
    -- Генерация последовательности дат для указанного месяца и года
    SELECT
        generate_series(
            make_date(p.target_year, p.target_month, 1),
            (make_date(p.target_year, p.target_month, 1) + interval '1 month - 1 day')::date,
            interval '1 day'
        )::date AS date
    FROM params p
),
current_year_data AS (
    -- Извлечение Q_in (Тепло прямой), Q_out (Тепло обратной), M_in (Расход прямой) и M_out (Расход обратной) за текущий месяц для measure_point_id = 779 (Котельная № 16)
    select
        dt.date,
        SUM(CASE WHEN mdd.measure_point_id = 779 AND elem->>'dataParameter' = 'Q_in' THEN (elem->>'value')::numeric END) AS q_in,
        SUM(CASE WHEN mdd.measure_point_id = 779 AND elem->>'dataParameter' = 'Q_out' THEN (elem->>'value')::numeric END) AS q_out,
        SUM(CASE WHEN mdd.measure_point_id = 779 AND elem->>'dataParameter' = 'M_in' THEN (elem->>'value')::numeric END) AS m_in,
        SUM(CASE WHEN mdd.measure_point_id = 779 AND elem->>'dataParameter' = 'M_out' THEN (elem->>'value')::numeric END) AS m_out
    FROM
        measure_points_data mdd
        CROSS JOIN params p
        CROSS JOIN LATERAL jsonb_array_elements(mdd.values::jsonb) AS elem
        join date_series dt 
        on (date_trunc('day', mdd.datetime) = dt.date AND EXTRACT(HOUR FROM mdd.datetime) < 8)
        or (date_trunc('day', mdd.datetime) = dt.date - interval '1 day' AND EXTRACT(HOUR FROM mdd.datetime) > 8)
    WHERE
        mdd.measure_point_id IN (779)
    GROUP BY
        dt.date
),
previous_year_data AS (
    -- Извлечение Q_in, Q_out, M_in и M_out за тот же месяц предыдущего года для measure_point_id = 779
    SELECT
        dt.date,
        SUM(CASE WHEN mdd.measure_point_id = 779 AND elem->>'dataParameter' = 'Q_in' THEN (elem->>'value')::numeric END) AS q_in,
        SUM(CASE WHEN mdd.measure_point_id = 779 AND elem->>'dataParameter' = 'Q_out' THEN (elem->>'value')::numeric END) AS q_out,
        SUM(CASE WHEN mdd.measure_point_id = 779 AND elem->>'dataParameter' = 'M_in' THEN (elem->>'value')::numeric END) AS m_in,
        SUM(CASE WHEN mdd.measure_point_id = 779 AND elem->>'dataParameter' = 'M_out' THEN (elem->>'value')::numeric END) AS m_out
    FROM
        measure_points_data mdd
        CROSS JOIN params p
        CROSS JOIN LATERAL jsonb_array_elements(mdd.values::jsonb) AS elem
        join date_series dt 
        on (date_trunc('day', mdd.datetime) = dt.date - interval '1 year' AND EXTRACT(HOUR FROM mdd.datetime) < 8)
        or (date_trunc('day', mdd.datetime) = dt.date - interval '1 year 1 day' AND EXTRACT(HOUR FROM mdd.datetime) > 8)
    WHERE
        mdd.measure_point_id IN (779)
    GROUP BY
        dt.date
),
current_year_data_day AS (
    -- Извлечение Q_in (Тепло прямой), Q_out (Тепло обратной), M_in (Расход прямой) и M_out (Расход обратной) за текущий месяц для measure_point_id = 779 (Котельная № 16)
    SELECT
        mdd.datetime::date AS date,
        MAX(CASE WHEN mdd.measure_point_id = 779 AND elem->>'dataParameter' = 'Q_in' THEN (elem->>'value')::numeric END) AS q_in,
        MAX(CASE WHEN mdd.measure_point_id = 785 AND elem->>'dataParameter' = 'M_in' THEN (elem->>'value')::numeric END) AS m_in_785,
        MAX(CASE WHEN mdd.measure_point_id = 779 AND elem->>'dataParameter' = 'Q_out' THEN (elem->>'value')::numeric END) AS q_out,
        MAX(CASE WHEN mdd.measure_point_id = 779 AND elem->>'dataParameter' = 'M_in' THEN (elem->>'value')::numeric END) AS m_in,
        MAX(CASE WHEN mdd.measure_point_id = 779 AND elem->>'dataParameter' = 'M_out' THEN (elem->>'value')::numeric END) AS m_out
    FROM
        measure_points_data_day mdd
        CROSS JOIN params p
        CROSS JOIN LATERAL jsonb_array_elements(mdd.values::jsonb) AS elem
    WHERE
        mdd.measure_point_id IN (779, 785)
        AND EXTRACT(YEAR FROM mdd.datetime) = p.target_year
        AND EXTRACT(MONTH FROM mdd.datetime) = p.target_month
    GROUP BY
        mdd.datetime::date
),
previous_year_data_day AS (
    -- Извлечение Q_in, Q_out, M_in и M_out за тот же месяц предыдущего года для measure_point_id = 779
    SELECT
        (mdd.datetime::date + interval '1 year')::date AS date,
        MAX(CASE WHEN mdd.measure_point_id = 779 AND elem->>'dataParameter' = 'Q_in' THEN (elem->>'value')::numeric END) AS q_in,
        MAX(CASE WHEN mdd.measure_point_id = 785 AND elem->>'dataParameter' = 'M_in' THEN (elem->>'value')::numeric END) AS m_in_785,
        MAX(CASE WHEN mdd.measure_point_id = 779 AND elem->>'dataParameter' = 'Q_out' THEN (elem->>'value')::numeric END) AS q_out,
        MAX(CASE WHEN mdd.measure_point_id = 779 AND elem->>'dataParameter' = 'M_in' THEN (elem->>'value')::numeric END) AS m_in,
        MAX(CASE WHEN mdd.measure_point_id = 779 AND elem->>'dataParameter' = 'M_out' THEN (elem->>'value')::numeric END) AS m_out
    FROM
        measure_points_data_day mdd
        CROSS JOIN params p
        CROSS JOIN LATERAL jsonb_array_elements(mdd.values::jsonb) AS elem
    WHERE
        mdd.measure_point_id IN (779, 785)
        AND EXTRACT(YEAR FROM mdd.datetime) = (p.target_year - 1)
        AND EXTRACT(MONTH FROM mdd.datetime) = p.target_month
    GROUP BY
        mdd.datetime::date
),
current_year_temp AS (
    -- Извлечение T_in (Температура прямой) за текущий месяц для measure_point_id = 787 (Береговая насосная)
    SELECT
        mdd.datetime::date AS date,
        MAX(CASE WHEN elem->>'dataParameter' = 'T_in' THEN (elem->>'value')::numeric END) AS t_in
    FROM
        measure_points_data_day mdd
        CROSS JOIN params p
        CROSS JOIN LATERAL jsonb_array_elements(mdd.values::jsonb) AS elem
    WHERE
        mdd.measure_point_id = 787
        AND EXTRACT(YEAR FROM mdd.datetime) = p.target_year
        AND EXTRACT(MONTH FROM mdd.datetime) = p.target_month
    GROUP BY
        mdd.datetime::date
),
previous_year_temp AS (
    -- Извлечение T_in за тот же месяц предыдущего года для measure_point_id = 787
    SELECT
        (mdd.datetime::date + interval '1 year')::date AS date,  -- Сдвиг даты на 1 год вперёд
        MAX(CASE WHEN elem->>'dataParameter' = 'T_in' THEN (elem->>'value')::numeric END) AS t_in
    FROM
        measure_points_data_day mdd
        CROSS JOIN params p
        CROSS JOIN LATERAL jsonb_array_elements(mdd.values::jsonb) AS elem
    WHERE
        mdd.measure_point_id = 787
        AND EXTRACT(YEAR FROM mdd.datetime) = (p.target_year - 1)
        AND EXTRACT(MONTH FROM mdd.datetime) = p.target_month
    GROUP BY
        mdd.datetime::date
),
current_year_temperatures AS (
    -- Извлечение T_Out и T_In для measure_point_id=779 (Котельная № 16) и T_In для measure_point_id=785 (Котельная № 16, новый (Подпитка)) за текущий год
    SELECT
        mdd.datetime::date AS date,
        MAX(CASE WHEN mdd.measure_point_id = 779 AND elem->>'dataParameter' = 'T_out' THEN (elem->>'value')::numeric END) AS t_out_779,
        MAX(CASE WHEN mdd.measure_point_id = 785 AND elem->>'dataParameter' = 'T_in' THEN (elem->>'value')::numeric END) AS t_in_785,
        MAX(CASE WHEN mdd.measure_point_id = 779 AND elem->>'dataParameter' = 'T_in' THEN (elem->>'value')::numeric END) AS t_in_779
    FROM
        measure_points_data_day mdd
        CROSS JOIN params p
        CROSS JOIN LATERAL jsonb_array_elements(mdd.values::jsonb) AS elem
    WHERE
        mdd.measure_point_id IN (779, 785)
        AND EXTRACT(YEAR FROM mdd.datetime) = p.target_year
        AND EXTRACT(MONTH FROM mdd.datetime) = p.target_month
    GROUP BY
        mdd.datetime::date
),
previous_year_temperatures AS (
    -- Извлечение T_Out и T_In для measure_point_id=779 и T_In для measure_point_id=785 за предыдущий год
    SELECT
        (mdd.datetime::date + interval '1 year')::date AS date,
        MAX(CASE WHEN mdd.measure_point_id = 779 AND elem->>'dataParameter' = 'T_out' THEN (elem->>'value')::numeric END) AS t_out_779,
        MAX(CASE WHEN mdd.measure_point_id = 785 AND elem->>'dataParameter' = 'T_in' THEN (elem->>'value')::numeric END) AS t_in_785,
        MAX(CASE WHEN mdd.measure_point_id = 779 AND elem->>'dataParameter' = 'T_in' THEN (elem->>'value')::numeric END) AS t_in_779
    FROM
        measure_points_data_day mdd
        CROSS JOIN params p
        CROSS JOIN LATERAL jsonb_array_elements(mdd.values::jsonb) AS elem
    WHERE
        mdd.measure_point_id IN (779, 785)
        AND EXTRACT(YEAR FROM mdd.datetime) = (p.target_year - 1)
        AND EXTRACT(MONTH FROM mdd.datetime) = p.target_month
    GROUP BY
        mdd.datetime::date
),
gas_cum_values AS (
    SELECT 
        "timestamp"::date AS date, 
        max(CASE WHEN sr.varname = 'DAILY_SPG761_P1_358_V' THEN sr.value::float END) AS gas_consumption_m3_8_to_7h,
        max(CASE WHEN sr.varname = 'DAILY_SPG761_P1_348_Q' THEN sr.value::float END) AS gas_consumption_hour
    FROM public.scada_rawdata sr
    GROUP BY "timestamp"::date
),
gas_consumption_m3_8_to_7h AS (
    -- Извлечение gas_consumption_m3_8_to_7h для каждого дня
    SELECT 
        ds.date, 
        (g1.gas_consumption_m3_8_to_7h - g2.gas_consumption_m3_8_to_7h - g1.gas_consumption_hour)::float AS gas_consumption_m3_8_to_7h
    FROM
        date_series ds   
        LEFT JOIN gas_cum_values g1 ON g1.date = ds.date
        JOIN gas_cum_values g2 on g1.date = g2.date + '1 day'::interval
),
electricity_consumption AS (
    -- Извлечение значений для расчета потребления электричества за текущий год 790 (котельная меркурий 233) 786 (Электроснабжение 1) 788 (Электроснабжение 2 / насосная)
    SELECT
        mdd.datetime::date AS date,
        MAX(CASE WHEN mdd.measure_point_id = 790 AND elem->>'dataParameter' = 'Ap' THEN (elem->>'value')::numeric END) AS ap_790,
        MAX(CASE WHEN mdd.measure_point_id = 786 AND elem->>'dataParameter' = 'Ap' THEN (elem->>'value')::numeric END) AS ap_786,
        MAX(CASE WHEN mdd.measure_point_id = 788 AND elem->>'dataParameter' = 'Ap' THEN (elem->>'value')::numeric END) AS ap_788
    FROM
        measure_points_data_day mdd
        CROSS JOIN params p
        CROSS JOIN LATERAL jsonb_array_elements(mdd.values::jsonb) AS elem
    WHERE
        mdd.measure_point_id IN (790, 786, 788)
        AND EXTRACT(YEAR FROM mdd.datetime) = p.target_year
        AND EXTRACT(MONTH FROM mdd.datetime) = p.target_month
    GROUP BY
        mdd.datetime::date
),
previous_year_electricity_consumption AS (
    -- Извлечение значений для расчета потребления электричества за предыдущий год
    SELECT
        (mdd.datetime::date + interval '1 year')::date AS date,
        MAX(CASE WHEN mdd.measure_point_id = 790 AND elem->>'dataParameter' = 'Ap' THEN (elem->>'value')::numeric END) AS ap_790,
        MAX(CASE WHEN mdd.measure_point_id = 786 AND elem->>'dataParameter' = 'Ap' THEN (elem->>'value')::numeric END) AS ap_786,
        MAX(CASE WHEN mdd.measure_point_id = 788 AND elem->>'dataParameter' = 'Ap' THEN (elem->>'value')::numeric END) AS ap_788
    FROM
        measure_points_data_day mdd
        CROSS JOIN params p
        CROSS JOIN LATERAL jsonb_array_elements(mdd.values::jsonb) AS elem
    WHERE
        mdd.measure_point_id IN (790, 786, 788)
        AND EXTRACT(YEAR FROM mdd.datetime) = (p.target_year - 1)
        AND EXTRACT(MONTH FROM mdd.datetime) = p.target_month
    GROUP BY
        mdd.datetime::date
),
outdoor_temp_data AS (
    -- Извлечение outdoor_temperature из static_measure_points_data для static_measure_point_id=1
    SELECT
        mdd.datetime::date AS date,
        MAX((mdd.value)::numeric) AS outdoor_temperature
    FROM
        static_measure_points_data mdd
        CROSS JOIN params p
    WHERE
        mdd.static_measure_point_id = 1
        AND EXTRACT(YEAR FROM mdd.datetime) = p.target_year
        AND EXTRACT(MONTH FROM mdd.datetime) = p.target_month
    GROUP BY
        mdd.datetime::date
),
previous_year_outdoor_temp_data AS (
    -- Извлечение outdoor_temperature из static_measure_points_data для static_measure_point_id=1
    SELECT
        (mdd.datetime::date + interval '1 year')::date AS date,
        MAX((mdd.value)::numeric) AS outdoor_temperature
    FROM
        static_measure_points_data mdd
        CROSS JOIN params p
    WHERE
        mdd.static_measure_point_id = 1
        AND EXTRACT(YEAR FROM mdd.datetime) = (p.target_year - 1)
        AND EXTRACT(MONTH FROM mdd.datetime) = p.target_month
    GROUP BY
        mdd.datetime::date
),
heating_network_operating_mode AS (
    -- Определение режима работы теплосети за текущий год на основе T_in для measure_point_id=779
    SELECT
        mdd.datetime::date AS date,
        (elem->>'value')::numeric AS heating_network_operating_mode
    FROM
        measure_points_data_day mdd
        CROSS JOIN params p
        CROSS JOIN LATERAL jsonb_array_elements(mdd.values::jsonb) AS elem
    WHERE
        mdd.measure_point_id = 779
        AND elem->>'dataParameter' = 'T_in'
        AND EXTRACT(YEAR FROM mdd.datetime) = p.target_year
        AND EXTRACT(MONTH FROM mdd.datetime) = p.target_month
),
previous_year_heating_network_operating_mode AS (
    -- Определение режима работы теплосети за предыдущий год (архив) на основе T_in для measure_point_id=779
    SELECT
        (mdd.datetime::date + interval '1 year')::date AS date,
        (elem->>'value')::numeric AS heating_network_operating_mode
    FROM
        measure_points_data_day mdd
        CROSS JOIN params p
        CROSS JOIN LATERAL jsonb_array_elements(mdd.values::jsonb) AS elem
    WHERE
        mdd.measure_point_id = 779
        AND elem->>'dataParameter' = 'T_in'
        AND EXTRACT(YEAR FROM mdd.datetime) = (p.target_year - 1)
        AND EXTRACT(MONTH FROM mdd.datetime) = p.target_month
)
SELECT
    ds.date,

    -- Отпуск тепловой энергии за текущий год
    coalesce(cy.q_in - cy.q_out, (cyd.q_in - cyd.q_out)/24*23) AS heat_energy_supply_8_to_7_gcal,

    -- Отпуск тепловой энергии за предыдущий год (архив)
    coalesce(py.q_in - py.q_out, (pyd.q_in - pyd.q_out)/24*23) AS archive_heat_energy_supply_8_to_7_gcal,

    -- Массовый расход теплоносителя за текущий год (тонн/ч)
    cyd.m_in / 24 AS heat_calculator_main_exit_boiler_room,

    -- Массовый расход теплоносителя за предыдущий год (тонн/ч)
    pyd.m_in / 24 AS archive_heat_calculator_main_exit_boiler_room,

    -- Подпитка за текущий год (м³/ч)
    round(cyd.m_in_785 / 24, 0)::int AS recharge_m3_hour,

    -- Подпитка за предыдущий год (архив, м³/ч)
    round(pyd.m_in_785 / 24, 0)::int AS archive_recharge_m3_hour,

    -- Температура в реке Плюсса за текущий год
    cyt.t_in AS temperature_river_plussa,

    -- Температура в реке Плюсса за предыдущий год (архив)
    pyt.t_in AS archive_temperature_river_plussa,

    -- Расход газа, м3 8 до 7ч
    sgc.gas_consumption_m3_8_to_7h,

    -- Потребление электричества (kW) за текущий год
    CASE 
        WHEN ec.ap_790 IS NULL AND ec.ap_786 IS NULL AND ec.ap_788 IS NULL THEN NULL
        ELSE COALESCE(ec.ap_790, 0) * 400 + COALESCE(ec.ap_786, 0) * 400 + COALESCE(ec.ap_788, 0) * 200
    END AS consumption_boiler_and_pump_house_electricity_kw,

    -- Потребление электричества (kW) за предыдущий год (архив)
    CASE 
        WHEN pec.ap_790 IS NULL AND pec.ap_786 IS NULL AND pec.ap_788 IS NULL THEN NULL
        ELSE COALESCE(pec.ap_790, 0) * 400 + COALESCE(pec.ap_786, 0) * 400 + COALESCE(pec.ap_788, 0) * 200
    END AS archive_consumption_boiler_and_pump_house_electricity_kw,

    -- Удельный расход газа
    sgc.gas_consumption_m3_8_to_7h / NULLIF(coalesce(cy.q_in - cy.q_out, cyd.q_in - cyd.q_out), 0) AS specific_gas_consumption,

    -- Удельный расход электричества
    ((COALESCE(ec.ap_790, 0) * 400) + (COALESCE(ec.ap_786, 0) * 400) + (COALESCE(ec.ap_788, 0) * 200)) / NULLIF(coalesce(cy.q_in - cy.q_out, cyd.q_in - cyd.q_out), 0) AS specific_consumption_electricity,

    -- Внешняя температура
    odt.outdoor_temperature AS outdoor_temperature,

    -- Внешняя температура за предыдущий год (архив)
    podt.outdoor_temperature AS archive_outdoor_temperature,
    
     -- Режим работы теплосети за текущий год
    hmo.heating_network_operating_mode,

    -- Режим работы теплосети за предыдущий год (архив)
    phmo.heating_network_operating_mode AS archive_heating_network_operating_mode

FROM
    date_series ds
    LEFT JOIN current_year_data cy ON ds.date = cy.date
    LEFT JOIN previous_year_data py ON ds.date = py.date
    LEFT JOIN current_year_data_day cyd ON ds.date = cyd.date
    LEFT JOIN previous_year_data_day pyd ON ds.date = pyd.date
    LEFT JOIN current_year_temp cyt ON ds.date = cyt.date
    LEFT JOIN previous_year_temp pyt ON ds.date = pyt.date
    LEFT JOIN current_year_temperatures ct ON ds.date = ct.date
    LEFT JOIN previous_year_temperatures pytemp ON ds.date = pytemp.date
    LEFT JOIN gas_consumption_m3_8_to_7h sgc ON ds.date = sgc.date
    LEFT JOIN electricity_consumption ec ON ds.date = ec.date
    LEFT JOIN previous_year_electricity_consumption pec ON ds.date = pec.date
    LEFT JOIN outdoor_temp_data odt ON ds.date = odt.date
    LEFT JOIN previous_year_outdoor_temp_data podt ON ds.date = podt.date
    LEFT JOIN heating_network_operating_mode hmo ON ds.date = hmo.date
    LEFT JOIN previous_year_heating_network_operating_mode phmo ON ds.date = phmo.date
ORDER BY
    ds.date;
"""
        try:
            records = await connection.fetch(query, year, month)
        except asyncpg.PostgresSyntaxError as e:
            logging.error(f"PostgreSQL Syntax Error: {e}")
            raise e

    heat_data_list = []
    for record in records:
        try:
            heat_data = BoilerRoomEngineerReport(
                date=record["date"],
                heat_energy_supply_8_to_7_gcal=record["heat_energy_supply_8_to_7_gcal"],
                archive_heat_energy_supply_8_to_7_gcal=record["archive_heat_energy_supply_8_to_7_gcal"],
                heat_calculator_main_exit_boiler_room=record["heat_calculator_main_exit_boiler_room"],
                archive_heat_calculator_main_exit_boiler_room=record["archive_heat_calculator_main_exit_boiler_room"],
                temperature_river_plussa=record["temperature_river_plussa"],
                archive_temperature_river_plussa=record["archive_temperature_river_plussa"],
                recharge_m3_hour=record["recharge_m3_hour"],
                archive_recharge_m3_hour=record["archive_recharge_m3_hour"],
                gas_consumption_m3_8_to_7h=record.get("gas_consumption_m3_8_to_7h"),
                consumption_boiler_and_pump_house_electricity_kw=record.get(
                    "consumption_boiler_and_pump_house_electricity_kw"
                ),
                archive_consumption_boiler_and_pump_house_electricity_kw=record.get(
                    "archive_consumption_boiler_and_pump_house_electricity_kw"
                ),
                specific_gas_consumption=record.get("specific_gas_consumption"),
                specific_consumption_electricity=record.get("specific_consumption_electricity"),
                outdoor_temperature=record.get("outdoor_temperature"),
                archive_outdoor_temperature=record.get("archive_outdoor_temperature"),
                heating_network_operating_mode=record.get("heating_network_operating_mode"),
                archive_heating_network_operating_mode=record.get("archive_heating_network_operating_mode"),
            )
            heat_data_list.append(heat_data)
        except ValidationError as e:
            logging.error(f"Ошибка проверки записи {record['date']}: {e}")

    heat_data_list = [BoilerRoomEngineerReport(**dict(record)) for record in records]

    # ✅ Сериализуем с помощью `pickle.dumps()`
    serialized_data = pickle.dumps(heat_data_list)
    await cache.set(cache_key, serialized_data, ttl=60 * 60 * 12)

    return heat_data_list
