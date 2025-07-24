import logging

import asyncpg
from pydantic import ValidationError

from app.models.v2.report import BoilerRoomEngineerReportV2


async def boiler_room_engineer_report_v2(self, temperature: float, direction: int):
    async with self.pool.acquire() as connection:
        # await connection.execute("REFRESH MATERIALIZED VIEW report_mat_view;")
        query = """
WITH params AS (
    SELECT   
        $1::numeric(4,1) AS reference_temperature,
        $2::int4 as direction
),
weather as (
    select 
        "date"::date,
        humidity,
        sum(outdoor_temperature) over (order by "date" desc RANGE BETWEEN interval '6 days' PRECEDING AND CURRENT row)/7 outdoor_temperature_avg,
        (min_wind_speed+max_wind_speed)/2::int avg_wind_speed,
        wind_direction
    from weather_data w
    where extract('hour' from "date") = 9
),
source as (
    select
        m.date,
        m.outdoor_temperature,
        m.heat_energy_supply_8_to_7_gcal,
        abs(reference_temperature - m.outdoor_temperature) t_delta,
        (reference_temperature - m.outdoor_temperature) t_delta_2
    from params
    join report_mat_view m on true
    join weather w 
    on w.date = m.date
),
row_order as (
    select 
    s.date,
    t_delta_2,
    row_number() over (order by t_delta) rn_t_delta,
    row_number() over (order by heat_energy_supply_8_to_7_gcal) rn_supply_avg_delta
from source s
join params p on true
where 
case 
when direction = 0 then t_delta_2 between -0.5 and 0.5
when direction > 0 then t_delta_2 > 0.5
when direction < 0 then t_delta_2 < -0.5
end)
select 
	m.outdoor_temperature,
    m.date, 
    m.heat_energy_supply_8_to_7_gcal::numeric(7,2),
    w.outdoor_temperature_avg::numeric(3,1),
    w.humidity,
    w.avg_wind_speed,
    w.wind_direction,
    m.temperature_river_plussa,
    m.heating_network_operating_mode::numeric(7,2),
    m.heat_calculator_main_exit_boiler_room::numeric(8,2),
    m.recharge_m3_hour,
    m.gas_consumption_m3_8_to_7h,
    m.consumption_boiler_and_pump_house_electricity_kw, 
    m.specific_consumption_electricity::numeric(7,2),
    m.specific_gas_consumption::numeric(7,2) 
from report_mat_view m
join weather w 
on w.date = m.date
join row_order r 
on m.date = r.date
join params on true
order by case when direction = 0
then m.outdoor_temperature
else null
end desc, rn_t_delta, rn_supply_avg_delta
limit 30;
"""
        try:
            records = await connection.fetch(query, temperature, direction)
        except asyncpg.PostgresSyntaxError as e:
            logging.error(f"PostgreSQL Syntax Error: {e}")
            raise e

        heat_data_list = []
        for record in records:
            try:
                heat_data = BoilerRoomEngineerReportV2(
                    outdoor_temperature=record.get("outdoor_temperature"),
                    date=record.get("date"),
                    heat_energy_supply_8_to_7_gcal=record.get("heat_energy_supply_8_to_7_gcal"),
                    outdoor_temperature_avg=record.get("outdoor_temperature_avg"),
                    humidity=record.get("humidity"),
                    avg_wind_speed=record.get("avg_wind_speed"),
                    wind_direction=record.get("wind_direction"),
                    temperature_river_plussa=record.get("temperature_river_plussa"),
                    heating_network_operating_mode=record.get("heating_network_operating_mode"),
                    heat_calculator_main_exit_boiler_room=record.get("heat_calculator_main_exit_boiler_room"),
                    recharge_m3_hour=record.get("recharge_m3_hour"),
                    gas_consumption_m3_8_to_7h=record.get("gas_consumption_m3_8_to_7h"),
                    consumption_boiler_and_pump_house_electricity_kw=record.get("consumption_boiler_and_pump_house_electricity_kw"), 
                    specific_consumption_electricity=record.get("specific_consumption_electricity"),
                    specific_gas_consumption=record.get("specific_gas_consumption"),
                )
                heat_data_list.append(heat_data)
            except ValidationError as e:
                logging.error(f"Ошибка проверки записи {record['date']}: {e}")

        return heat_data_list

