from django.urls import path

from . import views

urlpatterns = [
    path('', views.index, name='index'),
    path('releases',views.releases, name='releases' ),
    path('<int:post_id>/', views.post, name='post'),
]
