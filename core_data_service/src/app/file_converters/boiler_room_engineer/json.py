from fastapi.responses import JSONResponse, Response

from app.file_converters.base_strategy import ReportFormatStrategy
from app.models.report import BoilerRoomEngineerReport


class JSONStrategy(ReportFormatStrategy):
    def _convert(self, data: list[BoilerRoomEngineerReport], **kwargs) -> list[dict]:
        for item in data:
            if not isinstance(item, BoilerRoomEngineerReport):
                print(f"⚠️ Ошибка! Ожидался объект BoilerRoomEngineerReport, но получен {type(item)}: {item}")

        return [
            {**item.dict(), "date": item.date.isoformat() if item.date else None}
            for item in data
            if isinstance(item, BoilerRoomEngineerReport)
        ]

    def _get_response(self, data: bytes, filename="") -> Response:
        return JSONResponse(data)
