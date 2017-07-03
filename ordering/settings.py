USERNAME = 'user'
PASSWORD = 'pass'
HOST = '0.0.0.0'
PORT = 5000

CACHE_TYPE = 'simple'

ARROW_URL = 'List_of_Arrow_episodes'
FLASH_URL = 'List_of_The_Flash_episodes'
LEGENDS_URL = 'List_of_DC%27s_Legends_of_Tomorrow_episodes'
SUPERGIRL_URL = 'List_of_Supergirl_episodes'
VIXEN_URL = 'List_of_Vixen_episodes'
CONSTANTINE_URL = 'List_of_Constantine_episodes'

WIKIA_ROOT = 'http://arrow.wikia.com/wiki/'
WIKIPEDIA_ROOT = 'https://en.wikipedia.org/wiki/'

SHOWS = (
    {
        'id': 'arrow',
        'name': 'Arrow',
        'url': ARROW_URL,
        'root': WIKIA_ROOT
    },
    {
        'id': 'flash',
        'name': 'The Flash',
        'url': FLASH_URL,
        'root': WIKIA_ROOT
    },
    {
        'id': 'legends',
        'name': 'DC\'s Legends of Tomorrow',
        'url': LEGENDS_URL,
        'root': WIKIA_ROOT
    },
    {
        'id': 'supergirl',
        'name': 'Supergirl',
        'url': SUPERGIRL_URL,
        'root': WIKIPEDIA_ROOT
    },
    {
        'id': 'vixen',
        'name': 'Vixen',
        'url': VIXEN_URL,
        'root': ARROW_ROOT
    },
    {
        'id': 'constantine',
        'name': 'Constantine',
        'url': CONSTANTINE_URL,
        'root': ARROW_ROOT
    }
)

SHOW_DICT = {SHOWS[i]['id']: SHOWS[i] for i in xrange(len(SHOWS))}
SHOW_DICT.update({SHOWS[i]['name']: SHOWS[i] for i in xrange(len(SHOWS))})
