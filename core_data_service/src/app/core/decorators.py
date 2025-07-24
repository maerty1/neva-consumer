import re
from functools import wraps

import asyncpg

from app.core import exceptions


def asyncpg_exc_handler(func):
    """
    Декоратор для обработки исключений asyncpg.
    """

    @wraps(func)
    async def wrapper(*args, **kwargs):
        try:
            result = await func(*args, **kwargs)
            return result

        except asyncpg.UniqueViolationError as e:
            field = parse_field_name(r"Key \((?P<field>[\w\s,]+)\)=", str(e))
            raise exceptions.UniqueDbException(detail=f"Уже существует {field} с таким значением")

        except asyncpg.NotNullViolationError as e:
            field = parse_field_name(r'column "(?P<field>\w+)" of', str(e))
            raise exceptions.NotNullDbException(field=field, details=str(e))

        except asyncpg.ForeignKeyViolationError as e:
            field = foreign_key_parser(r"Key \((\w+)\)=\((.+)\)", str(e))
            raise exceptions.ForeignKeyViolation(field=field, details=str(e))

        # except asyncpg.OperationalError as e:
        #     port_and_host = extract_host_port(str(e))
        #     raise exceptions.OperationalError(details=str(e), connection=port_and_host)

        except asyncpg.UndefinedTableError as e:
            regex = r'relation "(?P<field>[a-zA-Z0-9_]+)" does not exist'
            field_name = parse_field_name(regex, str(e))
            raise exceptions.NotFoundException(details=str(e), message="Таблица не найдена", field=field_name)

    return wrapper


def parse_field_name(regex: str, text: str) -> str:
    """
    Извлекает имя поля из сообщения об ошибке базы данных,
    используя регулярное выражение.
    """

    matches = re.search(regex, text)
    return matches.group("field") if matches else "Неизвестное поле"


def foreign_key_parser(regex: str, text: str) -> str:
    """
    Извлекает имена и значения ограничения внешнего ключа
    из сообщения об ошибке базы данных с помощью регулярного выражения.
    """

    matches = re.search(regex, text)
    return f"{matches.group(1)}={matches.group(2)}" if matches else "Неизвестное поле"


def extract_host_port(error_message: str) -> str:
    """
    Извлекает поля хоста и порта из сообщения об ошибке PostgreSQL.
    """

    host_regexes = [
        r'running on host "(.*?)"',
        r'connection to server at "(.*?)" \((?:\d{1,3}\.){3}\d{1,3}\)',
    ]
    port_regexes = [
        r"port (\d+) failed",
        r"on port (\d+)\?",
        r"TCP/IP connections on port (\d+)\?",
    ]

    for port_regex in port_regexes:
        port_matches = re.search(port_regex, error_message)
        if port_matches:
            port = port_matches.group(1)
            break
    else:
        port = "Неизвестный порт"

    for host_regex in host_regexes:
        host_matches = re.search(host_regex, error_message)
        if host_matches:
            host = host_matches.group(1)
            break
    else:
        host = "Неизвестный хост"

    return f"{host}: {port}"
