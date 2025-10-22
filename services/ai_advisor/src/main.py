from fastapi import FastAPI

from src.api.advisor import router

advisor = FastAPI()
advisor.include_router(router)

advisor.openapi_tags = [
    {
        "name": "Perfume Advisor",
        "description": "AI Perfume Advisor",
    }
]
