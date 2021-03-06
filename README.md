# Search the web from the command line using a Google Custom Search engine

Note: Should work anywhere you have access to a Factotum (Plan9, 9front, plan9port)

## Installation
```
go get github.com/halfwit/google
go install github.com/halfwit/google
```

## Usage 

google [[-i] [-it type] [-is size] [-ic color] [-id scale]] [-m results] [-r url] [-a key] [-f type] [-e match] [-x exclude] [ -s safe] query
 - -m Number of results per query
 - -r Search for sites related to [url]
 - -a Use API key instead of factotum
 - -f File type [bmp|gif|png|jpg|svg|pdf]
 - -e Match string exactly
 - -x Phrase to exclude
 - -s Safe search [active|high|medium|off]
 - -i Image search
### Requires -i
 - -it Image type [clipart|face|lineart|news|photo]
 - -is Image size [huge|icon|large|medium|small|xlarge|xxlarge]
 - -ic Image color [black|blue|brown|gray|green|orange|ping|purple|red|teal|white|yellow]
 - -id Image scale [color|gray|mono]

## Authentication

To use this, you need an API key for Google Custom Search engine.
See https://developers.google.com/custom-search/v1/introduction, select "Get A Key"

Additionally, you also have to create a Custom Search Engine itself.
When you create a Custom Search Engine, you'll also get a `cx` key.

Store these keys in your Factotum
 - `gcli` - your API key (AIza...)
 - `gcse` - your CSE key

## Performance (Subjective)

Using this, I was able to complete a search with 50 results, in 0.03 seconds.
A similar search on the popular Python program `googler` was roughly 0.19 seconds.
