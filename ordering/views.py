from feedgen.entry import FeedEntry
from feedgen.feed import FeedGenerator
from quart import jsonify, make_response, render_template, request, url_for

from . import app
from .utils import _get_bool, _get_date, get_full_series_episode_list


@app.route('/', methods=['GET'])
async def index():
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

    return await render_template('index.html', **context)


@app.route('/api', methods=['GET'])
async def api():
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

    return jsonify(episode_list)


@app.route('/recent_episodes.atom')
async def recent_episodes():
    logo_link = url_for('static', filename='favicon.png', _external=True)

    feed = FeedGenerator()
    feed.title('Arrowverse.info - Recent Episodes')
    feed.id(request.url_root)
    feed.link(href=request.url)
    feed.logo(logo_link)
    feed.icon(logo_link)
    feed.language('en')

    hide_shows_list = request.args.getlist('hide_show')

    newest_first_episode_list = get_full_series_episode_list(excluded_series=hide_shows_list)[::-1]

    for episode in newest_first_episode_list[:15]:
        title = '{series} - {episode_id} - {episode_name}'.format(**episode)
        content = '{series} {episode_id} {episode_name} will air on {air_date}'.format(**episode)
        show_dict = app.config['SHOW_DICT_WITH_NAMES'][episode['series']]
        data_source = f"{show_dict['root']}{show_dict['url']}"

        feed_entry = FeedEntry()
        feed_entry.id(data_source)
        feed_entry.link({'href': data_source})
        feed_entry.title(title)
        feed_entry.content(content, type='text')
        feed_entry.author(uri=show_dict['root'])

        feed.add_entry(feed_entry)

    response = await make_response(feed.atom_str(pretty=True))
    response.headers['Content-Type'] = 'application/atom+xml'

    return response


@app.route('/newest_first/', methods=['GET'])
async def index_newest_first():
    context = {}

    episode_list = get_full_series_episode_list()
    episode_list = episode_list[::-1]

    context['table_content'] = episode_list
    context['show_list'] = app.config['SHOW_DICT']

    return await render_template('index.html', **context)


@app.route('/hide/<list:hide_list>/', methods=['GET'])
async def index_with_hidden(hide_list):
    context = {}

    episode_list = get_full_series_episode_list(hide_list)

    context['hidden_show_list'] = hide_list
    context['table_content'] = episode_list
    context['show_list'] = app.config['SHOW_DICT']

    return await render_template('index.html', **context)


@app.route('/hide/<list:hide_list>/newest_first/', methods=['GET'])
async def index_with_hidden_newest_first(hide_list):
    context = {}

    episode_list = get_full_series_episode_list(hide_list)
    episode_list = episode_list[::-1]

    context['hidden_show_list'] = hide_list
    context['table_content'] = episode_list
    context['show_list'] = app.config['SHOW_DICT']

    return await render_template('index.html', **context)
