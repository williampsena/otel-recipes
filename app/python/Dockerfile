FROM docker.io/python:3.10-alpine

ARG PORT=8000

WORKDIR /app

COPY app/python/requirements.txt /app

RUN pip install --trusted-host pypi.org --trusted-host files.pythonhosted.org -r requirements.txt

COPY app/python /app

RUN chmod +x /app/entrypoint.sh

EXPOSE ${PORT}

ENTRYPOINT [ "sh", "/app/entrypoint.sh" ]

