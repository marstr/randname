//go:build go1.7
// +build go1.7

package randname

import (
	"bufio"
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sync"
	"testing"

	"github.com/marstr/collection"
)

var pascalForm = regexp.MustCompile(`^[A-Z][a-z]+[A-Z][a-z]+\d{2}$`)
var camelForm = regexp.MustCompile(`^[a-z]+[A-Z][a-z]+\d{2}$`)
var hyphenedForm = regexp.MustCompile(`^[a-z]+-[a-z]+-\d{2}$`)

func TestGenerateNamedFormats(t *testing.T) {
	tests := []struct {
		name     string
		generate func() string
		want     *regexp.Regexp
	}{
		{"pascal", GeneratePascalCase, pascalForm},
		{"camel", GenerateCamelCase, camelForm},
		{"hyphened", GenerateHyphened, hyphenedForm},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.generate(); !tt.want.MatchString(got) {
				t.Errorf("%q didn't match expected pattern %q", got, tt.want)
			}
		})
	}
}

func TestAdjNoun_GenerateNamedFormats(t *testing.T) {
	adjectives := &collection.Dictionary{}
	adjectives.Add("bIG")
	nouns := &collection.Dictionary{}
	nouns.Add("cLOUD")

	methods := []struct {
		name     string
		generate func(AdjNoun) string
		prefix   string
	}{
		{"pascal", AdjNoun.GeneratePascalCase, "BigCloud"},
		{"camel", AdjNoun.GenerateCamelCase, "bigCloud"},
		{"hyphened", AdjNoun.GenerateHyphened, "big-cloud-"},
	}
	digits := []struct {
		value byte
		want  string
	}{
		{0, "00"},
		{9, "09"},
		{99, "99"},
	}

	for _, method := range methods {
		t.Run(method.name, func(t *testing.T) {
			for _, configured := range []bool{false, true} {
				name := "default_formatter"
				if configured {
					name = "custom_formatter"
				}
				t.Run(name, func(t *testing.T) {
					for _, digit := range digits {
						t.Run(digit.want, func(t *testing.T) {
							generator := AdjNoun{
								Adjectives:    adjectives,
								Nouns:         nouns,
								RandGenerator: bytes.NewReader(bytes.Repeat([]byte{digit.value}, 16)),
							}
							if configured {
								generator.Format = func(string, string, int) string {
									t.Fatal("named generation method called the configured formatter")
									return ""
								}
							}
							if got, want := method.generate(generator), method.prefix+digit.want; got != want {
								t.Errorf("got %q, want %q", got, want)
							}
						})
					}
				})
			}
		})
	}
}

func TestAdjNoun_GenerateNamedFormatsOverrideFormat(t *testing.T) {
	adjectives := &collection.Dictionary{}
	adjectives.Add("big")
	nouns := &collection.Dictionary{}
	nouns.Add("cloud")

	tests := []struct {
		name     string
		generate func(AdjNoun) string
		want     string
	}{
		{"pascal", AdjNoun.GeneratePascalCase, "BigCloud09"},
		{"camel", AdjNoun.GenerateCamelCase, "bigCloud09"},
		{"hyphened", AdjNoun.GenerateHyphened, "big-cloud-09"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			generator := AdjNoun{
				Adjectives: adjectives,
				Nouns:      nouns,
				Format: func(string, string, int) string {
					return "my-custom-format"
				},
				RandGenerator: bytes.NewReader(bytes.Repeat([]byte{9}, 48)),
			}
			if got := generator.Generate(); got != "my-custom-format" {
				t.Fatalf("Generate() = %q, want custom format", got)
			}
			if got := tt.generate(generator); got != tt.want {
				t.Errorf("named method = %q, want %q", got, tt.want)
			}
			if got := generator.Generate(); got != "my-custom-format" {
				t.Errorf("Generate() after named method = %q, want custom format", got)
			}
		})
	}
}

func TestAdjNoun_Generate(t *testing.T) {
	testCases := []struct {
		Format   AdjNounFormat
		Expected *regexp.Regexp
	}{
		{nil, pascalForm},
		{FormatCamelCaseAdjNoun, camelForm},
		{FormatPascalCaseAdjNoun, pascalForm},
	}

	var subject AdjNoun

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			subject.Format = tc.Format

			result := subject.Generate()
			if !tc.Expected.MatchString(result) {
				t.Errorf("%q didn't match expected pattern %q", result, tc.Expected)
			}
		})
	}
}

