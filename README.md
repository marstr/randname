# randname
[![Build Status](https://mstrobel.visualstudio.com/randname/_apis/build/status/randname-CI)](https://mstrobel.visualstudio.com/randname/_build/latest?definitionId=6)

A Go package which creates random names with an emphasis on flexibility and ease of use.

The default dictionaries are embedded at build time and do not need to be distributed separately. Go 1.16 or later is required.

# Contributing

Feel free to send PRs my way! Just make sure that your code passes `golint` and `go test`. If you add any publicly exported methods, please document them.

To contribute words, edit `adjectives.txt` or `nouns.txt`, keeping one word per line in alphabetical order. The files are embedded automatically on the next build; no Go code generation is needed.
