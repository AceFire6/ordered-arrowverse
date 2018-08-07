from flask import Flask
from flask_caching import Cache
from flask_compress import Compress
from flask_minify import minify

from .url_converters import ListConverter

app = Flask(__name__)

app.config.from_pyfile('settings.py')

app.cache = Cache(app)

# gzip responses
Compress(app)
# Minify HTML and any inline JS or CSS
minify(app, js=True)

app.url_map.converters['list'] = ListConverter

from . import filters
from . import views
