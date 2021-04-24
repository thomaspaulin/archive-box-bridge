from typing import List
from urllib.parse import urlsplit

from mistune import Renderer


def url_is_absolute(url):
    return bool(urlsplit(url).netloc)


class LinkRenderer(Renderer):
    def __init__(self, ignore: List[str]):
        self.urls = []
        self.ignore = ignore
        super().__init__()

    def is_ignored(self, url) -> bool:
        return urlsplit(url).netloc in self.ignore

    def is_duplicate(self, url) -> bool:
        return url in self.urls

    def link(self, link, title, text):
        if url_is_absolute(link):
            if self.is_duplicate(link) is False and self.is_ignored(link) is False:
                self.urls.append(link)
        else:
            raise ValueError("Invalid URL: " + link)
