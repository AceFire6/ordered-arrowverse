from flask import Flask
from flask.ext.cache import Cache

from ordering.url_converters import ListConverter

app = Flask(__name__)

app.config.from_pyfile('settings.py')

app.cache = Cache(app)

app.url_map.converters['list'] = ListConverter

from ordering import filters
from ordering import views
