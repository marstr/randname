//go:build go1.7
// +build go1.7

package randname

import (
	"bufio"
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

func TestAdjNoun_Generate(t *testing.T) {
	testCases := []struct {
		Format   AdjNounFormat
		Expected *regexp.Regexp
	}{
		{FormatCamelCaseAdjNoun, camelForm},
		{FormatPascalCaseAdjNoun, pascalForm},
	}

	subject := AdjNoun{}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			subject.Format = tc.Format

			result := subject.Generate()
			if !tc.Expected.MatchString(result) {
				t.Fail()
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

func BenchmarkAdjNoun_Generate(b *testing.B) {
	subject := &AdjNoun{}
	subject.Generate() // Prime the pump to separate default Dictionary population time.
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		subject.Generate()
	}
}
