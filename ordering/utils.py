from datetime import datetime


def get_season(series_num, episode_num):
    return int(episode_num) / int(series_num)


def get_episode_list(series_soup, series):
    episode_list = []
    season = 0
    for table in series_soup.find_all('table')[1:]:
        season += 1
        table = [row.strip().split('\n')
                 for row in table.getText().split('\n\n') if row]

        for row in table:
            try:
                row[-1] = air_date = datetime.strptime(row[-1], '%B %d, %Y')
            except ValueError:
                continue

            if air_date and 'TBA' not in row:
                episode = 'S{:>02}E{:>02}'.format(season, row[1])
                row = [series, episode, row[2].replace('"', ''), row[3]]
                episode_list.append(row)
    return episode_list


def sort_episodes(list_1, list_2):
    full_list = []
    full_list.extend(list_1)
    full_list.extend(list_2)

    full_list = sorted(full_list, key=lambda episode_list: episode_list[-1])

    for i in range(79, len(full_list), 2):
        if i + 1 < len(full_list):
            full_list[i + 1], full_list[i] = full_list[i], full_list[i + 1]

    count = 0
    for row in full_list:
        count += 1
        row.insert(0, count)
        row[-1] = row[-1].strftime('%B %d, %Y')

    return full_list
