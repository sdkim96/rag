import httpx
import logging
from typing import (
    Tuple, 
    Union, 
    Any, 
    TypeVar, 
    Literal, 
    Optional, 
    Generator
)

from pydantic import BaseModel

_ResponseT = Union[dict, BaseModel, str]
_T = TypeVar("_T", bound=_ResponseT)

LOGGER = logging.getLogger(__name__)

def get(
    cast_to: type[_T] = dict,
    *,
    client: httpx.Client | None,
    url: str,
    headers: dict[str, str] | None = None,
    params: dict[str, Any] | None = None,
    timeout: int | None = None,
    debug: bool = False
) -> Tuple[Optional[_T], Optional[str]]:
    """
    동기 GET 요청을 수행합니다.

    Args:
        cast_to: type[_T] - 응답을 변환할 타입 (기본값: dict)
        client: httpx.Client - 사용할 HTTP 클라이언트입니다.
        url: str - 요청할 URL입니다.
        headers: dict[str, str] | None - 요청 헤더입니다.
        params: dict[str, Any] | None - 요청할 파라미터입니다.
        timeout: int | None - 요청 타임아웃 (초 단위)입니다.
        debug: bool - 디버그 모드 (기본값: False)

    Returns:
        Tuple[Optional[_T], Optional[str]] - 변환된 응답 데이터와 에러 메시지 (에러가 없으면 None)

    """

    internal_client_flag = False
    if client is None:
        internal_client_flag = True
        client = httpx.Client()
        LOGGER.info("내부에서 http 클라이언트 생성")

    req = _build_request(
        client=client,
        method="GET",
        url=url,
        params=params,
        timeout=timeout,
        headers=headers,
        debug=debug
    )
    try:
        resp = client.send(req)
    except httpx.RequestError as e:
        return None, str(e)
    finally:
        if internal_client_flag:
            client.close()
            LOGGER.info("내부 http 클라이언트 종료")
    return _handle_response(resp, cast_to, debug)


def post(
    cast_to: type[_T] = dict,
    *,
    client: httpx.Client | None,
    url: str,
    headers: dict[str, str] | None = None,
    json_data: dict[str, Any] | None = None,
    timeout: int | None = None,
    debug: bool = False
) -> Tuple[Optional[_T], Optional[str]]:
    """
    동기 POST 요청을 수행합니다.

    매개변수:
        client: httpx.Client - 사용할 HTTP 클라이언트입니다.
        cast_to: type[_T] - 응답을 변환할 타입입니다.
        url: str - 요청할 URL입니다.
        headers: dict[str, str] | None - 요청 헤더입니다.
        json_data: dict[str, Any] | None - 요청할 JSON 데이터입니다.
        timeout: int | None - 요청 타임아웃 (초 단위)입니다.
        debug: bool - 디버그 모드 (기본값: False)

    반환값:
        변환된 응답 데이터와 에러 메시지 (에러가 없으면 None)
    """
    
    internal_client_flag = False
    if client is None:
        internal_client_flag = True
        client = httpx.Client()
        LOGGER.debug("내부에서 http 클라이언트 생성")

    req = _build_request(
        client=client,
        method="POST",
        url=url,
        json_data=json_data,
        timeout=timeout,
        headers=headers,
        debug=debug
    )
    LOGGER.info("요청객체 생성 완료")
    try:
        LOGGER.info("요청객체 전송")
        resp = client.send(req)
    except httpx.RequestError as e:
        LOGGER.error(f"HTTP 오류 발생: {str(e)}")
        return None, str(e)
    finally:
        if internal_client_flag:
            client.close()
            LOGGER.debug("내부 http 클라이언트 종료")

    LOGGER.info("응답객체 수신 완료")
    return _handle_response(resp, cast_to, debug)


