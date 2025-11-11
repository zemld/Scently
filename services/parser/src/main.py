import asyncio
import json
import os
import time
from pathlib import Path

import httpx
from apscheduler.schedulers.background import BackgroundScheduler

from src.app import get_all_perfumes, unite_perfumes
from src.canonization import Canonizer, NoteCanonizer
from src.models import PerfumeWithUnitedShops
from src.scraping.gold_apple import GoldApplePageParser, GoldAppleScrapper
from src.scraping.letu import LetuPageParser, LetuScrapper
from src.scraping.randewoo import RandewooPageParser, RandewooScrapper
from src.scraping.scrapper import Scrapper
from src.util import setup_logger

logger = setup_logger(
    __name__, log_file=Path.cwd() / "logs" / f"{__name__.split('.')[-1]}.log"
)

BATCH_SIZE = 100

UPLOAD_PERFUME_INFO = "http://perfume:8000/v1/perfumes/update"

TYPE_CANONIZER = Canonizer(Path.cwd() / "data/types")
SEX_CANONIZER = Canonizer(Path.cwd() / "data/sex")
FAMILY_CANONIZER = Canonizer(Path.cwd() / "data/families")
NOTES_CANONIZER = NoteCanonizer(Path.cwd() / "data/notes")


def collect_and_store_perfumes(shop_name: str, scrapper: Scrapper) -> None:
    perfumes = scrapper.scrap_all_accuratly()
    with open(f"data/collected_perfumes/{shop_name}_perfumes.json", "w") as f:
        json.dump(
            [perfume.to_dict() for perfume in perfumes],
            f,
            indent=4,
            ensure_ascii=False,
        )
    logger.info(
        f"Collected and stored perfumes | shop_name={shop_name} | count={len(perfumes)}"
    )


def collect_and_store_all_perfumes() -> None:
    scrappers = {
        "goldapple": GoldAppleScrapper(
            page_parser=GoldApplePageParser(
                type_canonizer=TYPE_CANONIZER,
                sex_canonizer=SEX_CANONIZER,
                family_canonizer=FAMILY_CANONIZER,
                notes_canonizer=NOTES_CANONIZER,
            ),
        ),
        "randewoo": RandewooScrapper(
            "https://randewoo.ru",
            page_parser=RandewooPageParser(
                type_canonizer=TYPE_CANONIZER,
                sex_canonizer=SEX_CANONIZER,
                family_canonizer=FAMILY_CANONIZER,
                notes_canonizer=NOTES_CANONIZER,
            ),
        ),
        "letu": LetuScrapper(
            "https://www.letu.ru",
            page_parser=LetuPageParser(
                type_canonizer=TYPE_CANONIZER,
                sex_canonizer=SEX_CANONIZER,
                family_canonizer=FAMILY_CANONIZER,
                notes_canonizer=NOTES_CANONIZER,
            ),
        ),
    }
    for shop_name, scrapper in scrappers.items():
        collect_and_store_perfumes(shop_name, scrapper)


async def try_to_upload_perfumes_to_database(
    perfumes: list[PerfumeWithUnitedShops],
    sem: asyncio.Semaphore,
    try_number: int = 0,
    max_retries: int = 3,
    client: httpx.AsyncClient | None = None,
) -> bool:
    if client is None:
        client = httpx.AsyncClient()
        return await try_to_upload_perfumes_to_database(
            perfumes, sem, try_number, max_retries, client
        )
    body = {"perfumes": [perfume.to_dict() for perfume in perfumes]}
    try:
        async with sem:
            response = await client.post(
                UPLOAD_PERFUME_INFO,
                json=body,
                timeout=30,
                headers={
                    "Authorization": f"Bearer {os.getenv('PERFUME_INTERNAL_TOKEN')}"
                },
            )
            response.raise_for_status()
            return True
    except (httpx.HTTPStatusError, httpx.RequestError) as e:
        if try_number == max_retries - 1:
            raise e
        await asyncio.sleep(2**try_number)
        return await try_to_upload_perfumes_to_database(
            perfumes, sem, try_number + 1, max_retries, client
        )


async def update_perfumes_in_database(
    collect: bool = True, batch_size: int = BATCH_SIZE
) -> None:
    if collect:
        await asyncio.to_thread(collect_and_store_all_perfumes)
    perfumes = unite_perfumes(get_all_perfumes(Path.cwd() / "data/collected_perfumes"))
    batches = [
        perfumes[i : i + batch_size] for i in range(0, len(perfumes), batch_size)
    ]

    sem = asyncio.Semaphore(10)
    async with httpx.AsyncClient() as client:
        await asyncio.gather(
            *(
                try_to_upload_perfumes_to_database(batch, sem=sem, client=client)
                for batch in batches
            ),
            return_exceptions=True,
        )


if __name__ == "__main__":
    logger.info("Starting scheduler")
    scheduler = BackgroundScheduler()

    def sheduled_update_wrapper() -> None:
        asyncio.run(update_perfumes_in_database())

    scheduler.add_job(
        sheduled_update_wrapper,
        "interval",
        days=5,
        replace_existing=True,
        max_instances=1,
    )
    scheduler.start()
    asyncio.run(update_perfumes_in_database(collect=False))
    try:
        while True:
            time.sleep(1)
    except KeyboardInterrupt:
        scheduler.shutdown()
        logger.info("Scheduler shutdown")
