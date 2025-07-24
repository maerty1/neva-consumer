from datetime import date
from typing import Optional

from pydantic import BaseModel, Field


class BoilerRoomEngineerReportV2(BaseModel):
    date: date
    outdoor_temperature: Optional[float] = Field(default=None, description="Температура наружного воздуха")
    heat_energy_supply_8_to_7_gcal: Optional[float] = Field(
        default=None, description="Отпуск тепловой энергии 8_7 Гкал"
    )
    humidity: Optional[float] = Field(default=None, description="Влажность")
    outdoor_temperature_avg: Optional[float] = Field(
        default=None, description="Средняя температура за прошедшую неделю"
    )
    avg_wind_speed: Optional[float] = Field(default=None, description="Средняя скорость ветра")
    wind_direction: Optional[float] = Field(default=None, description="Направление ветра")
    temperature_river_plussa: Optional[float] = Field(default=None, description="Температура в р.Плюсса")
    recharge_m3_hour: Optional[float] = Field(default=None, description="Подпитка, м3/час")
    gas_consumption_m3_8_to_7h: Optional[float] = Field(default=None, description="Потребление газа")
    consumption_boiler_and_pump_house_electricity_kw: Optional[float] = Field(
        default=None, description="Потербление электричества"
    )
    heating_network_operating_mode: Optional[float] = Field(default=None, description="Режим работы теплосети ℃")
    heat_calculator_main_exit_boiler_room: Optional[float] = Field(
        default=None, description="Массовый расход теплоносителя, почасово тонн / ч"
    )
    specific_consumption_electricity: Optional[float] = Field(default=None, description="Удельный расход электричества")
    specific_gas_consumption: Optional[float] = Field(default=None, description="Удельный расход газа")

    # fmt: off
    # temperature_river_plussa: Optional[float] = Field(default=None, description="Температура в реке Плюсса")
    # recharge_m3_hour: Optional[float] = Field(default=None, description="Подпитка м3/час")
    # gas_consumption_m3_8_to_7h: Optional[float] = Field(default=None, description="Расход газа, м3 8 до 7ч")
    # consumption_boiler_and_pump_house_electricity_kw: Optional[float] = Field(default=None, description="Расход электроэнергии котельной и насосной, кВт")

    # Пока пустые
    # archive_gas_consumption_m3_8_to_7h: Optional[float] = Field(default=None, description="Расход газа, м3 8 до 7ч (архив)")
    # archive_heat_energy_supply_8_to_7_gcal: Optional[float] = Field(default=None, description="Отпуск тепловой энергии 8_7 Гкал (архив)")
    # archive_heat_calculator_main_exit_boiler_room: Optional[float] = Field(default=None, description="Массовый расход теплоносителя, почасово тонн / ч (архив)")
    # archive_temperature_river_plussa: Optional[float] = Field(default=None, description="Температура в реке Плюсса (архив)")
    # archive_outdoor_temperature: Optional[float] = Field(default=None, description="Температура наружного воздуха (архив)")
    # archive_recharge_m3_hour: Optional[float] = Field(default=None, description="Подпитка м3/час (архив)")
    # archive_consumption_boiler_and_pump_house_electricity_kw: Optional[float] = Field(default=None, description="Расход электроэнергии котельной и насосной, кВт (архив)")
    # archive_heating_network_operating_mode: Optional[float] = Field(default=None, description="Режим работы теплосети ℃ (архив)")
    # fmt: on

    @staticmethod
    def round_value(value: Optional[float], places: int = 2) -> Optional[float]:
        if value is not None:
            return round(value, 1)
        return value

    def __init__(self, **data):
        super().__init__(**data)
        # fmt: off
        self.outdoor_temperature = self.round_value(self.outdoor_temperature)
        self.heat_energy_supply_8_to_7_gcal = self.round_value(self.heat_energy_supply_8_to_7_gcal)
        self.outdoor_temperature_avg = self.round_value(self.outdoor_temperature_avg)
        self.humidity = self.round_value(self.humidity)
        self.avg_wind_speed = self.round_value(self.avg_wind_speed)
        self.wind_direction = self.round_value(self.wind_direction)
        self.temperature_river_plussa = self.round_value(self.temperature_river_plussa)
        self.recharge_m3_hour = self.round_value(self.recharge_m3_hour)
        self.gas_consumption_m3_8_to_7h = self.round_value(self.gas_consumption_m3_8_to_7h)
        self.consumption_boiler_and_pump_house_electricity_kw = self.round_value(self.consumption_boiler_and_pump_house_electricity_kw)
        self.heating_network_operating_mode = self.round_value(self.heating_network_operating_mode)
        self.heat_calculator_main_exit_boiler_room = self.round_value(self.heat_calculator_main_exit_boiler_room)
        self.specific_consumption_electricity = self.round_value(self.specific_consumption_electricity)
        self.specific_gas_consumption =  self.round_value(self.specific_gas_consumption)

        # self.consumption_boiler_and_pump_house_electricity_kw = self.round_value(self.consumption_boiler_and_pump_house_electricity_kw)
        # self.gas_consumption_m3_8_to_7h = self.round_value(self.gas_consumption_m3_8_to_7h)
        # self.temperature_river_plussa = self.round_value(self.temperature_river_plussa)

        # self.archive_heat_energy_supply_8_to_7_gcal = self.round_value(self.archive_heat_energy_supply_8_to_7_gcal)
        # self.archive_heat_calculator_main_exit_boiler_room = self.round_value(self.archive_heat_calculator_main_exit_boiler_room)
        # self.archive_temperature_river_plussa = self.round_value(self.archive_temperature_river_plussa)
        # self.recharge_m3_hour = self.round_value(self.recharge_m3_hour)
        # self.archive_recharge_m3_hour = self.round_value(self.archive_recharge_m3_hour)
        # self.archive_consumption_boiler_and_pump_house_electricity_kw = self.round_value(self.archive_consumption_boiler_and_pump_house_electricity_kw)
        # self.archive_heating_network_operating_mode = self.round_value(self.archive_heating_network_operating_mode)
        # fmt: on


class MergedBoilerRoomEngineerReportV2(BaseModel):
    less: list[BoilerRoomEngineerReportV2]
    current: list[BoilerRoomEngineerReportV2]
    more: list[BoilerRoomEngineerReportV2]
