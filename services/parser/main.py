import asyncio
import json
import os
import time
from pathlib import Path

import httpx

from canonization.canonizer import Canonizer
from canonization.notes_canonizer import NoteCanonizer
from canonization.type_canonizer import TypeCanonizer
from scraping.gold_apple.gold_apple_parser import GoldApplePageParser
from scraping.gold_apple.scrapper import GoldAppleScrapper

UPLOAD_PERFUME_INFO = "http://perfume:8089/v1/perfumes/update"


def get_hard_update_key() -> str:
    hard_update_key: str
    try:
        file = os.getenv("HARD_UPDATE_PASSWORD_FILE")
        if file is None:
            hard_update_key = "default_key"
        else:
            with open(file) as f:
                hard_update_key = f.read().strip()
    except Exception as e:
        print(f"Error reading hard update key: {e}")
        hard_update_key = "default_key"
    return hard_update_key


async def _upload_perfumes_async(
    url: str, payload: dict[str, str | list[str] | int]
) -> str:
    timeout = httpx.Timeout(10.0, read=30.0)
    async with httpx.AsyncClient(timeout=timeout) as client:
        r = await client.post(
            url, json=payload, params={"hard": True, "password": get_hard_update_key()}
        )
        r.raise_for_status()
        text = r.text
        return text if isinstance(text, str) else str(text)


def update_perfumes(index: int | None = None) -> None:
    gold_apple_scrapper = GoldAppleScrapper(
        GoldApplePageParser(
            brand_canonizer=None,
            name_canonizer=None,
            type_canonizer=TypeCanonizer(Path.cwd() / "data/types"),
            sex_canonizer=Canonizer(Path.cwd() / "data/sex"),
            family_canonizer=Canonizer(Path.cwd() / "data/families"),
            notes_canonizer=NoteCanonizer(Path.cwd() / "data/notes"),
        )
    )
    if index:
        perfumes = gold_apple_scrapper.scrap_sitemap(index)
        with open(
            f"data/collected_perfumes/gold_apple_{index+1}.json", "w"
        ) as checkpoint:
            json.dump([p.to_dict() for p in perfumes], checkpoint)
    else:
        perfumes = gold_apple_scrapper.scrap_all_accuratly()
        with open("data/collected_perfumes/ga.json", "w") as f:
            json.dump([p.to_dict() for p in perfumes], f)
    # with open("goldapple_perfumes.json", "r") as f:
    #     perfumes = json.load(f)
    # payload = {"perfumes": perfumes}

    # try:
    #     asyncio.run(_upload_perfumes_async(UPLOAD_PERFUME_INFO, payload))
    # except Exception as e:
    #     print(f"Error uploading perfumes: {e}")


if __name__ == "__main__":
    with open("data/all_perfumes.json") as f:
        perfumes = json.load(f)

    time.sleep(15)

    try:
        asyncio.run(_upload_perfumes_async(UPLOAD_PERFUME_INFO, {"perfumes": perfumes}))
    except Exception as e:
        print(f"Error uploading perfumes: {e}")
    # parser = argparse.ArgumentParser(description="Update perfumes from sitemap")
    # parser.add_argument("-n", type=int, help="Sitemap index number")
    # args = parser.parse_args()

    # if args.n is not None:
    #     update_perfumes(args.n)
    # else:
    #     update_perfumes()

    # scheduler = BackgroundScheduler()
    # scheduler.add_job(update_perfumes, "cron", hour=3, minute=0, day_of_week="sun")
    # scheduler.start()

    # print("added cron")
    # try:
    #     while True:
    #         pass
    # except (KeyboardInterrupt, SystemExit):
    #     scheduler.shutdown()
