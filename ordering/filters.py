from flask import request

from . import app
from .constants import WIKIPEDIA


def url_form(episode_name):
    return episode_name.replace(' ', '_')


@app.template_filter('episode_url')
def episode_url_filter(episode_name, series):
    root_url = app.config['SHOW_DICT_WITH_NAMES'][series]['root']
    from_wikipedia = WIKIPEDIA in root_url

    if episode_name == 'Pilot' or from_wikipedia:
        return root_url + url_form(episode_name + ' (%s)' % series)
    else:
        return root_url + url_form(episode_name)


@app.context_processor
def inject_oldest_first_url():
    if request.url.endswith('/newest_first'):
        return {'oldest_first_url': '/'.join(request.url.split('/')[:-1])}
    else:
        return {'oldest_first_url': None}


@app.context_processor
def inject_newest_first():
    return {'newest_first': request.url.endswith('/newest_first')}


@app.context_processor
def inject_show_dict():
    return {'series_map': app.config['SHOW_DICT_WITH_NAMES']}
