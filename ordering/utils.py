from concurrent.futures import as_completed, ThreadPoolExecutor
from datetime import datetime

from dateutil.parser import parse as parse_date_string
from operator import itemgetter
import re

import requests
from bs4 import BeautifulSoup

from . import app
from .constants import (
    ARROW,
    BATWOMAN,
    BLACK_LIGHTNING,
    CONSTANTINE,
    FLASH,
    FREEDOM_FIGHTERS,
    LEGENDS_OF_TOMORROW,
    SUPERGIRL,
    VIXEN,
    WIKIPEDIA,
)

TWELVE_HOURS = 43200


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
            if series.upper() not in [CONSTANTINE, FREEDOM_FIGHTERS, BATWOMAN]:
                continue

        season += 1

        if not from_wikipedia:
            table = [
                row.strip().split('\n')
                for row in table.getText().split('\n\n') if row.strip()
            ]
        else:
            table_heading = table.find(name='tr', class_=None)
            table_headings = [
                heading.getText().replace(' ', '').lower()
                for heading in table_heading.children
            ]
            episode_num_index = table_headings.index('no.inseason')
            title_index = table_headings.index('title')
            air_date_index = table_headings.index('originalairdate')

            wikipedia_row_unpacker = itemgetter(episode_num_index, title_index, air_date_index)

            table = [
                [
                    episode_row_col.getText()
                    for episode_row_col in wikipedia_row_unpacker(episode_row.contents)
                ]
                for episode_row in table.find_all(class_='vevent')
            ]

        for row in table:
            # TODO: Make more robust - protects against rows that don't have enough data
            if len(row) < 2:
                continue

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
                row[-1] = air_date = parse_date_string(date).date()
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


def _swap_episode_rows(episode_list, index_1, index_2):
    episode_list[index_1], episode_list[index_2] = episode_list[index_2], episode_list[index_1]


def _handle_screening_day_error(episode_list):
    problem_episodes = (episode_list[78], episode_list[79])

    one_is_flash = any([x['series'].upper() == FLASH for x in problem_episodes])
    one_is_arrow = any([x['series'].upper() == ARROW for x in problem_episodes])

    both_are_episode_17 = all([x['episode_id'].endswith('E17') for x in problem_episodes])

    if one_is_arrow and one_is_flash and both_are_episode_17:
        _swap_episode_rows(episode_list, 78, 79)


def _handle_air_time_error(episode_list):
    # handles when two episodes air on the same day
    seasons = ['', 'Midseason', 'Midseason', 'Midseason', 'Midseason', 'Midseason', 'Midseason',
               'Summer', 'Summer', 'Summer', 'Fall', 'Fall', 'Fall']
    air_orders = {'Fall 2016': [(LEGENDS_OF_TOMORROW, VIXEN)],
                  'Midseason 2017': [(FLASH, LEGENDS_OF_TOMORROW)],
                  'Fall 2017': [(FLASH, LEGENDS_OF_TOMORROW)],
                  'Midseason 2018': [(FLASH, BLACK_LIGHTNING)],
                  'Fall 2018': [(ARROW, LEGENDS_OF_TOMORROW),
                                (FLASH, BLACK_LIGHTNING),
                                (SUPERGIRL, BLACK_LIGHTNING)],
                  'Midseason 2019': [(ARROW, BLACK_LIGHTNING),
                                     (LEGENDS_OF_TOMORROW, ARROW)],
                  'Fall 2019': [(BATWOMAN, SUPERGIRL),
                                (FLASH, ARROW)]}

    for i in range(len(episode_list)-1):
        curr_ep = episode_list[i]
        next_ep = episode_list[i+1]

        if curr_ep['air_date'] != next_ep['air_date']:
            continue

        air_date = curr_ep['air_date']
        air_season = f'{seasons[air_date.month]} {air_date.year}'

        if air_season not in air_orders.keys():  # handle summer releases
            continue

        pairs = air_orders[air_season]
        for series_1, series_2 in pairs:
            if curr_ep['series'].upper() == series_2 and next_ep['series'].upper() == series_1:
                _swap_episode_rows(episode_list, i, i+1)


