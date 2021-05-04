import json
import logging
import os
from typing import List

import mistune
import sys

from link_renderer import LinkRenderer, RendererOptions


class FinderOptions:
    def __init__(self, ignore_list=None, using_hugo=False) -> None:
        if ignore_list is None:
            ignore_list = []
        self.ignore_list = ignore_list
        self.using_hugo = using_hugo

def find_links(file_path: str, options: FinderOptions):
    # credit to https://github.com/MatMoore/markdown-external-link-finder
    render_opts = RendererOptions(ignore=options.ignore_list)
    link_renderer = LinkRenderer(render_opts)
    renderer = mistune.Markdown(renderer=link_renderer)

    with open(file_path) as f:
        text = f.read()
        renderer(text)

    return link_renderer.urls


def generate_output():
    try:
        files = os.environ["INPUT_FILES"].split(",")
    except KeyError:
        logging.error("The `files` input was not set. Expected a comma separated list of absolute file paths.")
        sys.exit(1)

    try:
        url_blacklist: List[str] = os.environ["INPUT_URL_BLACKLIST"].split(",")
    except KeyError:
        url_blacklist: List[str] = []

    opts = FinderOptions(ignore_list=url_blacklist, using_hugo=True)
    links: List[str] = []
    for file_path in files:
        if file_path.endswith(".md"):
            links = links + find_links(file_path, opts)

    return json.dumps({"links": links})


if __name__ == "__main__":
    logging.basicConfig(
        level=logging.INFO,
        format="%(asctime)s [%(levelname)s] %(message)s",
        handlers=[
            logging.FileHandler("debug.log"),
            logging.StreamHandler()
        ]
    )
    json_output = generate_output()
    logging.info(json_output)
    print(f"::set-output name=links::{json_output}")
