import json
import logging
import os
from typing import List

import mistune
import sys

from link_renderer import LinkRenderer


def find_links(file_path: str, ignore_list: List[str]):
    # credit to https://github.com/MatMoore/markdown-external-link-finder
    link_renderer = LinkRenderer(ignore_list)
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

    links: List[str] = []
    for file_path in files:
        if file_path.endswith(".md"):
            links = links + find_links(file_path, url_blacklist)

    return json.dumps({"links": links})


if __name__ == "__main__":
    logging.basicConfig(filename='archive.log', level=logging.INFO)
    json_output = generate_output()
    logging.info(json_output)
    print(f"::set-output name=links::{json_output}")
