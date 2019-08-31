from quart.routing import BaseConverter


class ListConverter(BaseConverter):
    def to_python(self, value):
        return value.split('+')

    def to_url(self, values):
        return '+'.join(super().to_url(value) for value in values)
