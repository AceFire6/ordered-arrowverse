from quart import Quart
from quart_compress import Compress
# from tortoise.contrib.quart import register_tortoise

from .settings import DATABASE_URL, DEBUG
from .url_converters import ListConverter

app = Quart(__name__)

app.config.from_pyfile('settings.py')

# gzip responses
Compress(app)

app.url_map.converters['list'] = ListConverter

from . import filters
from . import views

# register_tortoise(
#     app,
#     db_url=DATABASE_URL,
#     modules={'models': ['ordering.models']},
#     generate_schemas=DEBUG,
# )
