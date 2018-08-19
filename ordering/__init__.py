from flask import Flask
from flask_assets import Bundle, Environment
from flask_caching import Cache
from flask_compress import Compress
from flask_minify import minify

from .url_converters import ListConverter

app = Flask(__name__)

app.config.from_pyfile('settings.py')

app.cache = Cache(app)

js_assets = Bundle('js/cookie.js', 'js/index.js', filters='rjsmin', output='gen/bundled.js')
css_assets = Bundle('css/index.css', filters='cssmin', output='gen/bundled.css')

assets = Environment(app)
assets.register('js_all', js_assets)
assets.register('css_all', css_assets)

# gzip responses
Compress(app)
# Minify HTML and any inline JS or CSS
minify(app, js=True)

app.url_map.converters['list'] = ListConverter

from . import filters
from . import views
