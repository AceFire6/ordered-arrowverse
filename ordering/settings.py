USERNAME = 'user'
PASSWORD = 'pass'
HOST = '0.0.0.0'
PORT = 5000

CACHE_TYPE = 'simple'

ARROW_URL = 'http://arrow.wikia.com/wiki/List_of_Arrow_episodes'
FLASH_URL = 'http://arrow.wikia.com/wiki/List_of_The_Flash_episodes'
LEGENDS_URL = ('http://arrow.wikia.com/wiki/'
               'List_of_DC%27s_Legends_of_Tomorrow_episodes')

SHOWS = (
    {
        'name': 'Arrow',
        'url': ARROW_URL
    },
    {
        'name': 'The Flash',
        'url': FLASH_URL
    },
    {
        'name': 'DC\'s Legends of Tomorrow',
        'url': LEGENDS_URL
    },
)
