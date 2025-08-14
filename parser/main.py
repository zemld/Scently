from perfume import Perfume
from url_collectors.gold_apple import collect_links_from_sitemap
from perfume_info_collector.gold_apple import parse_pages_to_perfumes
from apscheduler.schedulers.background import BackgroundScheduler


def update_perfumes():
    collect_links_from_sitemap()
    perfumes = parse_pages_to_perfumes()
    # TODO: добавить обращение в api


if __name__ == "__main__":
    scheduler = BackgroundScheduler()
    scheduler.add_job(update_perfumes, "cron", hour=3, minute=0, day_of_week="sun")
    scheduler.start()

    print("added cron")
    try:
        while True:
            pass
    except (KeyboardInterrupt, SystemExit):
        scheduler.shutdown()
