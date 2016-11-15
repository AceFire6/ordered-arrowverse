import re
import requests

from bs4 import BeautifulSoup
from datetime import datetime
from operator import itemgetter

from ordering import app


def get_episode_list(series_soup, series):
    episode_list = []
    season = 0
    wikipedia = 'wikipedia' in app.config['SHOW_DICT'][series]['root']

    if not wikipedia:
        tables = series_soup.find_all('table')
    else:
        tables = series_soup.find_all('table', class_='wikiepisodetable')

    for table in tables:
        if 'series overview' in table.getText().lower():
            continue
        season += 1
        if not wikipedia:
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
            if wikipedia:
                row[-1] = row[-1].split('(')[0].replace(u'\xa0', ' ').strip()
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
                episode_id = 'S{:>02}E{:>02}'.format(season, episode_num)
                row = [series, episode_id, episode_name, air_date]
                episode_list.append(row)
    return episode_list


def sort_episodes(show_list_set):
    full_list = []
    for show_list in show_list_set:
        full_list.extend(show_list)

    full_list = sorted(full_list, key=lambda episode_list: episode_list[-1])

    # Fix screening time error caused by network
    # This fix corrects all the list errors.
    ep_17 = all(
        map(lambda x: x[1].endswith('E17'), (full_list[78], full_list[79]))
    )
    if len(full_list) > 80 and ep_17:
        full_list[78], full_list[79] = full_list[79], full_list[78]

    count = 0
    for row in full_list:
        count += 1
        row.insert(0, count)
        row[-1] = row[-1].strftime('%B %d, %Y')

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
                    BeautifulSoup(show_html), show['name']
            )
            show_list_set.append(show_list)

    return sort_episodes(show_list_set)
