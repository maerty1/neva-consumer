from dataclasses import dataclass

from fastapi import HTTPException


class ApiException(Exception):
    pass


@dataclass
class BadRequestException(ApiException):
    message: str = "Неверные данные запроса"
    code: str = "B001"
    details: list | str | None = None


class NoPermissionException(HTTPException):
    def __init__(self, detail: str):
        super().__init__(status_code=403, detail=detail)


@dataclass
class DbException(ApiException):
    message: str
    field: str | None = None
    code: str = "D001"
    details: str | None = None


@dataclass
class NotNullDbException(DbException):
    code: str = "D002"
    message: str = "Указан null, но поле не может быть null"


class UniqueDbException(HTTPException):
    def __init__(self, detail: str):
        super().__init__(status_code=400, detail=detail)


@dataclass
class ForeignKeyViolation(DbException):
    code: str = "D004"
    message: str = (
        "Добавление данного объекта невозможно, "
        "так как он ссылается на несуществующий объект. "
        "Проверьте его наличие в базе данных"
    )


class NotFoundException(HTTPException):
    def __init__(self, detail: str):
        super().__init__(status_code=404, detail=detail)


@dataclass
class OperationalError(DbException):
    code: str = "D006"
    message: str = "Не удалось подключиться к базе данных.  Проверьте настройки подключения"
    connection: str = ""


@dataclass
class InvalidTextRepresentation(DbException):
    code: str = "D007"
    message: str = "Неверное значение Enum"
    field: str | None = None


@dataclass
class UnexpectedServerException(ApiException):
    message: str = "Непредвиденная ошибка сервера"
    code: str = "S001"
    details: str | None = None


@dataclass
class InvalidJsonPathException(DbException):
    code: str = "D008"
    message: str = "Невалидный JSONPath параметр"
    field: str | None = None
