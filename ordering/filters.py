from ordering import app


def url_form(episode_name):
    return episode_name.replace(' ', '_')


@app.template_filter('episode_url')
def episode_url_filter(episode_name, series):
    root_url = app.config['ROOT_URL']
    if episode_name == 'Pilot':
        return root_url + url_form(episode_name + ' (%s)' % series)
    else:
        return root_url + url_form(episode_name)
