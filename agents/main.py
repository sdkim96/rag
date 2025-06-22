import time
import asyncio
from typing import AsyncGenerator

import uvicorn
from fastapi import FastAPI
from fastapi.responses import StreamingResponse

from pydantic import BaseModel

app = FastAPI()

class CompletionData(BaseModel):
    status: int
    token: str
    error: bool
    done: bool
    sec: str
    message_res: str | None = None
    appendable: str | None = None


# 모의 completion 함수
async def mock_completion_stream(is_escape_n: bool) -> AsyncGenerator[str, None]:
    start_time = time.time()

    # 초기 메시지
    yield CompletionData(
        status=200,
        token="⚙️ 유저의 질문을 받고 있습니다..",
        error=False,
        done=False,
        sec="0",
        message_res=None,
        appendable=None
    ).model_dump_json() + ("\n" if is_escape_n else "")

    await asyncio.sleep(0.5)

    # 의도 파악 메시지
    yield CompletionData(
        status=200,
        token="🔍 유저의 질문 파악 완료: normal 질문입니다.",
        error=False,
        done=False,
        sec=str(int(time.time() - start_time)),
        message_res=None,
        appendable=None
    ).model_dump_json() + ("\n" if is_escape_n else "")
    await asyncio.sleep(0.5)

    # 스트리밍 응답
    for i, chunk in enumerate(["안녕하세요.", "무엇을 도와드릴까요?", "오늘도 좋은 하루 되세요!"]):
        yield CompletionData(
            status=200,
            token=chunk,
            error=False,
            done=False,
            sec=str(int(time.time() - start_time)),
            message_res=None,
            appendable=None
        ).model_dump_json() + ("\n" if is_escape_n else "")
        await asyncio.sleep(0.8)

    # 완료 메시지
    yield CompletionData(
        status=200,
        token="✅ 응답 완료",
        error=False,
        done=True,
        sec=str(int(time.time() - start_time)),
        message_res=None,
        appendable=None
    ).model_dump_json() + ("\n" if is_escape_n else "")


@app.post("/escape_n/stream")
async def stream():
    return StreamingResponse(
        mock_completion_stream(is_escape_n = True),
        media_type="application/json",
        headers={
            "Cache-Control": "no-cache",
            "Connection": "keep-alive",
            "X-Accel-Buffering": "no"  # Nginx 등 프록시에서 버퍼링 비활성화
        }
    )

@app.post("/escape_n_false/stream")
async def nstream():
    return StreamingResponse(
        mock_completion_stream(is_escape_n = False),
        media_type="application/json",
        headers={
            "Cache-Control": "no-cache",
            "Connection": "keep-alive",
            "X-Accel-Buffering": "no"  # Nginx 등 프록시에서 버퍼링 비활성화
        }
    )

if __name__ == "__main__":
    
    uvicorn.run(app, host="0.0.0.0", port=8000)