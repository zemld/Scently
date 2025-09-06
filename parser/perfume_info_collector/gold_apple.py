import json
from bs4 import BeautifulSoup, element
import re
import requests
from models.perfume import Perfume
from pathlib import Path
from concurrent.futures import ThreadPoolExecutor, as_completed
from threading import Lock
from util.canonization.canonize_note import canonize_note

LOCK = Lock()
DIR = Path.cwd() / "collected_urls"
PROPERTIES_CNT = 14

SPLIT_NOTES_PATTERN = r",\s*|\s+и\s+|\s+-\s+|\s+–\s+"


def get_page_content(link: str) -> tuple[str, str]:
    headers = {"User-Agent": "Mozilla/5.0"}
    try:
        r = requests.get(link, timeout=20, headers=headers)
        r.raise_for_status()
        return (link, r.content)
    except Exception as e:
        print(f"Error fetching {link}: {e}")
        return (link, "")


def _is_brand_tag(tag: element.Tag) -> bool:
    return tag.has_attr("text") and tag.get("text") == "Бренд"


class Brand:
    name: str
    country: str

    def __init__(self, name: str = "Unknown", country: str = "Unknown"):
        self.name = name
        self.country = country


def get_brand_info(soup: BeautifulSoup) -> Brand:
    brand_tag = soup.find_all(_is_brand_tag)
    if not brand_tag:
        return ""
    try:
        brand_tag = brand_tag[0]
        brand_info = [
            tag.string.strip() for tag in brand_tag.find_all("div") if tag.string
        ]
        return Brand(brand_info[0], brand_info[1])
    except Exception as e:
        return Brand()


def _is_name_tag(tag: element.Tag) -> bool:
    return (
        tag.name == "span"
        and tag.has_attr("itemprop")
        and tag.get("itemprop") == "name"
        and tag.has_attr("class")
    )


def get_name(soup: BeautifulSoup) -> str:
    name_tag = soup.find_all(_is_name_tag)
    if not name_tag:
        return ""
    try:
        name_tag = name_tag[0]
        return name_tag.string.strip()
    except Exception as e:
        return ""


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
    notes_list = re.split(SPLIT_NOTES_PATTERN, notes)
    notes_list = [note.strip() for note in notes_list if note.strip()]
    notes_list[-1] = notes_list[-1].strip(".")
    canonized_notes = [canonize_note(note) for note in notes_list]
    canonized_notes = [note for note in canonized_notes if note]
    return canonized_notes


def get_volume(volume: str) -> int:
    int_rx = re.compile(r"\d+")
    return int(re.search(int_rx, volume).group(0)) if re.search(int_rx, volume) else 0


def get_properties(soup: BeautifulSoup) -> Perfume | None:
    properties = parse_properties(soup)
    if not properties or len(properties) < PROPERTIES_CNT:
        return None
    perfume = Perfume(
        perfume_type=properties[1].lower(),
        sex=properties[3].lower(),
        family=properties[5].lower(),
        upper_notes=get_notes(properties[7]),
        middle_notes=get_notes(properties[9]),
        base_notes=get_notes(properties[11]),
        volume=get_volume(properties[13]),
    )
    return perfume


def get_perfume(soup: BeautifulSoup) -> Perfume | None:
    perfume = get_properties(soup)
    if not perfume:
        return None
    perfume.brand = get_brand_info(soup).name
    perfume.name = get_name(soup)
    return perfume


def read_files_with_urls() -> list[str]:
    if not DIR.exists():
        return []
    files = list(DIR.glob("*.txt"))
    urls = []
    for file in files:
        with open(file, "r") as f:
            urls.extend(f.read().splitlines())
    return urls


def process_perfume(perfumes: list[Perfume], link: str, page: str):
    if not page:
        return
    soup = BeautifulSoup(page, "lxml")
    perfume = get_perfume(soup)
    perfume.link = link
    if perfume:
        with LOCK:
            perfumes.add(perfume)


def process_links(links: list[str]) -> set[Perfume]:
    perfumes = set()
    with ThreadPoolExecutor(max_workers=10) as executor:
        futures = {executor.submit(get_page_content, link): link for link in links}
        for fut in as_completed(futures):
            try:
                process_perfume(perfumes, fut.result()[0], fut.result()[1])
            except Exception as e:
                print(f"Error processing {fut}: {e}")
    return perfumes


def parse_pages_to_perfumes() -> list[Perfume]:
    links = read_files_with_urls()
    perfumes = process_links(links)
    return perfumes


if __name__ == "__main__":
    perfumes = parse_pages_to_perfumes()
    with open("goldapple_perfumes.json", "w") as f:
        json.dump([p.to_dict() for p in perfumes], f, ensure_ascii=False, indent=4)
