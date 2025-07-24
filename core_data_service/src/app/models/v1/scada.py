from datetime import datetime
from typing import Optional

from pydantic import BaseModel, Field, field_validator

MEASUREMENT_TYPE = {"ACCIDENT": "Авария", "MSD": "Мнемосхема"}


class LastData(BaseModel):
    timestamp: Optional[datetime] = Field(
        ..., description="Дата и время сбора пакета данных.", example="2024-09-03T12:34:56"
    )
    value: Optional[str] = Field(..., description="Значение.", example="2024-09-03T12:34:56")
    measurement_type_id: Optional[int] = Field(..., description="Уникальный идентификатор типа измерения.")
    measurement_name: Optional[str] = Field(..., description="Имя измерения.")
    measurement_type: Optional[str] = Field(..., description="Тип измерения.")

    @field_validator("measurement_type")
    def modify_measurement_type(cls, value):
        return MEASUREMENT_TYPE[value]


class GetLastData(BaseModel):
    name: Optional[str] = Field(..., description="Имя точки измерения.")
    packets: Optional[list[LastData]] = Field(..., description="Данные пакета.")
