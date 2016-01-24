import requests
from bs4 import BeautifulSoup
from flask import render_template

from ordering import app
from utils import get_episode_list, sort_episodes


@app.route('/', methods=['GET'])
@app.route('/<newest_first>', methods=['GET'])
@app.cache.memoize(timeout=43200)
def index(newest_first=None):
    context = {}

    show_list_set = []
    for show in app.config['SHOWS']:
        show_html = requests.get(show['url']).content
        show_list = get_episode_list(BeautifulSoup(show_html), show['name'])
        show_list_set.append(show_list)

    episode_list = sort_episodes(show_list_set)

    if newest_first == 'newest_first':
        episode_list = episode_list[::-1]

    context['newest_first'] = True if newest_first else False
    context['table_content'] = episode_list

    return render_template('index.html', **context)
