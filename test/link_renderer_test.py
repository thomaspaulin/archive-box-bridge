import unittest

from link_renderer import LinkRenderer


class LinkRendererTests(unittest.TestCase):

    def test_link_ignored(self):
        ignore_list = ['en.wikipedia.org']
        sut = LinkRenderer(ignore_list)
        self.assertTrue(sut.is_ignored("https://en.wikipedia.org/wiki/Python_%28programming_language%29"))

    def test_invalid_url_throws_error(self):
        sut = LinkRenderer([])
        with self.assertRaises(ValueError):
            sut.link("example.com/example", "", "")

    def test_no_duplicates(self):
        sut = LinkRenderer([])
        sut.link("https://en.wikipedia.org/wiki/Python_%28programming_language%29", "", "")
        sut.link("https://en.wikipedia.org/wiki/Python_%28programming_language%29", "", "")
        self.assertEqual(1, len(sut.urls))

    def test_valid_links_added(self):
        sut = LinkRenderer([])
        sut.link("https://en.wikipedia.org/wiki/Python_%28programming_language%29", "", "")
        sut.link("https://thomaspaulin.me/", "", "")
        self.assertEqual(2, len(sut.urls))
        self.assertIn("https://en.wikipedia.org/wiki/Python_%28programming_language%29", sut.urls)
        self.assertIn("https://thomaspaulin.me/", sut.urls)


if __name__ == '__main__':
    unittest.main()
