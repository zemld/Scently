from bs4 import BeautifulSoup
import re
import requests
from perfume import Perfume

# TODO: remove test links
links = [
    "https://goldapple.ru/7041000021-chance-eau-fraiche",
    "https://goldapple.ru/26250800001-spicebomb-fresh",
    "https://goldapple.ru/80790100005-amber",
    "https://goldapple.ru/7041100009-allure-homme-sport",
    "https://goldapple.ru/7041100012-allure-homme-sport",
    "https://goldapple.ru/80790100007-prada-amber",
    "https://goldapple.ru/80790100002-amber",
    "https://goldapple.ru/7040500014-n-19-poudre",
    "https://goldapple.ru/7040100046-n-5",
    "https://goldapple.ru/7042000021-bleu-de-chanel",
    "https://goldapple.ru/7041100013-allure-homme-sport",
    "https://goldapple.ru/7040600018-coco-mademoiselle",
    "https://goldapple.ru/7041100001-allure-homme-sport",
    "https://goldapple.ru/7040100053-n-5",
    "https://goldapple.ru/7040400006-allure",
    "https://goldapple.ru/7042000009-bleu-de-chanel",
    "https://goldapple.ru/7040900001-allure-sensuelle",
    "https://goldapple.ru/80791400002-candy-florale",
]


def get_page_content(link: str) -> str:
    headers = {"User-Agent": "Mozilla/5.0"}
    try:
        r = requests.get(link, timeout=20, headers=headers)
        r.raise_for_status()
        return r.text
    except Exception as e:
        print(f"Error fetching {link}: {e}")
        return ""


# TODO: добавить парсинг бренда + объема
# TODO: добавить перевод в объект

if __name__ == "__main__":
    page_content: str
    with open("test.txt", "r") as c:
        page_content = c.read()
    soup = BeautifulSoup(page_content, "lxml")
    properties_title_rx = re.compile("\s*Подробные характеристики\s*", re.I)
    properties_title = soup.find_all(string=properties_title_rx)
    section = properties_title[0].parent.parent if properties_title else None
    properties = section.find_all("span")
    for prop in properties:
        print(prop.string.strip())
