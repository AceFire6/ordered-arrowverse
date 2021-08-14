from flask_assets import Environment
from quart import Quart
from quart_compress import Compress
from quart_minify import Minify
# from tortoise.contrib.quart import register_tortoise
from webassets import Bundle

from .settings import DATABASE_URL, DEBUG
from .url_converters import ListConverter

app = Quart(__name__)

app.config.from_pyfile('settings.py')

js_assets = Bundle('js/cookie.js', 'js/index.js', filters='rjsmin', output='gen/bundled.js')
css_assets = Bundle('css/index.css', filters='cssmin', output='gen/bundled.css')

assets = Environment(app)
assets.register('js_all', js_assets)
assets.register('css_all', css_assets)

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
