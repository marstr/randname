package randname

import (
	"bytes"
	"crypto/rand"
	_ "embed"
	"fmt"
	"io"
	"math/big"
	"strings"
	"sync"
	"unicode"

	"github.com/marstr/collection"
)

// AdjNounFormat combines an adjective, noun, and number into a formatted name.
// AdjNoun passes a number from 0 to 99, inclusive.
type AdjNounFormat func(string, string, int) string

// AdjNoun generates a random adjective, noun, and number from 0 to 99, inclusive.
// Its zero value is ready to use and generates PascalCase names like BigCloud09.
type AdjNoun struct {
	// Adjectives supplies the adjective dictionary. Nil uses the embedded defaults;
	// a non-nil empty dictionary omits the adjective.
	Adjectives *collection.Dictionary
	// Format controls Generate's output. Nil uses FormatPascalCaseAdjNoun.
	Format AdjNounFormat
	// Nouns supplies the noun dictionary. Nil uses the embedded defaults;
	// a non-nil empty dictionary omits the noun.
	Nouns *collection.Dictionary
	// RandGenerator supplies randomness. Nil uses crypto/rand.Reader.
	RandGenerator io.Reader
}

var defaultFormat AdjNounFormat = FormatPascalCaseAdjNoun
var defaultRandGenerator io.Reader = rand.Reader

// NewAdjNoun returns a generator using PascalCase formatting and crypto/rand.Reader.
// Its dictionary fields are nil, so generation uses the embedded defaults.
func NewAdjNoun() *AdjNoun {
	return &AdjNoun{
		RandGenerator: defaultRandGenerator,
		Format:        defaultFormat,
	}
}

func defaultBuilder(words *string) func() *collection.Dictionary {
	var loader sync.Once
	defaultDictionary := &collection.Dictionary{}
	return func() *collection.Dictionary {
		loader.Do(func() {
			reader := strings.NewReader(*words)
			var word string
			for {
				if _, err := fmt.Fscanln(reader, &word); err != nil {
					break
				}
				defaultDictionary.Add(word)
			}
			*words = ""
		})
		return defaultDictionary
	}
}

//go:embed adjectives.txt
var defaultAdjectives string

//go:embed nouns.txt
var defaultNouns string

var getDefaultAdjectives = defaultBuilder(&defaultAdjectives)
var getDefaultNouns = defaultBuilder(&defaultNouns)

// Generate returns a random name using the embedded dictionaries and PascalCase
// formatting, like BigCloud09. The number ranges from 00 to 99, inclusive.
func Generate() string {
	return AdjNoun{}.Generate()
}

// GeneratePascalCase returns a random name using the embedded dictionaries,
// formatted like BigCloud09, with a number from 00 to 99, inclusive.
func GeneratePascalCase() string {
	return AdjNoun{}.GeneratePascalCase()
}

// GenerateCamelCase returns a random name using the embedded dictionaries,
// formatted like bigCloud09, with a number from 00 to 99, inclusive.
func GenerateCamelCase() string {
	return AdjNoun{}.GenerateCamelCase()
}

// GenerateHyphened returns a random name using the embedded dictionaries,
// formatted like big-cloud-09, with a number from 00 to 99, inclusive.
func GenerateHyphened() string {
	return AdjNoun{}.GenerateHyphened()
}

// Generate returns a random name using the generator's dictionaries, randomness,
// and Format. A nil Format uses PascalCase, like BigCloud09.
// The number passed to Format ranges from 0 to 99, inclusive.
func (adNoun AdjNoun) Generate() string {
	if adNoun.Format == nil {
		adNoun.Format = defaultFormat
	}
	return adNoun.Format(adNoun.getAdjective(), adNoun.getNoun(), adNoun.getDigit())
}

// GeneratePascalCase returns a random name using the generator's dictionaries
// and randomness, formatted like BigCloud09, with a number from 00 to 99, inclusive.
// It ignores Format without changing it.
func (adNoun AdjNoun) GeneratePascalCase() string {
	return FormatPascalCaseAdjNoun(adNoun.getAdjective(), adNoun.getNoun(), adNoun.getDigit())
}

