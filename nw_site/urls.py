from django.urls import path

from django.conf import settings
from django.conf.urls.static import static

from . import views

urlpatterns = [
    path('', views.index, name='index'),
    path('music', views.music, name='music'),
    path('vault', views.vault, name='vault'),
    path('archives', views.archives, name='archives'),
    path('releases',views.releases, name='releases' ),
    path('<int:post_id>/', views.post, name='post'),
    path('contact', views.contact, name='contact'),
    path('archives/<int:archives_id>/', views.archive_post, name='archives_post'),
] + static(settings.MEDIA_URL, document_root=settings.MEDIA_ROOT)
