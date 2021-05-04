import os
import unittest

from link_finder import find_links, FinderOptions


class LinkFinderTests(unittest.TestCase):

    def test_find_links_on_valid_markdown_file_returns_links(self):
        fp = os.path.join(os.curdir, "resources", "file.md")
        opts = FinderOptions([])
        links = find_links(fp, opts)
        self.assertIn("https://en.wikipedia.org", links)
        self.assertIn("https://www.example.com", links)
        self.assertIn("https://www.example.org", links)
        self.assertEqual(3, len(links))

    def test_find_links_in_hugo_markdown_file_with_shortcodes(self):
        fp = os.path.join(os.curdir, "resources", "hugo.md")
        opts = FinderOptions([])
        links = find_links(fp, opts)
        self.assertIn("https://en.wikipedia.org", links)
        self.assertIn("https://www.example.com", links)
        self.assertIn("https://www.example.org", links)
        self.assertEqual(3, len(links))
    # todo test cases when env vars aren't set, and non-markdown files


if __name__ == '__main__':
    unittest.main()