// GenerateCamelCase returns a random name using the generator's dictionaries
// and randomness, formatted like bigCloud09, with a number from 00 to 99, inclusive.
// It ignores Format without changing it.
func (adNoun AdjNoun) GenerateCamelCase() string {
	return FormatCamelCaseAdjNoun(adNoun.getAdjective(), adNoun.getNoun(), adNoun.getDigit())
}

// GenerateHyphened returns a random name using the generator's dictionaries
// and randomness, formatted like big-cloud-09, with a number from 00 to 99, inclusive.
// It ignores Format without changing it.
func (adNoun AdjNoun) GenerateHyphened() string {
	return FormatHyphenedAdjNoun(adNoun.getAdjective(), adNoun.getNoun(), adNoun.getDigit())
}

func (adNoun AdjNoun) getAdjective() string {
	if adNoun.Adjectives == nil {
		adNoun.Adjectives = getDefaultAdjectives()
	}
	if collection.Any(adNoun.Adjectives) {
		if adNoun.RandGenerator == nil {
			adNoun.RandGenerator = defaultRandGenerator
		}
		randomLocation, _ := rand.Int(adNoun.RandGenerator, big.NewInt(adNoun.Adjectives.Size()))
		return collection.ElementAt(adNoun.Adjectives, uint(randomLocation.Uint64())).(string)
	}
	return ""
}

func (adNoun AdjNoun) getNoun() string {
	if adNoun.Nouns == nil {
		adNoun.Nouns = getDefaultNouns()
	}
	if collection.Any(adNoun.Nouns) {
		if adNoun.RandGenerator == nil {
			adNoun.RandGenerator = defaultRandGenerator
		}
		position, _ := rand.Int(adNoun.RandGenerator, big.NewInt(adNoun.Nouns.Size()))
		return collection.ElementAt(adNoun.Nouns, uint(position.Uint64())).(string)
	}
	return ""
}

func (adNoun AdjNoun) getDigit() int {
	if adNoun.RandGenerator == nil {
		adNoun.RandGenerator = defaultRandGenerator
	}
	result, _ := rand.Int(adNoun.RandGenerator, big.NewInt(100))
	return int(result.Int64())
}

// FormatCamelCaseAdjNoun formats an adjective, noun, and number like bigCloud09.
// Nonnegative numbers are zero-padded to at least two decimal digits.
func FormatCamelCaseAdjNoun(adjective, noun string, digit int) string {
	pascal := FormatPascalCaseAdjNoun(adjective, noun, digit)
	return strings.ToLower(pascal[:1]) + pascal[1:]
}

// FormatPascalCaseAdjNoun formats an adjective, noun, and number like BigCloud09.
// Nonnegative numbers are zero-padded to at least two decimal digits.
func FormatPascalCaseAdjNoun(adjective, noun string, digit int) string {
	builder := bytes.Buffer{}
	if len(adjective) > 0 {
		builder.WriteRune(unicode.ToUpper(rune(adjective[0])))
		builder.WriteString(strings.ToLower(adjective[1:]))
	}

	if len(noun) > 0 {
		builder.WriteRune(unicode.ToUpper(rune(noun[0])))
		builder.WriteString(strings.ToLower(noun[1:]))
	}

	fmt.Fprintf(&builder, "%02d", digit)

	return builder.String()
}

// FormatHyphenedAdjNoun formats an adjective, noun, and number like big-cloud-09.
// Nonnegative numbers are zero-padded to at least two decimal digits.
func FormatHyphenedAdjNoun(adjective, noun string, digit int) string {
	builder := bytes.Buffer{}
	if len(adjective) > 0 {
		builder.WriteString(strings.ToLower(adjective))
		builder.WriteRune('-')
	}
	if len(noun) > 0 {
		builder.WriteString(strings.ToLower(noun))
		builder.WriteRune('-')
	}
	builder.WriteString(fmt.Sprintf("%02d", digit))

	return builder.String()
}
