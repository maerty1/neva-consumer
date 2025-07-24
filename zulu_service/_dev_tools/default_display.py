# TechnicalAt: ат
# Kilogram: кг
# GigaCalorie: ГКал
# CubicMeter м3
# Ton: т
# KilowattHour: кВт·ч
# KilovarHour: кВар·ч

# P_delta — Перепад давления
# P_in, P_out, P_cw — Давление
# V_in, V_out — Объемный расход
# M_in, M_out, M_delta — Массовый расход
# Q_in, Q_out, Q_delta — Тепловая энергия

default_units = {
  "Heat": {
    "dataParameterUnits": {
      "Pressure": "TechnicalAt",
      "Mass": "Kilogram",
      "Volume": "CubicMeter",
      "Heat": "GigaCalorie",
      "PressureDrop": "TechnicalAt"
    }
  },
  "HotWater": {
    "dataParameterUnits": {
      "Pressure": "TechnicalAt",
      "Mass": "Ton",
      "Volume": "CubicMeter",
      "Heat": "GigaCalorie",
      "PressureDrop": "TechnicalAt"
    }
  },
  "ColdWater": {
    "dataParameterUnits": {
      "Pressure": "TechnicalAt",
      "Mass": "Ton",
      "Volume": "CubicMeter",
      "Heat": "GigaCalorie",
      "PressureDrop": "TechnicalAt"
    }
  },
  "Steam": {
    "dataParameterUnits": {
      "Pressure": "TechnicalAt",
      "Mass": "Ton",
      "Volume": "CubicMeter",
      "Heat": "GigaCalorie",
      "PressureDrop": "TechnicalAt"
    }
  },
  "Gas": {
    "dataParameterUnits": {
      "Pressure": "TechnicalAt",
      "Mass": "Ton",
      "Volume": "CubicMeter",
      "PressureDrop": "TechnicalAt"
    }
  },
  "Electricity": {
    "dataParameterUnits": {
      "ActiveElectricalEnergy": "KilowattHour",
      "ReactiveElectricalEnergy": "KilovarHour"
    }
  },
  "Sewage": {
    "dataParameterUnits": {
      "Pressure": "TechnicalAt",
      "Mass": "Ton",
      "Volume": "CubicMeter",
      "Heat": "GigaCalorie",
      "PressureDrop": "TechnicalAt"
    }
  },
  "Control": {
    "dataParameterUnits": {
      "Pressure": "TechnicalAt",
      "Mass": "Ton",
      "Volume": "CubicMeter",
      "Heat": "GigaCalorie",
      "PressureDrop": "TechnicalAt"
    }
  }
}

# Что, если сначала в одних ед изм приходило, а теперь в других?

# Функция: на вход (T_in, T_out, Q_in, Q_out) на выход — Т, c, Атм

{
  1: "TechnicalAt",
  2: "Ton",
  3: "CubicMeter",
  4: "GigaCalorie",
  5: "Kilogram"
}

{
  "V_in": "Volume",
  "V_out": "Volume", 
  "M_in": "Mass",
  "M_out": "Mass",
  "M_delta": "Mass",
}