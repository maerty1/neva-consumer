from datetime import datetime
from enum import Enum
from typing import Optional

from pydantic import BaseModel, Field


class PeriodType(Enum):
    DAY = "day"
    HOUR = "hour"


class MeasurePointTypeParameters(BaseModel):
    type_key: str = Field(...)
    minzoom: int = Field(...)
    maxzoom: int = Field(...)


class MeasurePointType(BaseModel):
    title: str
    parameters: MeasurePointTypeParameters


class GetMeasurePointsWithLastData(BaseModel):
    measure_point_id: int = Field(..., description="Уникальный идентификатор измерительной точки.")
    address: str = Field(..., description="Адрес объекта", example="Улица Пушкина")
    title: str = Field(..., description="КУМИ")
    lat: Optional[float] = Field(None, description="Широта местоположения измерительной точки.", example=55.7558)
    lon: Optional[float] = Field(None, description="Долгота местоположения измерительной точки.", example=37.6173)
    full_title: str = Field(..., description="Полное название измерительной точки.", example="ул. Пушкина д. 2а - ВУМИ")
    packet_datetime: datetime = Field(
        ..., description="Дата и время сбора последнего пакета данных.", example="2024-09-03T00:00:00"
    )
    type: Optional[MeasurePointType] = None
    packet: list[dict] = Field(
        ...,
        description="Последний пакет данных",
        example=[
            {
                "value": 18.8,
                "isCalc": True,
                "isReset": False,
                "isBad": False,
                "isInterpolated": False,
                "dataParameter": "T_in",
            }
        ],
    )


class PacketData(BaseModel):
    packet_datetime: datetime = Field(
        ..., description="Дата и время сбора пакета данных.", example="2024-09-03T12:34:56"
    )
    packet: list[dict] = Field(
        ...,
        description="Данные пакета.",
        example=[
            {
                "value": 18.8,
                "isCalc": True,
                "isReset": False,
                "isBad": False,
                "isInterpolated": False,
                "dataParameter": "T_in",
            }
        ],
    )


class GetMeasurePointWithLastData(BaseModel):
    measure_point_id: int = Field(..., description="Уникальный идентификатор измерительной точки.")
    lat: Optional[float] = Field(None, description="Широта местоположения измерительной точки.", example=55.7558)
    lon: Optional[float] = Field(None, description="Долгота местоположения измерительной точки.", example=37.6173)
    full_title: str = Field(..., description="Полное название измерительной точки.", example="ул. Пушкина д. 2а - ВУМИ")
    type: Optional[MeasurePointType] = None
    packets: list[PacketData] = Field(..., description="Список пакетов данных с их датами и содержимым.")


class MeasureTypes(BaseModel):
    id: int
    measure_type: str = Field(None, description="Название переменной в ЛЭРС", example="T_in")
    title: str = Field(None, description="Название на русском", example="Температура прямой")
    unit: str = Field(None, description="Единицы измерения", example="°C")


class GetPointsDataRequestMeasurement(BaseModel):
    i: str = Field(None, description="In", example="T_in")
    o: str = Field(None, description="Out", example="T_out")


class GetPointsDataRequest(BaseModel):
    elem_id: int
    measurements: dict[str, GetPointsDataRequestMeasurement]


class GetPointsDataResponseMeasurement(BaseModel):
    i: Optional[float]
    o: Optional[float]

    @staticmethod
    def round_value(value: Optional[float], places: int = 2) -> Optional[float]:
        if value is not None:
            return round(value, 1)
        return value

    def __init__(self, **data):
        super().__init__(**data)
        self.i = self.round_value(self.i)
        self.o = self.round_value(self.o)


class GetPointsDataResponse(BaseModel):
    elem_id: int
    iscopied: bool = False
    measurements: dict[str, GetPointsDataResponseMeasurement]


class GetPointsDataResponseMeasurementWithTimestamp(BaseModel):
    i: float | None
    o: float | None
    timestamp: datetime | None

    @staticmethod
    def round_value(value: Optional[float], places: int = 2) -> Optional[float]:
        if value is not None:
            return round(value, 2)
        return value

    def __init__(self, **data):
        super().__init__(**data)
        self.i = self.round_value(self.i)
        self.o = self.round_value(self.o)


class GetPointsDataResponseWithHistory(BaseModel):
    elem_id: int
    measurements: dict[str, list[GetPointsDataResponseMeasurementWithTimestamp]]


class GetPointsDataKeyvalueRequest(BaseModel):
    elem_id: int
    measurements: dict[str, str] = Field(
        ..., example={"Alarm_Zagaz.Metan_1%": "Авария X", "Alarm_Obschaja_avarija": "Авария Y"}
    )


class GetPointsDataKeyvalueResponse(BaseModel):
    measurements: dict[str, dict[str, str]] = Field(
        ..., example={"1": {"ASAP_1_Alarm_Obschaja_avarija": "1", "ASAP_1_Alarm_Zagaz.Metan_1%": "0"}}
    )
