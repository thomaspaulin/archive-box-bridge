#!/bin/bash
export ARCHIVE_BOX_HOST="TODO"

curl --connect-timeout 20 \
	-H 'Content-Type: application/json' \
	-d '{"links": ["https://en.wikipedia.org", "https://www.example.com", "https://www.example.org"]}' \
	--post301 \
	-iLv \
	"$ARCHIVE_BOX_HOST/archive-links"
