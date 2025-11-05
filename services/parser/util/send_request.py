import time

import requests
from bs4 import BeautifulSoup
from playwright.sync_api import (
    TimeoutError as PlaywrightTimeoutError,
)
from playwright.sync_api import sync_playwright


def get_page(
    link: str,
    use_playwright: bool = True,
) -> BeautifulSoup | None:
    if use_playwright:
        return _get_page_with_playwright(link)
    return _get_page_with_requests(link)


def _get_page_with_requests(link: str) -> BeautifulSoup | None:
    headers = {"User-Agent": "Mozilla/5.0"}

    try:
        r = requests.get(link, headers=headers, timeout=30)
        r.raise_for_status()
        return BeautifulSoup(r.content, _define_bs_type_from_link(link))
    except Exception as e:
        print(e)
        return None


def _get_page_with_playwright(link: str) -> BeautifulSoup | None:
    try:
        with sync_playwright() as p:
            browser = p.chromium.launch(
                headless=True,
                args=[
                    "--disable-blink-features=AutomationControlled",
                    "--disable-dev-shm-usage",
                    "--no-sandbox",
                ],
            )

            context = browser.new_context(
                viewport={"width": 1920, "height": 1080},
                user_agent=(
                    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) "
                    "AppleWebKit/537.36 (KHTML, like Gecko) "
                    "Chrome/120.0.0.0 Safari/537.36"
                ),
                locale="ru-RU",
            )

            context.add_init_script(
                """
                Object.defineProperty(navigator, 'webdriver', {
                    get: () => undefined
                });
            """
            )

            page = context.new_page()

            try:
                page.goto(link, wait_until="networkidle", timeout=60000)
            except PlaywrightTimeoutError:
                pass

            try:
                page.wait_for_selector("body", timeout=15000)
                page.wait_for_selector("dl", timeout=10000)
            except PlaywrightTimeoutError:
                pass

            time.sleep(3)

            html_content = page.content()

            browser.close()

            return BeautifulSoup(html_content, _define_bs_type_from_link(link))

    except Exception:
        import traceback

        traceback.print_exc()
        return None


def _define_bs_type_from_link(link: str) -> str:
    if link.endswith("xml"):
        return "xml"
    return "lxml"
