#!/usr/bin/env fish

set urls (cat urls.txt)

for url in $urls
    go run main.go ankify --type=url --cards=10 --tag=learning $url
end
