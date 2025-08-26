from models.perfume import Perfume
from url_collectors.gold_apple import collect_links_from_sitemap
from perfume_info_collector.gold_apple import parse_pages_to_perfumes
from apscheduler.schedulers.background import BackgroundScheduler
import asyncio
import httpx
import os

UPLOAD_PERFUME_INFO = "http://db-api:8089/v1/perfumes/update"


def get_hard_update_key() -> str:
    hard_update_key: str
    try:
        file = os.getenv("HARD_UPDATE_PASSWORD_FILE")
        with open(file, "r") as f:
            hard_update_key = f.read().strip()
    except Exception as e:
        print(f"Error reading hard update key: {e}")
        hard_update_key = "default_key"
    return hard_update_key


async def _upload_perfumes_async(url: str, payload: dict) -> str:
    timeout = httpx.Timeout(10.0, read=30.0)
    async with httpx.AsyncClient(timeout=timeout) as client:
        r = await client.post(
            url, json=payload, params={"hard": True, "password": get_hard_update_key()}
        )
        r.raise_for_status()
        return r.text


def update_perfumes():
    collect_links_from_sitemap()
    perfumes = parse_pages_to_perfumes()
    print(f"Before glue: {len(perfumes)} perfumes")

    payload = {"perfumes": [p.to_dict() for p in perfumes]}

    try:
        asyncio.run(_upload_perfumes_async(UPLOAD_PERFUME_INFO, payload))
    except Exception as e:
        print(f"Error uploading perfumes: {e}")


if __name__ == "__main__":
    update_perfumes()

    scheduler = BackgroundScheduler()
    scheduler.add_job(update_perfumes, "cron", hour=3, minute=0, day_of_week="sun")
    scheduler.start()

    print("added cron")
    try:
        while True:
            pass
    except (KeyboardInterrupt, SystemExit):
        scheduler.shutdown()
