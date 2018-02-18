USERNAME = 'user'
PASSWORD = 'pass'
HOST = '0.0.0.0'
PORT = 5000

CACHE_TYPE = 'simple'

ARROW_URL = 'List_of_Arrow_episodes'
CONSTANTINE_URL = 'List_of_Constantine_episodes'
FLASH_URL = 'List_of_The_Flash_episodes'
FREEDOM_FIGHTERS_URL = 'List_of_Freedom_Fighters:_The_Ray_episodes'
LEGENDS_URL = "List_of_DC's_Legends_of_Tomorrow_episodes"
SUPERGIRL_URL = 'List_of_Supergirl_episodes'
VIXEN_URL = 'List_of_Vixen_episodes'

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
        'id': 'constantine',
        'name': 'Constantine',
        'url': CONSTANTINE_URL,
        'root': WIKIA_ROOT
    },
    {
        'id': 'flash',
        'name': 'The Flash',
        'url': FLASH_URL,
        'root': WIKIA_ROOT
    },
    {
        'id': 'freedom-fighters',
        'name': 'Freedom Fighters: The Ray',
        'url': FREEDOM_FIGHTERS_URL,
        'root': WIKIA_ROOT
    },
    {
        'id': 'legends',
        'name': "DC's Legends of Tomorrow",
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
        'root': WIKIA_ROOT
    },
)

SHOW_DICT = {SHOWS[i]['id']: SHOWS[i] for i in range(len(SHOWS))}

SHOW_DICT_WITH_NAMES = {SHOWS[i]['name']: SHOWS[i] for i in range(len(SHOWS))}
SHOW_DICT_WITH_NAMES.update(SHOW_DICT)
