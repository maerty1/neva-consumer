import logging

import fastapi.exceptions
from fastapi.responses import JSONResponse
from pydantic import ValidationError

from app.core import exceptions


def no_permission_exception_handler(request, exc: exceptions.NoPermissionException) -> JSONResponse:
    msg = {
        "Message": f"User {request.headers['x-telematix-login']} has no permission to access {request.url.path}",
        "Details": exc.details,
    }
    logging.error(msg)

    return JSONResponse(
        status_code=403,
        content={
            "Message": exc.message,
            "Code": exc.code,
        },
    )


# def unique_db_exception_handler(request, exc: exceptions.UniqueDbException) -> JSONResponse:
#     msg = {"Message": exc.message, "Field": exc.field, "Details": exc.details}
#     logging.error(msg)
#
#     return JSONResponse(
#         status_code=400,
#         content={
#             "Message": exc.message,
#             "Code": exc.code,
#             "Field": exc.field,
#         },
#     )


def not_null_db_exception_handler(request, exc: exceptions.NotNullDbException) -> JSONResponse:
    msg = {"Message": exc.message, "Field": exc.field, "Details": exc.details}
    logging.error(msg)

    return JSONResponse(
        status_code=404,
        content={
            "Message": exc.message,
            "Code": exc.code,
            "Field": exc.field,
        },
    )


def foreign_key_db_exception_handler(request, exc: exceptions.ForeignKeyViolation) -> JSONResponse:
    msg = {"Message": exc.message, "Field": exc.field, "Details": exc.details}
    logging.error(msg)

    return JSONResponse(
        status_code=400,
        content={
            "Message": exc.message,
            "Code": exc.code,
            "Field": exc.field,
        },
    )


def operational_error_exception_handler(request, exc: exceptions.OperationalError) -> JSONResponse:
    msg = {"Message": exc.message, "Details": exc.details}
    logging.error(msg)

    return JSONResponse(
        status_code=500,
        content={
            "Message": exc.message,
            "Code": exc.code,
            **({"Field": exc.connection} if exc.connection != "Неизвестный хост: Неизвестный порт" else {}),
        },
    )


def bad_request_exception_handler(request, exc: exceptions.BadRequestException) -> JSONResponse:
    msg = {"Message": exc.message, "Details": exc.details}
    logging.error(msg)

    return JSONResponse(
        status_code=400,
        content={
            "Message": exc.message,
            "Code": exc.code,
            "Details": exc.details,
        },
    )


def unexpected_server_exception_handler(request, exc: exceptions.UnexpectedServerException):
    msg = {"Message": exc.message, "Details": exc.details}
    logging.error(msg)

    return JSONResponse(
        status_code=500,
        content={
            "Message": exc.message,
            "Code": exc.code,
        },
    )


def validation_error_exception_handler(request, exc: ValidationError) -> JSONResponse:
    errors = []
    for error in exc.errors():
        field = ".".join(str(x) for x in error["loc"])
        message = error["msg"]
        errors.append({"field": field, "message": message})

    return bad_request_exception_handler(request, exceptions.BadRequestException(details=errors))


def configure_exceptions_handlers(app):
    exceptions_to_handlers = [
        (exceptions.BadRequestException, bad_request_exception_handler),
        (exceptions.ForeignKeyViolation, foreign_key_db_exception_handler),
        (exceptions.NoPermissionException, no_permission_exception_handler),
        (exceptions.NotNullDbException, not_null_db_exception_handler),
        (exceptions.OperationalError, operational_error_exception_handler),
        # (exceptions.UniqueDbException, unique_db_exception_handler),
        (exceptions.UnexpectedServerException, unexpected_server_exception_handler),
        (ValidationError, validation_error_exception_handler),
        (fastapi.exceptions.ResponseValidationError, validation_error_exception_handler),
    ]

    for exception, handler in exceptions_to_handlers:
        app.add_exception_handler(exception, handler)
