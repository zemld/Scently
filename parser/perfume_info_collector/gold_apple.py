from bs4 import BeautifulSoup, element
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


# TODO: добавить парсинг бренда


def _is_volume_box(tag: element.Tag) -> bool:
    return tag.has_attr("style") and tag["style"] == "--icon-gap:5px;"


def parse_volume(soup: BeautifulSoup) -> list[int]:
    volume_variants = soup.find_all(_is_volume_box)
    volumes = []
    for variant in volume_variants:
        spans = variant.find_all("span")
        for span in spans:
            if span.string:
                volumes.append(
                    [int(volume) for volume in re.findall(r"\d+", span.string)]
                )
    volumes = sorted([volume[0] for volume in volumes if volume])
    return volumes


def parse_properties(soup: BeautifulSoup) -> list[str]:
    properties_title_rx = re.compile("Подробные характеристики", re.I)
    properties_title = soup.find_all(string=properties_title_rx)
    if not properties_title:
        return []
    try:
        section = properties_title[0].parent.parent if properties_title else None
        raw_properties = section.find_all("span")
        properties = []
        for prop in raw_properties:
            properties.append(prop.string.strip())
        return properties
    except Exception as e:
        return []


def get_notes(notes: str) -> list[str]:
    notes = notes.lower()
    notes_list = notes.split(",")
    notes_list = [note.strip() for note in notes_list if note.strip()]
    return notes_list


def get_properties(soup: BeautifulSoup) -> Perfume | None:
    properties = parse_properties(soup)
    if not properties:
        return None
    perfume = Perfume(
        perfume_type=properties[1].lower(),
        sex=properties[3].lower(),
        family=properties[5].lower(),
        upper_notes=get_notes(properties[7]),
        middle_notes=get_notes(properties[9]),
        base_notes=get_notes(properties[11]),
    )
    return perfume


if __name__ == "__main__":
    page_content: str
    with open("test.txt", "r") as c:
        page_content = c.read()
    soup = BeautifulSoup(page_content, "lxml")
    perfume = get_properties(soup)
    perfume.volume = parse_volume(soup)
    print(perfume)
