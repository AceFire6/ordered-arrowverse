import re
import requests

from bs4 import BeautifulSoup
from datetime import datetime
from operator import itemgetter

from . import app
from .constants import ARROW, CONSTANTINE, FLASH, FREEDOM_FIGHTERS, SUPERGIRL, WIKIPEDIA


def get_episode_list(series_soup, series):
    episode_list = []
    season = 0
    from_wikipedia = WIKIPEDIA in app.config['SHOW_DICT_WITH_NAMES'][series]['root']

    if not from_wikipedia:
        tables = series_soup.find_all('table')
    else:
        tables = series_soup.find_all('table', class_='wikiepisodetable')

    for table in tables:
        table_name = table.getText().lower()
        if 'series overview' in table_name:
            continue

        if 'season' not in table_name:
            if series.upper() not in [CONSTANTINE, FREEDOM_FIGHTERS]:
                continue

        season += 1

        if not from_wikipedia:
            table = [
                row.strip().split('\n')
                for row in table.getText().split('\n\n') if row.strip()
            ]
        else:
            table = [
                [j.getText() for j in itemgetter(1, 3, 5, 11)(i.contents)]
                for i in table.find_all(class_='vevent')
            ]

        for row in table:
            if from_wikipedia:
                row[-1] = row[-1].split('(')[0].replace('\xa0', ' ').strip()

            episode_name = row[-2].replace('"', '')
            if '[' in episode_name:
                episode_name = episode_name.split('[')[0]

            episode_num = row[-3]
            try:
                date = row[-1]
                reference = re.search(r'\[\d+\]$', row[-1])
                date = date[:reference.start()] if reference else date
                row[-1] = air_date = datetime.strptime(date, '%B %d, %Y')
            except ValueError:
                continue

            if air_date and 'TBA' not in row:
                episode_id = f'S{season:>02}E{episode_num:>02}'
                episode_data = {
                    'series': series,
                    'episode_id': episode_id,
                    'episode_name': episode_name,
                    'air_date': air_date,
                }
                episode_list.append(episode_data)

    return episode_list


def _swap_episode_rows(epsidoe_list, index_1, index_2):
    epsidoe_list[index_1], epsidoe_list[index_2] = epsidoe_list[index_2], epsidoe_list[index_1]


def _handle_screening_day_error(episode_list):
    problem_episodes = (episode_list[78], episode_list[79])

    one_is_flash = any([x['series'].upper() == FLASH for x in problem_episodes])
    one_is_arrow = any([x['series'].upper() == ARROW for x in problem_episodes])

    both_are_episode_17 = all([x['episode_id'].endswith('E17') for x in problem_episodes])

    if one_is_arrow and one_is_flash and both_are_episode_17:
        _swap_episode_rows(episode_list, 78, 79)


def _handle_crisis_on_earth_x_order_error(episode_list):
    arrow_episode_index = None
    supergirl_episode_index = None

    for index in range(len(episode_list)):
        show_name = episode_list[index]['series'].upper()
        if show_name not in [ARROW, SUPERGIRL]:
            continue

        episode_name = episode_list[index]['episode_name']
        if not episode_name.startswith('Crisis on Earth-X, Part'):
            continue

        if show_name == ARROW:
            arrow_episode_index = index
        elif show_name == SUPERGIRL:
            supergirl_episode_index = index

    _swap_episode_rows(episode_list, arrow_episode_index, supergirl_episode_index)


def _handle_john_con_noir_episode(episode_list):
    for index in range(len(episode_list)):
        show_name = episode_list[index]['series'].upper()
        if show_name != CONSTANTINE:
            continue

        episode_name = episode_list[index]['episode_name']
        if episode_name == 'Trash':
            episode_list.pop(index)
            break


def sort_episodes(show_list_set):
    full_list = []
    shows_in_list = []

    for show_list in show_list_set:
        full_list.extend(show_list)

        show_name = show_list[0]['series'].upper()
        shows_in_list.append(show_name)

    full_list = sorted(full_list, key=lambda episode: episode['air_date'])

    # Fix screening time error caused by network
    # This fix corrects all the list errors.
    if len(full_list) > 80 and FLASH in shows_in_list and ARROW in shows_in_list:
        _handle_screening_day_error(full_list)

    if ARROW in shows_in_list and SUPERGIRL in shows_in_list:
        _handle_crisis_on_earth_x_order_error(full_list)

    if CONSTANTINE in shows_in_list:
        _handle_john_con_noir_episode(full_list)

    count = 0
    for row in full_list:
        count += 1
        row['row_number'] = count
        row['air_date'] = f'{row["air_date"]:%B %d, %Y}'

    return full_list


@app.cache.memoize(timeout=43200)
def get_url_content(url):
    return requests.get(url).content


def get_full_series_episode_list(excluded_series=list()):
    show_list_set = []
    for show in app.config['SHOWS']:
        if show['id'] not in excluded_series:
            show_html = get_url_content(show['root'] + show['url'])
            show_list = get_episode_list(
                BeautifulSoup(show_html, 'html.parser'),
                show['name']
            )
            show_list_set.append(show_list)

    return sort_episodes(show_list_set)
