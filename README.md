# Arrowverse Series Ordering
[![Build Status](https://travis-ci.org/AceFire6/ordered-arrowverse.svg?branch=main)](https://travis-ci.org/AceFire6/ordered-arrowverse)

This is a project that aims to centralize the series in the Arrowverse that
have crossovers to make it easy to watch the episodes in the correct
order.

### How It Works:

This project uses data from the episode lists on the [Arrow Wiki](http://arrow.fandom.com)
and the episodes are ordered by air date.

Data is also drawn from Wikipedia when necessary.

### Development

The project has been tested in Python 3.9 and will (currently) not run on Python 3.10+

During development, the environment variable `REDIS_URL` can be set to "None"
removing the requirement of having a Redis instance up and running. This
effectively disables caching.

Setting up the environment:

```
poetry install
```

Running a development instance:

There are two options:

```
QUART_ENV=development QUART_DEBUG=true QUART_APP=./ordering/ poetry run quart run
```

and 

```
poetry shell
QUART_ENV=development QUART_DEBUG=true QUART_APP=./ordering/ quart run
```


### Currently Supported Series:

* Arrow
* Batwoman
* Black Lightning
* Constantine
* DC's Legends of Tomorrow
* The Flash
* Freedom Fighters: The Ray
* Stargirl
* Supergirl
* Superman & Lois
* Vixen
