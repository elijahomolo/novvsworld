from django.urls import path

from django.conf import settings
from django.conf.urls.static import static

from . import views

urlpatterns = [
    path('', views.index, name='index'),
    path('music', views.music, name='music'),
    path('vault', views.vault, name='vault'),
    path('releases',views.releases, name='releases' ),
    path('<int:post_id>/', views.post, name='post'),
] + static(settings.MEDIA_URL, document_root=settings.MEDIA_ROOT)
