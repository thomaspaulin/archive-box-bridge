import unittest

from link_renderer import LinkRenderer, RendererOptions


class LinkRendererTests(unittest.TestCase):

    def test_link_ignored(self):
        ignore_list = ['en.wikipedia.org']
        opts = RendererOptions(ignore_list)
        sut = LinkRenderer(opts)
        self.assertTrue(sut.is_ignored("https://en.wikipedia.org/wiki/Python_%28programming_language%29"))

    def test_invalid_url_throws_error(self):
        opts = RendererOptions([])
        sut = LinkRenderer(opts)
        with self.assertRaises(ValueError):
            sut.link("example.com/example", "", "")

    def test_no_duplicates(self):
        opts = RendererOptions([])
        sut = LinkRenderer(opts)
        sut.link("https://en.wikipedia.org/wiki/Python_%28programming_language%29", "", "")
        sut.link("https://en.wikipedia.org/wiki/Python_%28programming_language%29", "", "")
        self.assertEqual(1, len(sut.urls))

    def test_valid_links_added(self):
        opts = RendererOptions([])
        sut = LinkRenderer(opts)
        sut.link("https://en.wikipedia.org/wiki/Python_%28programming_language%29", "", "")
        sut.link("https://thomaspaulin.me/", "", "")
        self.assertEqual(2, len(sut.urls))
        self.assertIn("https://en.wikipedia.org/wiki/Python_%28programming_language%29", sut.urls)
        self.assertIn("https://thomaspaulin.me/", sut.urls)

    def test_ignores_hugo_shortcodes(self):
        opts = RendererOptions([])
        sut = LinkRenderer(opts)
        sut.link('{{< relref "file" >}}', "", "")
        self.assertEqual(0, len(sut.urls))

if __name__ == '__main__':
    unittest.main()
