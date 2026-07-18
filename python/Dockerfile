FROM python:3.12-slim

WORKDIR /app

COPY pyproject.toml requirements-dev.txt ./
RUN pip install --no-cache-dir -r requirements-dev.txt

COPY . .

CMD ["pytest", "-v"]
