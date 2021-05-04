from typing import List
from urllib.parse import urlsplit

from mistune import Renderer


def url_is_absolute(url):
    return bool(urlsplit(url).netloc)


class RendererOptions:
    def __init__(self, ignore: List[str]) -> None:
        if ignore is None:
            self.ignore = []
        else:
            self.ignore = ignore
        super().__init__()


class LinkRenderer(Renderer):
    def __init__(self, opts: RendererOptions):
        self.urls = []
        self.opts = opts
        super().__init__()

    def is_ignored(self, url) -> bool:
        return urlsplit(url).netloc in self.opts.ignore

    def is_duplicate(self, url) -> bool:
        return url in self.urls

    def is_hugo_shortcode(self, url: str) -> bool:
        # todo this is implementation specific. Refactor it away should another implementation ever be involved
        return url.startswith("{{<") and url.endswith(">}}")

    def link(self, link, title, text):
        if self.is_hugo_shortcode(link):
            pass
        elif url_is_absolute(link):
            if self.is_duplicate(link) is False and self.is_ignored(link) is False:
                self.urls.append(link)
        else:
            raise ValueError("Invalid URL: " + link)
