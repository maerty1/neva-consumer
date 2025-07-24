from abc import ABC, abstractmethod
from typing import Any

from fastapi.responses import Response


class ReportFormatStrategy(ABC):
    @abstractmethod
    def _convert(self, data: dict[str, Any], **kwargs) -> Any:
        pass

    @abstractmethod
    def _get_response(self, data: Any, filename: str = "") -> Response:
        pass

    def convert_and_response(self, data: dict[str, Any], filename: str = "", **kwargs) -> Response:
        if not data:
            return Response(status_code=204)

        converted_data = self._convert(data, **kwargs)
        return self._get_response(converted_data, filename)
