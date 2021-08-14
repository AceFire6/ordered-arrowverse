from tortoise import fields
from tortoise.models import Model

from .settings import Shows


class Episode(Model):
    name = fields.CharEnumField(Shows)
    show = fields.ForeignKeyField('models.Show', related_name='episodes', on_delete=fields.RESTRICT)
    synopsis = fields.TextField()
    link = fields.CharField(max_length=255)
    air_date = fields.DatetimeField()

    def __str__(self):
        return self.name


class Show(Model):
    name = fields.TextField()
    html_id = fields.CharField(max_length=50)
    source_url = fields.CharField(max_length=255)

    def __str__(self):
        return self.name


class ShowEpisodeOrder(Model):
    show_order = fields.JSONField(on_delete=fields.RESTRICT)
    start = fields.DatetimeField(null=True)
    end = fields.DatetimeField(null=True)
