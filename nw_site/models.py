from django.db import models

# Create your models here.
class Post(models.Model):
    title = models.CharField(max_length=50)
    description = models.CharField(max_length=100, default=title)
    body = models.TextField()
    image_upload = models.ImageField()
    pub_date = models.DateTimeField('date published')
    link_to = models.CharField(max_length=100, default='/')
    template = models.CharField(max_length=50, default='post.html')
    def __str__(self):
        return self.body

class Archives(models.Model):
    title = models.CharField(max_length=50)
    description = models.CharField(max_length=100, default=title)
    body = models.TextField()
    image_upload = models.ImageField()
    pub_date = models.DateTimeField('date published')
    link_to = models.CharField(max_length=100, default='/archives')
    template = models.CharField(max_length=50, default='archive_post.html')
    def __str__(self):
        return self.body
