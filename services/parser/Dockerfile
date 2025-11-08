FROM python:3.13-slim-bookworm

ENV PATH="/root/.local/bin:${PATH}"

WORKDIR /perfume-updater

RUN apt-get update && apt-get install -y --no-install-recommends \
    curl ca-certificates \
    && rm -rf /var/lib/apt/lists/*

RUN curl -LsSf https://astral.sh/uv/install.sh | sh

COPY . .

RUN uv sync --no-dev

RUN uv run --no-default-groups playwright install --with-deps chromium

CMD ["uv", "run", "--no-default-groups", "python", "-m", "src.main"]