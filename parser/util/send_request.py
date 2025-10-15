import requests
from bs4 import BeautifulSoup


def get_page(link: str) -> BeautifulSoup:
    headers = {"User-Agent": "Mozilla/5.0"}

    try:
        r = requests.get(link, headers=headers, timeout=30)
        r.raise_for_status()
        return BeautifulSoup(r.content, _define_bs_type_from_link(link))
    except Exception as e:
        print(e)
        return None


def _define_bs_type_from_link(link: str) -> str:
    if link.endswith("xml"):
        return "xml"
    return "lxml"
