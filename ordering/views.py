from datetime import datetime

from flask import render_template, request, url_for
from werkzeug.contrib.atom import AtomFeed

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


@app.route('/recent_episodes.atom')
def recent_episodes():
    feed = AtomFeed(
        title='Recent Episodes',
        feed_url=request.url,
        url=request.url_root,
        logo=url_for('static', filename='favicon.png', _external=True),
        icon=url_for('static', filename='favicon.png', _external=True),
    )

    hide_shows_list = request.args.getlist('hide_show')

    newest_first_episode_list = get_full_series_episode_list(excluded_series=hide_shows_list)[::-1]

    for episode in newest_first_episode_list[:15]:
        title = '{series} - {episode_id} - {episode_name}'.format(**episode)
        content = '{series} {episode_id} {episode_name} will air on {air_date}'.format(**episode)
        show_dict = app.config['SHOW_DICT_WITH_NAMES'][episode['series']]
        data_source = f"{show_dict['root']}{show_dict['url']}"

        feed.add(
            title=title,
            content=content,
            content_type='text',
            url=data_source,
            author=show_dict['root'],
            updated=datetime.now(),
        )

    return feed.get_response()


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
