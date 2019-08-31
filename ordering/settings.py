import os
from datetime import timedelta

USERNAME = 'user'
PASSWORD = 'pass'
HOST = '0.0.0.0'
PORT = 5000

DEBUG = os.getenv('FLASK_DEBUG', False)

CACHE_TYPE = 'redis' if not DEBUG else 'null'
CACHE_REDIS_URL = os.getenv('REDIS_URL', 'redis://localhost:6379')

SEND_FILE_MAX_AGE_DEFAULT = timedelta(weeks=1)

ARROW_URL = 'List_of_Arrow_episodes'
CONSTANTINE_URL = 'List_of_Constantine_episodes'
FLASH_URL = 'List_of_The_Flash_(The_CW)_episodes'
FREEDOM_FIGHTERS_URL = 'List_of_Freedom_Fighters:_The_Ray_episodes'
LEGENDS_URL = "List_of_DC's_Legends_of_Tomorrow_episodes"
SUPERGIRL_URL = 'List_of_Supergirl_episodes'
VIXEN_URL = 'List_of_Vixen_episodes'
BLACK_LIGHTNING_URL = 'Black_Lightning_(TV_series)'
BATWOMAN_URL = 'List_of_Batwoman_episodes'

FANDOM_ROOT = 'http://arrow.fandom.com/wiki/'
WIKIPEDIA_ROOT = 'https://en.wikipedia.org/wiki/'

SHOWS = (
    {
        'id': 'arrow',
        'name': 'Arrow',
        'url': ARROW_URL,
        'root': FANDOM_ROOT
    },
    {
        'id': 'constantine',
        'name': 'Constantine',
        'url': CONSTANTINE_URL,
        'root': FANDOM_ROOT
    },
    {
        'id': 'flash',
        'name': 'The Flash',
        'url': FLASH_URL,
        'root': FANDOM_ROOT
    },
    {
        'id': 'freedom-fighters',
        'name': 'Freedom Fighters: The Ray',
        'url': FREEDOM_FIGHTERS_URL,
        'root': FANDOM_ROOT
    },
    {
        'id': 'legends',
        'name': "DC's Legends of Tomorrow",
        'url': LEGENDS_URL,
        'root': FANDOM_ROOT
    },
    {
        'id': 'supergirl',
        'name': 'Supergirl',
        'url': SUPERGIRL_URL,
        'root': FANDOM_ROOT
    },
    {
        'id': 'vixen',
        'name': 'Vixen',
        'url': VIXEN_URL,
        'root': FANDOM_ROOT
    },
    {
        'id': 'black-lightning',
        'name': 'Black Lightning',
        'url': BLACK_LIGHTNING_URL,
        'root': WIKIPEDIA_ROOT
    },
    {
        'id': 'batwoman',
        'name': 'Batwoman',
        'url': BATWOMAN_URL,
        'root': FANDOM_ROOT
    },
)

SHOW_DICT = {SHOWS[i]['id']: SHOWS[i] for i in range(len(SHOWS))}

SHOW_DICT_WITH_NAMES = {SHOWS[i]['name']: SHOWS[i] for i in range(len(SHOWS))}
SHOW_DICT_WITH_NAMES.update(SHOW_DICT)
