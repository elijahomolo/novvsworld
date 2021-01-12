from django.db import models

# Create your models here.
class Post(models.Model):
    title = models.CharField(max_length=50)
    body = models.TextField()
    image_upload = models.ImageField(upload_to='uploads/')
    pub_date = models.DateTimeField('date published')
    def __str__(self):
        return self.body
