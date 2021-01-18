from django.shortcuts import get_object_or_404, render
from django.template import loader

from django.http import HttpResponse

from .models import Post
from .models import Archives

# Create your views here.

def index(request):
    latest_post_list = Post.objects.order_by('-pub_date')[:5]
    context = {'latest_post_list': latest_post_list}
    return render(request, 'index.html', context)

def music(request):
    return render(request, 'music.html')

def vault(request):
    return render(request, 'vault.html')

def archives(request):
    latest_archives_list = Archives.objects.order_by('-pub_date')
    context = {'latest_archives_list': latest_archives_list}
    return render(request, 'archives.html', context)

def contact(request):
    return render(request, 'contact.html')

def post(request, post_id):
    post = get_object_or_404(Post, id=post_id)
    _t = Post.objects.filter(id=post_id).order_by('template').values('template').first()
    template = _t['template']
    return render(request, template, {'post': post})

def archive_post(request, archives_id):
    archives_post = get_object_or_404(Archives, id=archives_id)
    _t = Archives.objects.filter(id=archives_id).order_by('template').values('template').first()
    template = _t['template']
    return render(request, template, {'archives_post': archives_post})



def releases(request):
    #template = loader.get_template('ex_mwingine.html')
    return render(request, 'release.html')
