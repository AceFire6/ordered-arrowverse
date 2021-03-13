from datetime import timedelta
from enum import Enum

from environs import Env

env = Env()

DEBUG = env.bool('QUART_DEBUG', default=False)

REDIS_URL = env('REDIS_URL', 'redis://localhost:6379')

SEND_FILE_MAX_AGE_DEFAULT = timedelta(weeks=1)


class Shows(str, Enum):
    ARROW = 'Arrow'
    BATWOMAN = 'Batwoman'
    BLACK_LIGHTNING = 'Black Lightning'
    CONSTANTINE = 'Constantine'
    FLASH = 'The Flash'
    FREEDOM_FIGHTERS = 'Freedom Fighters: The Ray'
    LEGENDS = "DC's Legends of Tomorrow"
    SUPERGIRL = 'Supergirl'
    SUPERMAN_AND_LOIS = 'Superman & Lois'
    VIXEN = 'Vixen'

    def __str__(self):
        return self.value


ARROW_URL = 'List_of_Arrow_episodes'
BATWOMAN_URL = 'List_of_Batwoman_episodes'
BLACK_LIGHTNING_URL = 'List_of_Black_Lightning_episodes'
CONSTANTINE_URL = 'List_of_Constantine_episodes'
FLASH_URL = 'List_of_The_Flash_(The_CW)_episodes'
FREEDOM_FIGHTERS_URL = 'List_of_Freedom_Fighters:_The_Ray_episodes'
LEGENDS_URL = "List_of_DC's_Legends_of_Tomorrow_episodes"
SUPERGIRL_URL = 'List_of_Supergirl_episodes'
SUPERMAN_AND_LOIS_URL = 'List_of_Superman_%26_Lois_episodes'
VIXEN_URL = 'List_of_Vixen_episodes'

FANDOM_ROOT = 'http://arrow.fandom.com/wiki/'
WIKIPEDIA_ROOT = 'https://en.wikipedia.org/wiki/'

DATABASE_URL = env('DATABASE_URL', default='postgres://localhost:5432/arrowverse_db')
if not DEBUG:
    DATABASE_URL = f'{DATABASE_URL}?ssl=require'

SHOWS = (
    {
        'id': 'arrow',
        'name': Shows.ARROW,
        'url': ARROW_URL,
        'root': FANDOM_ROOT
    },
    {
        'id': 'batwoman',
        'name': Shows.BATWOMAN,
        'url': BATWOMAN_URL,
        'root': FANDOM_ROOT
    },
    {
        'id': 'black-lightning',
        'name': Shows.BLACK_LIGHTNING,
        'url': BLACK_LIGHTNING_URL,
        'root': WIKIPEDIA_ROOT
    },
    {
        'id': 'constantine',
        'name': Shows.CONSTANTINE,
        'url': CONSTANTINE_URL,
        'root': FANDOM_ROOT
    },
    {
        'id': 'flash',
        'name': Shows.FLASH,
        'url': FLASH_URL,
        'root': FANDOM_ROOT
    },
    {
        'id': 'freedom-fighters',
        'name': Shows.FREEDOM_FIGHTERS,
        'url': FREEDOM_FIGHTERS_URL,
        'root': FANDOM_ROOT
    },
    {
        'id': 'legends',
        'name': Shows.LEGENDS,
        'url': LEGENDS_URL,
        'root': FANDOM_ROOT
    },
    {
        'id': 'supergirl',
        'name': Shows.SUPERGIRL,
        'url': SUPERGIRL_URL,
        'root': FANDOM_ROOT
    },
    {
        'id': 'superman-and-lois',
        'name': Shows.SUPERMAN_AND_LOIS,
        'url': SUPERMAN_AND_LOIS_URL,
        'root': FANDOM_ROOT
    },
    {
        'id': 'vixen',
        'name': Shows.VIXEN,
        'url': VIXEN_URL,
        'root': FANDOM_ROOT
    },
)

SHOW_DICT = {SHOWS[i]['id']: SHOWS[i] for i in range(len(SHOWS))}

SHOW_DICT_WITH_NAMES = {SHOWS[i]['name']: SHOWS[i] for i in range(len(SHOWS))}
SHOW_DICT_WITH_NAMES.update(SHOW_DICT)
