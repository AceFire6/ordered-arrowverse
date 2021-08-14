from functools import wraps
from logging import getLogger

import orjson
from redis import Redis

from ordering.settings import REDIS_URL


cache = Redis.from_url(REDIS_URL)
logger = getLogger(__name__)


def serialized_response(response):
    def handle_bytes(value):
        if isinstance(value, bytes):
            return value.decode('utf-8')

        raise TypeError

    return orjson.dumps(response, default=handle_bytes)


def safe_cache_content(timeout=None, backup=False, hash_args=False):
    def decorator(func):
        @wraps(func)
        def wrapper(*args, **kwargs):
            inputs = f'{args}-{kwargs}'
            if hash_args:
                inputs = hash(inputs)

            cache_key = f'{func.__name__}-{inputs}'

            cached_value = cache.get(cache_key)
            if cached_value:
                return orjson.loads(cached_value.decode('utf-8'))

            try:
                response = func(*args, **kwargs)
            except Exception as exception:
                logger.exception(f'Failed to run {func} - using cached backup response')

                cached_json_lts_response = cache.get(f'{cache_key}-backup')
                if cached_json_lts_response is None:
                    logger.error('No valid cached backup response - returning error response')
                    raise exception

                return orjson.loads(cached_json_lts_response.decode('utf-8'))
            else:
                json_response = serialized_response(response)
                cache.set(cache_key, json_response, ex=timeout)

                if backup:
                    cache.set(f'{cache_key}-backup', json_response)

                return response

        return wrapper

    return decorator
