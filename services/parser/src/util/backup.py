import json
from pathlib import Path
from typing import Any, TypeVar

from src.models.perfume import PerfumeFromConcreteShop

T = TypeVar("T")


class BackupManager:
    def __init__(self, shop_name: str, backup_dir: Path | None = None):
        self.shop_name = shop_name
        if backup_dir is None:
            backup_dir = Path.cwd() / "data" / "backups"
        self.backup_dir = backup_dir
        self.backup_dir.mkdir(parents=True, exist_ok=True)

        self.links_file = self.backup_dir / f"{shop_name}_links.json"
        self.perfumes_file = self.backup_dir / f"{shop_name}_perfumes.json"

    def load_links(self) -> set[str]:
        if not self.links_file.exists():
            return set()

        try:
            with open(self.links_file, encoding="utf-8") as f:
                data = json.load(f)
                if isinstance(data, list):
                    return set(data)
                return set()
        except (OSError, json.JSONDecodeError):
            return set()

    def save_links(self, links: set[str]) -> None:
        with open(self.links_file, "w", encoding="utf-8") as f:
            json.dump(sorted(list(links)), f, indent=2, ensure_ascii=False)

    def add_links(self, new_links: list[str]) -> set[str]:
        existing_links = self.load_links()
        existing_links.update(new_links)
        self.save_links(existing_links)
        return existing_links

    def load_perfumes(self) -> dict[str, dict[str, Any]]:
        if not self.perfumes_file.exists():
            return {}

        try:
            with open(self.perfumes_file, encoding="utf-8") as f:
                data: Any = json.load(f)
                if isinstance(data, list):
                    perfumes_dict: dict[str, dict[str, Any]] = {}
                    for perfume_dict in data:
                        if isinstance(perfume_dict, dict):
                            key = self._get_perfume_key(perfume_dict)
                            if key:
                                perfumes_dict[key] = perfume_dict
                    return perfumes_dict
                return {}
        except (OSError, json.JSONDecodeError):
            return {}

    def _get_perfume_key(self, perfume_dict: dict[str, Any]) -> str:
        shop_info: Any = perfume_dict.get("shop_info", {})
        if not isinstance(shop_info, dict):
            shop_info = {}
        volumes: Any = shop_info.get("volumes_with_prices", [])
        if isinstance(volumes, list) and volumes:
            first_volume: Any = volumes[0]
            if isinstance(first_volume, dict):
                link: Any = first_volume.get("link", "")
                if isinstance(link, str) and link:
                    return link
        brand: Any = perfume_dict.get("brand", "")
        name: Any = perfume_dict.get("name", "")
        sex: Any = perfume_dict.get("sex", "")
        if (
            isinstance(brand, str)
            and isinstance(name, str)
            and isinstance(sex, str)
            and brand
            and name
            and sex
        ):
            return f"{brand.lower()} {name.lower()} {sex.lower()}"
        return ""

    def save_perfumes(self, perfumes_dict: dict[str, dict[str, Any]]) -> None:
        perfumes_list = list(perfumes_dict.values())
        with open(self.perfumes_file, "w", encoding="utf-8") as f:
            json.dump(perfumes_list, f, indent=4, ensure_ascii=False)

    def add_perfumes(
        self, new_perfumes: list[PerfumeFromConcreteShop]
    ) -> dict[str, dict[str, Any]]:
        existing_perfumes = self.load_perfumes()

        for perfume in new_perfumes:
            perfume_dict = perfume.to_dict()
            key = self._get_perfume_key(perfume_dict)
            if key:
                existing_perfumes[key] = perfume_dict

        self.save_perfumes(existing_perfumes)
        return existing_perfumes
