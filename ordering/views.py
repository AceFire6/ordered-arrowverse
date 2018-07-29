from flask import render_template

from . import app
from .utils import get_full_series_episode_list


@app.route('/', methods=['GET'])
def index():
    context = {}

    episode_list = get_full_series_episode_list()

    context['table_content'] = episode_list
    context['show_list'] = app.config['SHOW_DICT']

    return render_template('index.html', **context)


@app.route('/newest_first/', methods=['GET'])
def index_newest_first():
    context = {}

    episode_list = get_full_series_episode_list()
    episode_list = episode_list[::-1]

    context['table_content'] = episode_list
    context['show_list'] = app.config['SHOW_DICT']

    return render_template('index.html', **context)


@app.route('/hide/<list:hide_list>/', methods=['GET'])
def index_with_hidden(hide_list):
    context = {}

    episode_list = get_full_series_episode_list(hide_list)

    context['hidden_list'] = hide_list
    context['table_content'] = episode_list
    context['show_list'] = app.config['SHOW_DICT']

    return render_template('index.html', **context)


@app.route('/hide/<list:hide_list>/newest_first/', methods=['GET'])
def index_with_hidden_newest_first(hide_list):
    context = {}

    episode_list = get_full_series_episode_list(hide_list)
    episode_list = episode_list[::-1]

    context['hidden_list'] = hide_list
    context['table_content'] = episode_list
    context['show_list'] = app.config['SHOW_DICT']

    return render_template('index.html', **context)