def stream(
    stream_to: Literal['file', 'sse', 'bytes'] = 'bytes',
    *,
    client: httpx.Client | None,
    url: str,
    method: Literal["GET", "POST"] = "GET",
    headers: dict[str, str] | None = None,
    params: dict[str, Any] | None = None,
    json_data: dict[str, Any] | None = None,
    timeout: int | None = None,
    debug: bool = False
) -> Generator[Any, None, None]:
    """
    다양한 유형의 스트리밍 응답을 지원하는 동기 함수입니다.

    Args:
        client: HTTPX 동기 클라이언트
        stream_to: 스트리밍할 데이터의 타입 (기본값: bytes)
        url: 요청할 URL
        method: HTTP 메서드 (GET 또는 POST)
        headers: 요청 헤더 (기본값: None)
        params: GET 요청 시 사용할 파라미터 (기본값: None)
        json_data: POST 요청 시 사용할 JSON 데이터 (기본값: None)
        timeout: 요청 타임아웃 (초 단위) (기본값: None)
        debug: 디버그 모드 (기본값: False)

    Returns:
        응답 generator 또는 None
    """
    LOGGER.info("스트리밍 요청 시작, 스트리밍 요청은 로깅이 디버그 모드일때 CURL출력을 지원하지 않습니다.")
    error = "ERROR: "
    if stream_to not in ['file', 'sse', 'bytes']:
        error += "stream_to는 'file', 'sse', 'bytes' 중 하나여야 합니다."
        LOGGER.error(error)
        yield error

    internal_client_flag = False
    if client is None:
        internal_client_flag = True
        client = httpx.Client()
        LOGGER.info("내부에서 http 클라이언트 생성")
    
    request_args: dict[str, Any] = {
        "method": method,
        "url": url,
        "timeout": timeout,
        "headers": headers,
    }
    match method:
        case "POST":
            request_args["json"] = json_data
        case "GET":
            request_args["params"] = params
        case _:
            error += "유효하지 않은 HTTP 메서드입니다."
            LOGGER.error(error)
            yield error
    LOGGER.info(f"스트리밍 요청객체 생성 완료")
    LOGGER.info("스트리밍 요청객체 전송")
    try:
        with client.stream(**request_args) as resp:
            resp.raise_for_status()
            LOGGER.info("스트리밍 응답객체 수신 완료. 스트리밍 시작.")

            if stream_to == 'file':
                for chunk in resp.iter_bytes():
                    yield chunk
            elif stream_to == 'sse':
                for line in resp.iter_lines():
                    yield line 
            elif stream_to == 'bytes':
                for byte in resp.iter_bytes():
                    yield byte

    except httpx.HTTPError as e:
        error += str(e)
        LOGGER.error(f"HTTP 오류 발생: {error}")
        yield error
    except Exception as e:
        error += str(e)
        LOGGER.error(f"HTTP 오류 발생: {error}")
        yield error
    finally:
        LOGGER.info("스트리밍 요청 종료")
        if internal_client_flag:
            client.close()
            LOGGER.info("내부 http 클라이언트 종료")

def delete(
    cast_to: type[_T] = dict,
    *,
    client: httpx.Client | None,
    url: str,
    headers: dict[str, str] | None = None,
    params: dict[str, Any] | None = None,
    json_data: dict[str, Any] | None = None,
    timeout: int | None = None,
    debug: bool = False
) -> Tuple[Optional[_T], Optional[str]]:
    """
    동기 DELETE 요청을 수행합니다.

    Args:
        client: httpx.Client - 사용할 HTTP 클라이언트입니다.
        url: str - 요청할 URL입니다.
        headers: dict[str, str] | None - 요청 헤더입니다.
        params: dict[str, Any] | None - 요청할 파라미터입니다.
        json_data: dict[str, Any] | None - 요청할 JSON 데이터입니다.
        timeout: int | None - 요청 타임아웃 (초 단위)입니다.
        debug: bool - 디버그 모드 (기본값: False)

    Returns:
        Tuple[Optional[bool], Optional[str]] - 성공 여부와 에러 메시지 (에러가 없으면 None)
    """
    
    internal_client_flag = False
    if client is None:
        internal_client_flag = True
        client = httpx.Client()
        LOGGER.info("내부에서 http 클라이언트 생성")

    req = _build_request(
        client=client,
        method="DELETE",
        url=url,
        params=params,
        json_data=json_data,
        timeout=timeout,
        headers=headers,
        debug=debug
    )
    
    try:
        resp = client.send(req)
    except httpx.RequestError as e:
        return None, str(e)
    finally:
        if internal_client_flag:
            client.close()
            LOGGER.info("내부 http 클라이언트 종료")
    
    if resp.is_error:
        return None, f"HTTP 오류 발생: {resp.status_code} {resp.text}"
    
    return _handle_response(resp, cast_to, debug)


