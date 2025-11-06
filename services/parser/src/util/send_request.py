import time
from typing import TypedDict

import requests
from bs4 import BeautifulSoup
from playwright.sync_api import (
    Error as PlaywrightError,
)
from playwright.sync_api import (
    Page,
    sync_playwright,
)
from playwright.sync_api import (
    TimeoutError as PlaywrightTimeoutError,
)


class PageConfig(TypedDict):
    selector: str
    needs_scroll: bool


def get_page(
    link: str,
    use_playwright: bool = False,
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
                page.goto(link, wait_until="domcontentloaded", timeout=60000)
            except (PlaywrightTimeoutError, PlaywrightError):
                pass

            if not _safe_wait_for_selector(
                page, "body", timeout=10000, fallback_sleep=0
            ):
                browser.close()
                return None

            _wait_for_page_content(page, link)

            try:
                page.wait_for_load_state("networkidle", timeout=10000)
            except (PlaywrightTimeoutError, PlaywrightError):
                try:
                    page.wait_for_load_state("load", timeout=5000)
                except (PlaywrightTimeoutError, PlaywrightError):
                    pass

            time.sleep(0.5)

            html_content = None
            max_retries = 3
            for attempt in range(max_retries):
                try:
                    html_content = page.content()
                    break
                except (PlaywrightError, Exception):
                    if attempt < max_retries - 1:
                        time.sleep(1)
                        continue
                    try:
                        page.wait_for_load_state("load", timeout=3000)
                        time.sleep(1)
                        html_content = page.content()
                    except Exception:
                        browser.close()
                        return None

            if html_content is None:
                browser.close()
                return None

            browser.close()

            return BeautifulSoup(html_content, _define_bs_type_from_link(link))

    except Exception:
        import traceback

        traceback.print_exc()
        return None


def _safe_wait_for_selector(
    page: Page, selector: str, timeout: int = 30000, fallback_sleep: int = 3
) -> bool:
    try:
        page.wait_for_selector(selector, timeout=timeout)
        return True
    except (PlaywrightTimeoutError, PlaywrightError):
        time.sleep(fallback_sleep)
        return False


def _scroll_page(page: Page, scroll_delay: float = 1.0) -> None:
    page.evaluate("window.scrollTo(0, document.body.scrollHeight)")
    time.sleep(scroll_delay)
    page.evaluate("window.scrollTo(0, 0)")
    time.sleep(scroll_delay)


def _wait_for_page_content(page: Page, link: str) -> None:
    site_configs: dict[str, dict[str, PageConfig]] = {
        "letu.ru": {
            "catalog": {
                "selector": "a.product-tile__item-container",
                "needs_scroll": True,
            },
            "product": {
                "selector": ".product-detail-v3__name-title, .product-detail-v3__offer",
                "needs_scroll": False,
            },
        },
        "goldapple.ru": {
            "catalog": {
                "selector": "article[data-scroll-id]",
                "needs_scroll": True,
            },
            "product": {
                "selector": "h1, [data-transaction-name='ga-pdp-title']",
                "needs_scroll": False,
            },
        },
    }

    site = None
    page_type = None

    for site_name in site_configs:
        if site_name in link:
            site = site_name
            break

    if site:
        if site == "letu.ru":
            if "/browse/" in link:
                page_type = "catalog"
            elif "/product/" in link:
                page_type = "product"
        elif site == "goldapple.ru":
            if "/parfjumerija" in link or "producttype" in link:
                page_type = "catalog"
            else:
                page_type = "product"

    if site and page_type:
        config = site_configs[site][page_type]
        found = _safe_wait_for_selector(page, config["selector"])
        if found and config["needs_scroll"]:
            _scroll_page(page)
    else:
        time.sleep(2)

    time.sleep(2)


def _define_bs_type_from_link(link: str) -> str:
    if link.endswith("xml"):
        return "xml"
    return "lxml"
