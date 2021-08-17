from quart import request, url_for

from . import app, settings
from .constants import WIKIPEDIA
from .settings import STATIC_CACHE_ID
from .utils import _get_bool


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


@app.template_global('static_url')
def static_url(filename):
    static_file_url = url_for('static', filename=filename)

    return f'{static_file_url}?{STATIC_CACHE_ID}'


@app.context_processor
def inject_oldest_first_url():
    if request.url.endswith('/newest_first/'):
        return {'oldest_first_url': request.url.replace('/newest_first/', '/')}
    else:
        return {'oldest_first_url': None}


@app.context_processor
def inject_newest_first():
    newest_first = request.args.get('newest_first', default=False, type=_get_bool)
    return {'newest_first': newest_first}


@app.context_processor
def inject_using_old_site():
    using_old_site = request.host == settings.OLD_SITE_HOST
    return {'using_old_site': using_old_site}


@app.context_processor
def inject_show_dict():
    return {'series_map': app.config['SHOW_DICT_WITH_NAMES']}


@app.context_processor
def inject_show_list():
    return {'show_list': app.config['SHOW_DICT']}