func TestDefaultDictionaries(t *testing.T) {
	tests := []struct {
		file string
		get  func() *collection.Dictionary
	}{
		{"adjectives.txt", getDefaultAdjectives},
		{"nouns.txt", getDefaultNouns},
	}

	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {
			file, err := os.Open(tt.file)
			if err != nil {
				t.Fatal(err)
			}
			defer file.Close()

			dict := tt.get()
			scanner := bufio.NewScanner(file)
			var count int64
			for scanner.Scan() {
				word := scanner.Text()
				if !dict.Contains(word) {
					t.Errorf("missing word %q", word)
				}
				count++
			}
			if err := scanner.Err(); err != nil {
				t.Fatal(err)
			}
			if dict.Size() != count {
				t.Errorf("got %d words, want %d", dict.Size(), count)
			}
		})
	}
}

func TestDefaultBuilderLazyAndConcurrent(t *testing.T) {
	words := "first\nsecond"
	get := defaultBuilder(&words)
	if words == "" {
		t.Fatal("words loaded before first use")
	}
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			dict := get()
			if dict.Size() != 2 || !dict.Contains("first") || !dict.Contains("second") {
				t.Error("default dictionary was not populated")
			}
		}()
	}
	wg.Wait()
	if words != "" {
		t.Fatal("source string was not released")
	}
}

func TestDefaultDictionariesStandalone(t *testing.T) {
	if os.Getenv("RANDNAME_STANDALONE_TEST") == "1" {
		if !getDefaultAdjectives().Contains("able") || !getDefaultNouns().Contains("aardvark") {
			t.Fatal("embedded dictionaries were not loaded")
		}
		if getDefaultAdjectives().Contains("impostor") || getDefaultNouns().Contains("impostor") {
			t.Fatal("working directory dictionaries were loaded")
		}
		if name := Generate(); !pascalForm.MatchString(name) {
			t.Fatalf("unexpected generated name %q", name)
		}
		return
	}

	dir := t.TempDir()
	for _, name := range []string{"adjectives.txt", "nouns.txt"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("impostor\n"), 0600); err != nil {
			t.Fatal(err)
		}
	}
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(executable, "-test.run=^TestDefaultDictionariesStandalone$")
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "RANDNAME_STANDALONE_TEST=1")
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("standalone test failed: %v\n%s", err, output)
	}
}

func TestGenerateAdjNounFormats(t *testing.T) {
	tests := []struct {
		name      string
		adjective string
		noun      string
		digit     int
		pascal    string
		camel     string
		hyphened  string
	}{
		{"lowercase", "big", "cloud", 9, "BigCloud09", "bigCloud09", "big-cloud-09"},
		{"mixed_case", "bIG", "cLOUD", 42, "BigCloud42", "bigCloud42", "big-cloud-42"},
		{"single_letters", "B", "C", 0, "BC00", "bC00", "b-c-00"},
		{"empty_adjective", "", "cloud", 9, "Cloud09", "cloud09", "cloud-09"},
		{"empty_noun", "big", "", 9, "Big09", "big09", "big-09"},
		{"empty_words", "", "", 9, "09", "09", "09"},
		{"negative_digit", "big", "cloud", -1, "BigCloud-1", "bigCloud-1", "big-cloud--1"},
		{"wide_digit", "big", "cloud", 123, "BigCloud123", "bigCloud123", "big-cloud-123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formats := []struct {
				name   string
				format AdjNounFormat
				want   string
			}{
				{"pascal", FormatPascalCaseAdjNoun, tt.pascal},
				{"camel", FormatCamelCaseAdjNoun, tt.camel},
				{"hyphened", FormatHyphenedAdjNoun, tt.hyphened},
			}
			for _, format := range formats {
				t.Run(format.name, func(t *testing.T) {
					if got := format.format(tt.adjective, tt.noun, tt.digit); got != format.want {
						t.Errorf("got %q, want %q", got, format.want)
					}
				})
			}
		})
	}
}

func BenchmarkAdjNoun_Generate(b *testing.B) {
	subject := &AdjNoun{}
	subject.Generate() // Prime the pump to separate default Dictionary population time.
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = subject.Generate()
	}
}

func BenchmarkFormatCamelCaseAdjNoun(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = FormatCamelCaseAdjNoun("big", "cloud", 9)
	}
}

func BenchmarkFormatPascalCaseAdjNoun(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = FormatPascalCaseAdjNoun("big", "cloud", 9)
	}
}

func BenchmarkFormatHyphenedAdjNoun(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = FormatHyphenedAdjNoun("big", "cloud", 9)
	}
}
