from pathlib import Path
import json
import sys
from redis import Redis
import argparse

IS_TESTING = True


def _read_env_file(path: Path) -> dict:
    with open(path, "r") as f:
        return {line.split("=")[0]: line.split("=")[1].strip("\n") for line in f.readlines()}


def _get_redis_client() -> Redis:
    redis_env = _read_env_file(Path(__file__).parent.parent / "secrets/config_storage.env")
    r = Redis(
        host="localhost" if IS_TESTING else str(redis_env.get("REDIS_HOST")),
        port=int(redis_env.get("REDIS_PORT", 6379)),
        password=None if IS_TESTING else redis_env.get("REDIS_PASSWORD"),
    )
    return r


def read_config(path: Path) -> dict:
    with open(path, "r") as f:
        return json.load(f)


def upload_config(service_name: str, config: dict) -> None:
    r = _get_redis_client()
    r.set(service_name, json.dumps(config).encode("utf-8"))


if __name__ == "__main__":
    argparser = argparse.ArgumentParser()
    argparser.add_argument("-f", type=str)
    args = argparser.parse_args()
    if not args.f:
        print("Error: -f is required")
        sys.exit(1)
    try:
        config = read_config(Path(__file__).parent / "configs" / f"{args.f}.json")
        print(config)
        upload_config(args.f, config)
        print("Config uploaded successfully")
    except Exception as e:
        print(f"Error uploading config: {e}")
        sys.exit(1)

    try:
        r = _get_redis_client()
        print(r.get(args.f))
    except Exception as e:
        print(f"Error reading config: {e}")
        sys.exit(1)
