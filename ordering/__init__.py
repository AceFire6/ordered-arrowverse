from quart import Quart
from quart_compress import Compress
from quart_minify import Minify
# from tortoise.contrib.quart import register_tortoise

from .settings import DATABASE_URL, DEBUG
from .url_converters import ListConverter

app = Quart(__name__)

app.config.from_pyfile('settings.py')

# Minify HTML and any inline JS or CSS
Minify(app, js=True)

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
