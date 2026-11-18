# randname
[![Build Status](https://github.com/marstr/randname/actions/workflows/ci.yml/badge.svg)](https://github.com/marstr/randname/actions/workflows/ci.yml)

A Go package which creates random names with an emphasis on flexibility and ease of use.

The default dictionaries are embedded at build time and do not need to be distributed separately. Go 1.16 or later is required.

# Contributing

Feel free to send PRs my way! Just make sure that your code passes `go test`. If you add any publicly exported methods, please document them.

To contribute words, edit `adjectives.txt` or `nouns.txt`, keeping one word per line in alphabetical order. The files are embedded automatically on the next build; no Go code generation is needed.
