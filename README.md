# randname
[![Go Reference](https://pkg.go.dev/badge/github.com/marstr/randname.svg)](https://pkg.go.dev/github.com/marstr/randname) [![Build Status](https://github.com/marstr/randname/actions/workflows/ci.yml/badge.svg)](https://github.com/marstr/randname/actions/workflows/ci.yml)

A Go package which creates random names with an emphasis on flexibility and ease of use.

The default dictionaries are embedded at build time and do not need to be distributed separately. Go 1.16 or later is required.

## Installation

From your Go project's directory, add the dependency:

```sh
go get github.com/marstr/randname@v0.2
```

If you're starting a new project, run `go mod init example.com/myproject` first, replacing the module path with your own.

## Quick start

Save this as `main.go` and run it with `go run .`:

```go
package main

import (
	"fmt"

	"github.com/marstr/randname"
)

func main() {
	fmt.Println(randname.Generate())
}
```

`Generate()` combines a random adjective, noun, and number from 00 to 99 in PascalCase, such as `BigCloud09`. No configuration is needed.

### Choose a format

Use these helpers in place of `Generate()` to choose how the name is formatted:

```go
fmt.Println(randname.GeneratePascalCase()) // e.g. BigCloud09
fmt.Println(randname.GenerateCamelCase())  // e.g. bigCloud09
fmt.Println(randname.GenerateHyphened())   // e.g. big-cloud-09
```

Each call generates a new random name; the comments illustrate the formats, not fixed output.

### Configure a generator

The zero value of `AdjNoun` is ready to use. Nil `Adjectives` and `Nouns` use the embedded dictionaries, nil `RandGenerator` uses `crypto/rand.Reader`, and nil `Format` uses PascalCase. A non-nil empty dictionary omits that word.

```go
generator := randname.AdjNoun{Format: randname.FormatHyphenedAdjNoun}
fmt.Println(generator.Generate())          // e.g. big-cloud-09
fmt.Println(generator.GenerateCamelCase()) // e.g. bigCloud09
fmt.Println(generator.Generate())          // still hyphened, e.g. small-river-42
```

The named formatting methods ignore `Format` for that call without changing it. A custom `AdjNounFormat` receives an adjective, noun, and integer from 0 to 99; it controls the output, including number padding.

### Use a prefix

For a name with your own prefix and a random alphanumeric suffix:

```go
fmt.Println(randname.GenerateWithPrefix("worker-", 6)) // e.g. worker-aB3xY7
```

The second argument is the suffix length; passing zero uses the default length of six.

See the [Go reference](https://pkg.go.dev/github.com/marstr/randname) for custom dictionaries, formatting functions, and other options.

# Contributing

Feel free to send PRs my way! Just make sure that your code passes `go test`. If you add any publicly exported methods, please document them.

To contribute words, edit `adjectives.txt` or `nouns.txt`, keeping one word per line in alphabetical order. The files are embedded automatically on the next build; no Go code generation is needed.
