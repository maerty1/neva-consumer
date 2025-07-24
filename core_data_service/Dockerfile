FROM python:3.11

WORKDIR /code

RUN apt-get update && \
    apt-get install -y wkhtmltopdf
 
RUN chmod 0700 /tmp

RUN pip install poetry

COPY ./pyproject.toml ./poetry.lock /code/
RUN poetry config virtualenvs.create false && poetry install --only main

COPY ./src /code/src
ENV PYTHONPATH=/code/src:/code/src/app

WORKDIR /code/src

CMD ["uvicorn", "app:app", "--host", "0.0.0.0", "--port", "8002"]