async def aget(
    cast_to: type[_T] | None = None,
    *,
    client: httpx.AsyncClient | None,
    url: str,
    headers: dict[str, str] | None = None,
    params: dict[str, Any] | None = None,
    timeout: int | None = None,
    debug: bool = False
) -> Tuple[Optional[_T], Optional[str]]:
    """
    비동기 GET 요청을 수행합니다.

    매개변수:
        client: httpx.AsyncClient - 사용할 비동기 HTTP 클라이언트입니다.
        cast_to: Optional[type[_T]] - 응답을 변환할 타입
        opts: Opts[GET] - GET 요청을 위한 옵션 객체
        url: str - 요청할 URL입니다.
        headers: dict[str, str] | None - 요청 헤더입니다.
        params: dict[str, Any] | None - 요청할 파라미터입니다.
        timeout: int | None - 요청 타임아웃 (초 단위)입니다.
        debug: bool - 디버그 모드 (기본값: False)

    반환값:
        변환된 응답 데이터 또는 None, 에러 메시지 또는 None
    """
    internal_client_flag = False
    if client is None:
        internal_client_flag = True
        client = httpx.AsyncClient()
        LOGGER.info("내부에서 http 클라이언트 생성")

    req = _build_request(
        client=client,
        method="GET",
        url=url,
        params=params,
        timeout=timeout,
        headers=headers,
        debug=debug
    )
    try:
        resp = await client.send(req)
    except httpx.RequestError as e:
        return None, str(e)
    finally:
        if internal_client_flag:
            await client.aclose()
            LOGGER.info("내부 http 클라이언트 종료")
    
    return await _handle_async_response(resp, cast_to, debug)


async def apost(
    cast_to: type[_T] | None = None,
    *,
    client: httpx.AsyncClient | None,
    url: str,
    headers: dict[str, str] | None = None,
    json_data: dict[str, Any] | None = None,
    timeout: int | None = None,
    debug: bool = False
) -> Tuple[Optional[_T], Optional[str]]:
    """
    비동기 POST 요청을 수행합니다.

    매개변수:
        cast_to: Optional[type[_T]] - 응답을 변환할 타입
        client: httpx.AsyncClient - 사용할 비동기 HTTP 클라이언트입니다.
        url: str - 요청할 URL입니다.
        headers: dict[str, str] | None - 요청 헤더입니다.
        json_data: dict[str, Any] | None - 요청할 JSON 데이터입니다.
        timeout: int | None - 요청 타임아웃 (초 단위)입니다.
        debug: bool - 디버그 모드 (기본값: False)

    반환값:
        변환된 응답 데이터 또는 None, 에러 메시지 또는 None
    """
    internal_client_flag = False
    if client is None:
        internal_client_flag = True
        client = httpx.AsyncClient()
        LOGGER.info("내부에서 http 클라이언트 생성")

    req = _build_request(
        client=client,
        method="POST",
        url=url,
        json_data=json_data,
        timeout=timeout,
        headers=headers,
        debug=debug
    )
    try:
        resp = await client.send(req)
    except httpx.RequestError as e:
        return None, str(e)
    finally:
        if internal_client_flag:
            await client.aclose()
            LOGGER.info("내부 http 클라이언트 종료")
    
    return await _handle_async_response(resp, cast_to, debug)


