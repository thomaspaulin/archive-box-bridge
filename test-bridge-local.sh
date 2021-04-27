#!/bin/bash

curl --connect-timeout 20 \
	-H 'Content-Type: application/json' \
	-d '{"links": ["https://en.wikipedia.org", "https://www.example.com", "https://www.example.org"]}' \
	-iLv \
	"http://localhost:3344/"
