from datetime import datetime


def get_episode_list(series_soup, series):
    episode_list = []
    season = 0
    for table in series_soup.find_all('table'):
        if 'series overview' in table.getText().lower():
            continue
        season += 1
        table = [row.strip().split('\n')
                 for row in table.getText().split('\n\n') if row]

        for row in table:
            episode_name = row[-2].replace('"', '')
            episode_num = row[-3]
            try:
                row[-1] = air_date = datetime.strptime(row[-1], '%B %d, %Y')
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
    full_list[78], full_list[79] = full_list[79], full_list[78]

    count = 0
    for row in full_list:
        count += 1
        row.insert(0, count)
        row[-1] = row[-1].strftime('%B %d, %Y')

    return full_list
