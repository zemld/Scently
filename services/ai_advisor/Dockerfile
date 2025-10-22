FROM python:3.13-slim-bookworm

ENV PYTHONPATH="/ai_advisor/src"

WORKDIR /ai_advisor

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY requirements.txt ./
RUN pip install --no-cache-dir -r requirements.txt

COPY . .

CMD ["uvicorn", "src.main:advisor", "--host", "0.0.0.0", "--port", "8000"]