import logging
from logging.handlers import RotatingFileHandler
from pathlib import Path


def setup_logger(
    name: str,
    log_level: int = logging.INFO,
    log_file: str | Path | None = None,
    log_format: str = "%(asctime)s - %(name)s - %(levelname)s - %(message)s",
) -> logging.Logger:
    logger = logging.getLogger(name)

    if logger.hasHandlers():
        return logger

    logger.setLevel(log_level)
    logger.propagate = False

    console_handler = logging.StreamHandler()
    _set_handler_settings(console_handler, log_level, log_format)
    logger.addHandler(console_handler)

    if log_file:
        try:
            log_path = Path(log_file)
            log_path.parent.mkdir(parents=True, exist_ok=True)
            file_handler = RotatingFileHandler(
                str(log_path),
                maxBytes=10 * 1024 * 1024,
                backupCount=5,
                encoding="utf-8",
            )
            _set_handler_settings(file_handler, log_level, log_format)
            logger.addHandler(file_handler)
        except (OSError, PermissionError) as e:
            logger.warning(f"Failed to create file handler for {log_file}: {e}")

    return logger


def _set_handler_settings(
    handler: logging.Handler, log_level: int, log_format: str
) -> None:
    formatter = logging.Formatter(log_format)
    handler.setLevel(log_level)
    handler.setFormatter(formatter)
