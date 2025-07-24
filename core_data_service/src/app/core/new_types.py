from uuid import UUID


class UUIDStr(str):
    @classmethod
    def __get_validators__(cls):
        yield cls.validate

    @classmethod
    def validate(cls, v):
        if isinstance(v, UUID):
            return v
        elif isinstance(v, str):
            try:
                return UUID(v)
            except ValueError:
                raise ValueError(f"Некорректный UUID: {v}")
        raise ValueError("UUID должен быть строкой или экземпляром UUID")

    # Метод для представления UUID в виде строки при сериализации
    def __str__(self):
        return str(self.uuid)
