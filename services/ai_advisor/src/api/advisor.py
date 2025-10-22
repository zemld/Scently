import json
import os
from typing import Annotated, Any

from fastapi import APIRouter, HTTPException, Query
from openai import OpenAI
from openai.types.chat import ChatCompletionMessageParam

from src.models.perfume import Perfume, PerfumeOut

router = APIRouter(prefix="/v1/advise")

api_key = os.getenv("API_KEY")
folder_id = os.getenv("FOLDER_ID")
base_url = "https://llm.api.cloud.yandex.net/v1"

system_prompt = """You are an expert perfume recommender. 
             Your task is to suggest alternative perfumes 
             based on the user's favorite one. 
             Do not include explanations, text, or markdown — only JSON."""


def create_model_url(folder_id: str) -> str:
    return f"gpt://{folder_id}/gpt-oss-20b/latest"


def create_user_prompt(brand: str, name: str) -> str:
    return f"""User's favorite perfume:
            Brand: {brand}
            Name: {name}

            Return exactly 6 other perfumes that the user might also like.
            Each item must include:
            - brand
            - name

            Respond strictly in this JSON format:
            [
            {{"brand": "string", "name": "string"}},
            {{"brand": "string", "name": "string"}},
            {{"brand": "string", "name": "string"}},
            {{"brand": "string", "name": "string"}},
            {{"brand": "string", "name": "string"}},
            {{"brand": "string", "name": "string"}}
            ]"""


def get_messages(perfume: Perfume) -> list[ChatCompletionMessageParam]:
    return [
        {
            "role": "system",
            "content": system_prompt,
        },
        {
            "role": "user",
            "content": create_user_prompt(perfume.brand, perfume.name),
        },
    ]


def get_response_format() -> dict[str, Any]:
    return {
        "type": "json_schema",
        "json_schema": {
            "name": "PerfumeList",
            "schema": {
                "type": "array",
                "items": {
                    "type": "object",
                    "properties": {
                        "brand": {"type": "string"},
                        "name": {"type": "string"},
                    },
                    "required": ["brand", "name"],
                    "additionalProperties": False,
                },
            },
        },
    }


@router.get(
    "/",
    response_model=PerfumeOut,
    tags=["Perfume Advisor"],
    summary="""Get advise for a perfume which contains other perfumes
    that the user might also like""",
    description="""Get advise for a perfume which contains other perfumes
    that the user might also like""",
    responses={
        200: {
            "description": "Success",
            "model": PerfumeOut,
        },
        400: {
            "description": "Bad Request - invalid input parameters",
        },
        500: {
            "description": "Internal Server Error",
        },
    },
)
async def get_advise(
    brand: Annotated[
        str,
        Query(
            min_length=1, max_length=128, description="Perfume brand", example="Chanel"
        ),
    ],
    name: Annotated[
        str,
        Query(
            min_length=1, max_length=128, description="Perfume name", example="Chance"
        ),
    ],
) -> PerfumeOut:
    if not brand.strip() or not name.strip():
        raise HTTPException(
            status_code=400,
            detail="Brand and name cannot be empty or contain only whitespace",
        )

    try:
        perfume = Perfume(brand=brand.strip(), name=name.strip())

        client = OpenAI(
            api_key=api_key,
            base_url=base_url,
        )

        response = client.chat.completions.create(
            model=create_model_url(folder_id),
            messages=get_messages(perfume),
            response_format=get_response_format(),
            max_tokens=500,
            temperature=0.4,
            stream=False,
        )

        json_response = json.loads(response.choices[0].message.content)
        perfumes = [Perfume(**item) for item in json_response]

        return PerfumeOut(perfumes=perfumes)

    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Ошибка API: {str(e)}")
