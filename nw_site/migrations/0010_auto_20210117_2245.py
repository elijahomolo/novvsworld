# Generated by Django 3.1.2 on 2021-01-17 22:45

from django.db import migrations, models


class Migration(migrations.Migration):

    dependencies = [
        ('nw_site', '0009_archives_template'),
    ]

    operations = [
        migrations.AlterField(
            model_name='archives',
            name='template',
            field=models.CharField(default='archives.html', max_length=50),
        ),
    ]