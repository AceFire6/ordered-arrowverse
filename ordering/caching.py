from functools import wraps
import json
from logging import getLogger

from redis import Redis

from ordering.settings import REDIS_URL


cache = Redis.from_url(REDIS_URL)
logger = getLogger(__name__)


def safe_cache_content(timeout=None):
    def decorator(func):
        @wraps(func)
        def wrapper(*args, **kwargs):
            hashed_json = hash(f'{json.dumps(args)}-{json.dumps(kwargs)})')
            cache_key = f'{func.__name__}-{hashed_json}'

            cached_value = cache.get(cache_key)
            if cached_value:
                return json.loads(cached_value.decode('utf-8'))

            try:
                response = func(*args, **kwargs)
            except Exception as exception:
                logger.exception(f'Failed to run {func} - using cached LTS response')

                cached_json_lts_response = cache.get(f'{cache_key}-lts')
                if cached_json_lts_response is None:
                    logger.error('No valid cached LTS response - returning error response')
                    raise exception

                return json.loads(cached_json_lts_response.decode('utf-8'))
            else:
                json_response = json.dumps(response)
                cache.set(cache_key, json_response, ex=timeout)
                cache.set(f'{cache_key}-lts', json_response)

                return response

        return wrapper

    return decorator
