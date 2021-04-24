# Markdown Link Finder

This action prints "Hello World" or "Hello" + the name of a person to greet to the log.

## Inputs

### `files`
**Required** A comma separated list of files to search within. These should be the absolute file paths.  Any files which do not end in *.md will be ignored'

### `url_blacklist`
A comma separated list of hosts to which should not be archived. This represents URLs which are considered reliable and are unlikely rot e.g., en.wikipedia.org'. The default setting permits every URL.

## Outputs

### `links`
The links found within the provided markdown files.

## Example Usage
```yaml
uses: thomaspaulin/markdown-link-finder@v1
with:
  files: 'README.md,example.md'
  url_blacklist: 'en.wikipedia.org,news.ycombinator.com'
```