def _convert_to_curl(
    request: httpx.Request,
):
    """
    httpx.Request 객체를 curl 명령어로 변환합니다.

    매개변수:
        request: httpx.Request - httpx 모듈의 요청 객체입니다.

    반환값:
        curl 명령어 문자열
    """
    curl_command = f"curl -X {request.method} '{request.url}'"
    
    for key, value in request.headers.items():
        curl_command += f" -H '{key}: {value}'"
    
    content = request.content.decode("utf-8")
    if content:
        curl_command += f" -d '{content}'"
    
    return curl_command


def _handle_response(
    response: httpx.Response,
    cast_to: type[_T],
    debug: bool,
) -> Tuple[Optional[_T], Optional[str]]:
    """
    동기 HTTP 응답을 처리합니다.

    매개변수:
        response: httpx.Response - httpx 모듈의 응답 객체입니다.
        cast_to: type[_T] - 응답 본문을 변환할 타입입니다.

    반환값:
        변환된 응답 데이터와 에러 메시지 (에러가 없으면 None)
    """
    try:
        response.raise_for_status()
    except httpx.HTTPStatusError as e:
        LOGGER.error(f"HTTP 오류 발생: {str(e)}")
        error_message = str(e) + ", " + str(response.text)
        return None, error_message

    try:
        if issubclass(cast_to, BaseModel):
            json_data= response.json()
            return cast_to.model_validate(json_data), None
        elif cast_to is str:
            return response.text, None  # type: ignore
        elif cast_to is dict:
            return response.json(), None  # type: ignore
        else:
            return cast_to(response.json()), None
    except Exception as e:
        LOGGER.error(f"응답 변환 오류: {str(e)}")
        return None, str(e)



def _build_request(
    client: httpx.Client | httpx.AsyncClient,
    method: Literal["GET", "POST", "DELETE"],
    url: str,
    debug: bool,
    params: dict[str, Any] | None = None,
    json_data: dict[str, Any] | None = None,
    timeout: int | None = None,
    headers: dict[str, str] | None = None,
    stream: bool = False,
) -> httpx.Request:    
    """
    요청을 빌드합니다.
    """
    match method:
        case "POST":
            request = client.build_request(
                method,
                url,
                json=json_data,
                params=params,
                headers=headers,
                timeout=timeout
            )
        case "GET":
            request = client.build_request(
                method,
                url,
                params=params,
                headers=headers,
                timeout=timeout
            )
        case "DELETE":
            request = client.build_request(
                method,
                url,
                params=params,
                headers=headers,
                timeout=timeout,
                json=json_data
            )

    curl = _convert_to_curl(request)
    if len(curl) > 1000:
        curl = curl[:1000] + "..."
    LOGGER.debug(f"CURL 명령어: {curl}")
    return request


async def _handle_async_response(
    response: httpx.Response,
    cast_to: type[_T] | None,
    debug: bool,
) -> Tuple[Optional[_T], Optional[str]]:
    """
    비동기 HTTP 응답을 처리합니다.

    매개변수:
        response: httpx.Response - 응답 객체
        cast_to: Optional[type[_T]] - 변환할 타입 (없으면 기본 파싱만 수행)

    반환값:
        변환된 응답 데이터 또는 None, 에러 메시지 또는 None
    """
    try:
        response.raise_for_status()
    except httpx.HTTPStatusError as e:
        LOGGER.error("HTTP 오류 발생: %s", str(e))
        return None, str(e)

    try:
        if cast_to is None:
            return response.json(), None
        if issubclass(cast_to, BaseModel):
            return cast_to(**response.json()), None
        elif cast_to is str:
            return response.text, None  # type: ignore
        else:
            return cast_to(response.json()), None
    except Exception as e:
        LOGGER.error("HTTP 오류 발생: %s", str(e))
        return None, str(e)