from flask import Flask
from flask.ext.cache import Cache

app = Flask(__name__)

app.config.from_pyfile('settings.py')

app.cache = Cache(app)

from ordering import filters
from ordering import views
