from django.shortcuts import render
from django.template import loader

from django.http import HttpResponse

from .models import Post

# Create your views here.

def index(request):
    latest_post_list = Post.objects.order_by('-pub_date')[:5]
    context = {'latest_post_list': latest_post_list}
    return render(request, 'index.html', context)

def music(request):
    return render(request, 'music.html')

def vault(request):
    return render(request, 'vault.html')

def post(request, post_id):
    return HttpResponse("You're looking at post %s." % post_id)

def releases(request):
    #template = loader.get_template('ex_mwingine.html')
    return render(request, 'release.html')
