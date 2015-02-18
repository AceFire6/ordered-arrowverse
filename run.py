#!virtualenv/bin/python

from ordering import app

app.run(host=app.config['HOST'], port=app.config['PORT'], debug=True)
