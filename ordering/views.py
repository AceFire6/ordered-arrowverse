from flask import render_template, request

from . import app
from .utils import _get_bool, _get_date, get_full_series_episode_list


@app.route('/', methods=['GET'])
def index():
    context = {}

    newest_first = request.args.get('newest_first', default=False, type=_get_bool)
    hide_shows_list = request.args.getlist('hide_show')
    from_date = request.args.get('from_date', default=None, type=_get_date)
    to_date = request.args.get('to_date', default=None, type=_get_date)

    episode_list = get_full_series_episode_list(
        excluded_series=hide_shows_list,
        from_date=from_date,
        to_date=to_date,
    )

    if newest_first:
        episode_list = episode_list[::-1]

    context['table_content'] = episode_list
    context['hidden_show_list'] = hide_shows_list
    context['from_date'] = from_date
    context['to_date'] = to_date

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
