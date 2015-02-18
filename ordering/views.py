import requests
from bs4 import BeautifulSoup
from flask import render_template

from . import app
from utils import get_episode_list, sort_episodes


@app.route('/', methods=['GET'])
def index():
    context = {}

    arrow_html = requests.get(
        'http://arrow.wikia.com/wiki/List_of_Arrow_episodes').content
    flash_html = requests.get(
        'http://arrow.wikia.com/wiki/List_of_The_Flash_episodes').content

    arrow_list = get_episode_list(BeautifulSoup(arrow_html), 'Arrow')
    flash_list = get_episode_list(BeautifulSoup(flash_html), 'The Flash')

    full_list = sort_episodes(arrow_list, flash_list)

    context['table_content'] = full_list

    return render_template('index.html', **context)
