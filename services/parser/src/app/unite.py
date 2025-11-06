import json
from pathlib import Path

from src.models import PerfumeFromConcreteShop, PerfumeKey, PerfumeWithUnitedShops


def get_all_perfumes(path: Path) -> list[PerfumeFromConcreteShop]:
    all_perfumes: list[PerfumeFromConcreteShop] = []
    for json_file in path.glob("*.json"):
        with open(json_file, encoding="utf-8") as f:
            perfumes_data = json.load(f)
            if isinstance(perfumes_data, list):
                all_perfumes.extend(
                    [PerfumeFromConcreteShop.from_dict(p) for p in perfumes_data]
                )
            else:
                all_perfumes.append(PerfumeFromConcreteShop.from_dict(perfumes_data))
    return all_perfumes


def _get_priority(shop_info_list: list[PerfumeFromConcreteShop.ShopInfo]) -> int:
    shop_priorities = {
        "gold apple": 1,
        "randewoo": 2,
        "letu": 3,
    }
    return min(
        [
            shop_priorities.get(shop_info.shop_name.lower(), 100)
            for shop_info in shop_info_list
        ]
    )


def unite_perfumes(
    perfumes: list[PerfumeFromConcreteShop],
) -> list[PerfumeWithUnitedShops]:
    united_perfumes = dict[PerfumeKey, PerfumeWithUnitedShops]()
    for perfume in perfumes:
        key = PerfumeKey(perfume)
        if key in united_perfumes:
            current_priority = _get_priority(united_perfumes[key].shops)
            new_priority = _get_priority([perfume.shop_info])
            if new_priority < current_priority:
                united_perfumes[key].properties = perfume.properties
            united_perfumes[key].shops.append(perfume.shop_info)
        else:
            united_perfumes[key] = PerfumeWithUnitedShops(perfume)

    return list(united_perfumes.values())
