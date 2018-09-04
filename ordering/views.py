from datetime import date

from flask import render_template, request

from . import app
from .utils import _get_bool, _get_date, get_full_series_episode_list


@app.route('/', methods=['GET'])
def index():
    context = {}

    newest_first = request.args.get('newest_first', default=False, type=_get_bool)
    hide_shows_list = request.args.getlist('hide_show')
    after_date = request.args.get('after', default=date(2012, 9, 12), type=_get_date)
    before_date = request.args.get('before', default=date.today(), type=_get_date)

    episode_list = get_full_series_episode_list(
        excluded_series=hide_shows_list,
        after_date=after_date,
        before_date=before_date,
    )

    if newest_first:
        episode_list = episode_list[::-1]

    context['table_content'] = episode_list
    context['show_list'] = app.config['SHOW_DICT']
    context['hidden_show_list'] = hide_shows_list
    context['newest_first'] = newest_first

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

    context['hidden_show_list'] = hide_list
    context['table_content'] = episode_list
    context['show_list'] = app.config['SHOW_DICT']

    return render_template('index.html', **context)


@app.route('/hide/<list:hide_list>/newest_first/', methods=['GET'])
def index_with_hidden_newest_first(hide_list):
    context = {}

    episode_list = get_full_series_episode_list(hide_list)
    episode_list = episode_list[::-1]

    context['hidden_show_list'] = hide_list
    context['table_content'] = episode_list
    context['show_list'] = app.config['SHOW_DICT']

    return render_template('index.html', **context)