def _handle_crisis_on_earth_x_order_error(episode_list, shows_in_list):
    earth_x_show_order = (SUPERGIRL, ARROW, FLASH, LEGENDS_OF_TOMORROW)
    earth_x_ordered_shows = [show for show in earth_x_show_order if show in shows_in_list]
    episode_indices = {show: None for show in earth_x_ordered_shows}

    for index, episode in enumerate(episode_list):
        show_name = episode['series'].upper()
        if show_name not in earth_x_show_order:
            continue

        episode_name = episode['episode_name']
        if not episode_name.startswith('Crisis on Earth-X, Part'):
            continue

        episode_indices[show_name] = index

    if ARROW in episode_indices and SUPERGIRL in episode_indices:
        arrow_index = episode_indices[ARROW]
        supergirl_index = episode_indices[SUPERGIRL]
        first_index, second_index = sorted([arrow_index, supergirl_index])

        episode_list[supergirl_index], episode_list[arrow_index] = episode_list[first_index], episode_list[second_index]  # noqa: E501

    if FLASH in episode_indices and LEGENDS_OF_TOMORROW in episode_indices:
        flash_index = episode_indices[FLASH]
        legends_index = episode_indices[LEGENDS_OF_TOMORROW]
        first_index, second_index = sorted([flash_index, legends_index])

        episode_list[flash_index], episode_list[legends_index] = episode_list[first_index], episode_list[second_index]  # noqa: E501


def _handle_john_con_noir_episode(episode_list):
    for index in range(len(episode_list)):
        show_name = episode_list[index]['series'].upper()
        if show_name != CONSTANTINE:
            continue

        episode_name = episode_list[index]['episode_name']
        if episode_name == 'Trash':
            episode_list.pop(index)
            break


def sort_and_filter_episodes(show_list_set, from_date=None, to_date=None):
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

    _handle_crisis_on_earth_x_order_error(full_list, shows_in_list)

    _handle_air_time_error(full_list)

    if CONSTANTINE in shows_in_list:
        _handle_john_con_noir_episode(full_list)

    filtered_list = _filter_on_air_date(full_list, from_date, to_date)

    count = 0
    for row in filtered_list:
        count += 1
        row['row_number'] = count
        row['air_date'] = f'{row["air_date"]:%B %d, %Y}'

    return filtered_list


def _filter_on_air_date(episode_list, from_date, to_date):
    if not from_date and not to_date:
        return episode_list

    if from_date:
        episode_list = [episode for episode in episode_list if episode['air_date'] >= from_date]

    if to_date:
        episode_list = [episode for episode in episode_list if episode['air_date'] <= to_date]

    return episode_list


@app.cache.memoize(timeout=TWELVE_HOURS)
def get_url_content(url):
    return requests.get(url).content


def get_show_list_from_show_html(show_name, show_html):
    show_soup = BeautifulSoup(show_html, 'html.parser')
    show_list = get_episode_list(show_soup, show_name)
    return show_list


@app.cache.memoize(timeout=TWELVE_HOURS)
def get_full_series_episode_list(excluded_series=None, from_date=None, to_date=None):
    excluded_series = [] if excluded_series is None else excluded_series
    shows_to_get = [show for show in app.config['SHOWS'] if show['id'] not in excluded_series]
    show_lists = []

    with ThreadPoolExecutor(max_workers=len(shows_to_get)) as executor:
        named_show_futures = {
            executor.submit(get_url_content, show['root'] + show['url']): show['name']
            for show in shows_to_get
        }

        for show_future in as_completed(named_show_futures):
            show_name = named_show_futures[show_future]
            show_html = show_future.result()

            show_list = get_show_list_from_show_html(show_name, show_html)
            show_lists.append(show_list)

    return sort_and_filter_episodes(show_lists, from_date=from_date, to_date=to_date)


def _get_bool(arg):
    if isinstance(arg, bool):
        return arg

    return arg == 'True'


def _get_date(arg):
    if isinstance(arg, datetime):
        return arg

    if arg is None:
        return None

    return parse_date_string(arg).date()
