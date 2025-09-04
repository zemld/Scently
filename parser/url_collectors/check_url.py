import requests
from bs4 import BeautifulSoup, element
import time
from requests import HTTPError
import re

PRODUCT_URL_RE = re.compile(r"^https://goldapple\.ru/\d{6,}-[a-z0-9-]+$", re.IGNORECASE)

_PERFUME = [
    re.compile(r"\bпарфюмированная\s+вода\b", re.IGNORECASE),
    re.compile(r"\bпарфюмерная\s+вода\b", re.IGNORECASE),
    re.compile(r"\bтуалетная\s+вода\b", re.IGNORECASE),
    re.compile(r"\bэкстракт\s+духов\b", re.IGNORECASE),
    re.compile(r"\bдухи\b", re.IGNORECASE),
    re.compile(r"\beau\s*de\s*parfum\b", re.IGNORECASE),
    re.compile(r"\beau\s*de\s*toilette\b", re.IGNORECASE),
    re.compile(r"\beau\s*de\s*cologne\b", re.IGNORECASE),
    re.compile(r"\bEDP\b", re.IGNORECASE),
    re.compile(r"\bEDT\b", re.IGNORECASE),
    re.compile(r"\bEDC\b", re.IGNORECASE),
]


def is_product_url(href: str) -> bool:
    return bool(PRODUCT_URL_RE.match(href))


def _is_perfume_html(html: str) -> bool:
    soup = BeautifulSoup(html, "lxml")
    return any(rx.search(soup.title.string.strip()) for rx in _PERFUME)


class PerfumeResponse:
    url: str
    is_perfume: bool
    is_retry_needed: bool


def is_perfume_url(url: str, timeout: int = 12) -> PerfumeResponse:
    time.sleep(1)
    answer = PerfumeResponse()
    answer.url = url
    answer.is_retry_needed = False
    try:
        r = requests.get(url, timeout=timeout, headers={"User-Agent": "Mozilla/5.0"})
        r.raise_for_status()
        answer.is_perfume = _is_perfume_html(r.content)
        return answer
    except HTTPError as e:
        print(f"HTTPError: {e}")
        answer.is_retry_needed = e.response.status_code == 429
        answer.is_perfume = False
    except Exception:
        return False
