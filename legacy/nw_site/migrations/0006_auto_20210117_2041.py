# Generated by Django 3.1.2 on 2021-01-17 20:41

from django.db import migrations, models


class Migration(migrations.Migration):

    dependencies = [
        ('nw_site', '0005_auto_20210117_2033'),
    ]

    operations = [
        migrations.AlterField(
            model_name='archives',
            name='image_upload',
            field=models.ImageField(height_field=200, upload_to='', width_field=200),
        ),
        migrations.AlterField(
            model_name='post',
            name='image_upload',
            field=models.ImageField(height_field=200, upload_to='', width_field=200),
        ),
    ]