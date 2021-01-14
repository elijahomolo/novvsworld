from django.db import models

# Create your models here.
class Post(models.Model):
    title = models.CharField(max_length=50)
    body = models.TextField()
    image_upload = models.ImageField()
    pub_date = models.DateTimeField('date published')
    link_to = models.CharField(max_length=100, default='/')
    def __str__(self):
        return self.body

class Archives(models.Model):
    title = models.CharField(max_length=50)
    body = models.TextField()
    image_upload = models.ImageField()
    pub_date = models.DateTimeField('date published')
    link_to = models.CharField(max_length=100, default='/archives')
    def __str__(self):
        return self.body
