from flask import Flask
from flask_caching import Cache

from .url_converters import ListConverter

app = Flask(__name__)

app.config.from_pyfile('settings.py')

app.cache = Cache(app)

app.url_map.converters['list'] = ListConverter

from . import filters
from . import views
