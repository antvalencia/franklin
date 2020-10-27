# franklin
A little API for news feed entries

---

## Example feed list
See feed initialization example file **./models/data/feed.csv**

CSV format is:
> topic,source,title,URL

headers **topic**, **source** & **title** are optional categories for convenience.  **URL** must be a valid RSS or Atom feed URL

---

## Setup database

> ./scripts/db_build.sh

But init env vars found therein.

---

## Run a test to add some news items
From test:
> go run add_news_items.go
