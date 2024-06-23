FROM docker.io/python:3.10-alpine

WORKDIR /app

COPY app/python/requirements.txt /app

RUN pip install --trusted-host pypi.org --trusted-host files.pythonhosted.org -r requirements.txt

COPY app/python /app

RUN chmod +x /app/entrypoint.sh

ENTRYPOINT [ "sh", "/app/entrypoint.sh" ]

