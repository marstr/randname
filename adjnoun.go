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

	"github.com/marstr/collection"
)

// AdjNounFormat is a type of function that aggregates an Adjective, Noun, and Digit into a single formatted string.
type AdjNounFormat func(string, string, int) string

// AdjNoun creates a random name of the form adjectiveNameDigit
type AdjNoun struct {
	Adjectives    *collection.Dictionary
	Format        AdjNounFormat
	Nouns         *collection.Dictionary
	RandGenerator io.Reader
}

var defaultFormat AdjNounFormat = FormatPascalCaseAdjNoun
var defaultRandGenerator io.Reader = rand.Reader

// NewAdjNoun creates a new instance of AdjNoun that is populated with all of the defaults.
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

// Generate fetches a random name of the default dictionary and format.
func Generate() string {
	return AdjNoun{}.Generate()
}

// GeneratePascalCase fetches a random name using the default dictionary, formatted like: BigCloud9
func GeneratePascalCase() string {
	return AdjNoun{}.GeneratePascalCase()
}

// GenerateCamelCase fetches a random name using the default dictionary, formatted like: bigCloud9
func GenerateCamelCase() string {
	return AdjNoun{}.GenerateCamelCase()
}

// GenerateHyphened fetches a random name using the default dictionary, formatted like: big-cloud-9
func GenerateHyphened() string {
	return AdjNoun{}.GenerateHyphened()
}

// Generate creates a new randomly generated name with the
func (adNoun AdjNoun) Generate() string {
	if adNoun.Format == nil {
		adNoun.Format = defaultFormat
	}
	return adNoun.Format(adNoun.getAdjective(), adNoun.getNoun(), adNoun.getDigit())
}

// GeneratePascalCase randomly generates an adjective, noun, and number then returns them
// as a string formatted like: BigCloud9
func (adNoun AdjNoun) GeneratePascalCase() string {
	return FormatPascalCaseAdjNoun(adNoun.getAdjective(), adNoun.getNoun(), adNoun.getDigit())
}

// GenerateCamelCase randomly generates an adjective, noun, and number combination using the provided
// dictionaries, then returns them as a string formatted like: bigCloud9
func (adNoun AdjNoun) GenerateCamelCase() string {
	return FormatCamelCaseAdjNoun(adNoun.getAdjective(), adNoun.getNoun(), adNoun.getDigit())
}

// GenerateHyphened randomly generates an adjective, noun, and number combination using the provided
// dictionaries, then returns them as a string formatted like: big-cloud-9
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

// FormatCamelCaseAdjNoun formats an adjective, noun, and digit in the following way: bigCloud9
func FormatCamelCaseAdjNoun(adjective, noun string, digit int) string {
	pascal := FormatPascalCaseAdjNoun(adjective, noun, digit)
	return strings.ToLower(pascal[:1]) + pascal[1:]
}

// FormatPascalCaseAdjNoun formats an adjective, noun, and digit in the following way: BigCloud9
func FormatPascalCaseAdjNoun(adjective, noun string, digit int) string {
	builder := bytes.Buffer{}
	if len(adjective) > 0 {
		builder.WriteString(strings.ToUpper(adjective[:1]))
		builder.WriteString(strings.ToLower(adjective[1:]))
	}

	if len(noun) > 0 {
		builder.WriteString(strings.ToUpper(noun[:1]))
		builder.WriteString(strings.ToLower(noun[1:]))
	}

	builder.WriteString(fmt.Sprintf("%02d", digit))

	return builder.String()
}

// FormatHyphenedAdjNoun formats an adjective, noun, and digit in the following way: big-cloud-9
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
