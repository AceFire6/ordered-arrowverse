from flask import render_template

from . import app
from .utils import get_full_series_episode_list


@app.route('/', methods=['GET'])
@app.route('/<newest_first>', methods=['GET'])
def index(newest_first=None):
    context = {}

    episode_list = get_full_series_episode_list()

    if newest_first == 'newest_first':
        episode_list = episode_list[::-1]

    context['table_content'] = episode_list
    context['show_list'] = app.config['SHOW_DICT']

    return render_template('index.html', **context)


@app.route('/hide/<list:hide_list>/', methods=['GET'])
@app.route('/hide/<list:hide_list>/<newest_first>', methods=['GET'])
def index_with_hidden(hide_list, newest_first=None):
    context = {}

    episode_list = get_full_series_episode_list(hide_list)

    if newest_first == 'newest_first':
        episode_list = episode_list[::-1]

    context['hidden_list'] = hide_list
    context['table_content'] = episode_list
    context['show_list'] = app.config['SHOW_DICT']

    return render_template('index.html', **context)
