package randname_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/marstr/randname"
)

func ExamplePrefixed_Generate() {
	generator := randname.Prefixed{
		Prefix: "randnameRocks",
	}
	fmt.Println(generator.Generate())
}

func ExamplePrefixed() {
	generator := randname.Prefixed{
		Prefix:     "randname-",
		Acceptable: randname.ArabicNumerals,
		Len:        10,
	}

	for i := 0; i < 10; i++ {
		fmt.Println(generator.Generate())
	}
}

func ExampleGenerateWithPrefix() {
	for i := 0; i < 10; i++ {
		fmt.Println(randname.GenerateWithPrefix("randname-", 10))
	}
}

func TestPrefixed_Generate_NoRepeats(t *testing.T) {
	const trials = 100

	seen := make(map[string]struct{}, trials)

	subject := randname.Prefixed{
		Len: uint8(40),
	}

	for i := 0; i < trials; i++ {
		tuple := subject.Generate()
		t.Log("Generated: ", tuple)

		if _, ok := seen[tuple]; ok {
			t.Logf("previously generated value, %q encountered again.", tuple)
			t.Fail()
		}
	}
}

func TestPrefixed_Generate(t *testing.T) {
	testCases := []randname.Prefixed{
		randname.Prefixed{},
		randname.Prefixed{
			Prefix: "randname-",
		},
		randname.Prefixed{
			Len: 1,
		},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			got := tc.Generate()

			if !strings.HasPrefix(got, tc.Prefix) {
				t.Logf("%q does not have expected prefix %q", got, tc.Prefix)
				t.Fail()
			}

			expectedRunes := tc.Len
			if expectedRunes == 0 {
				expectedRunes = randname.PrefixedDefaultLen
			}

			trimmed := strings.TrimPrefix(got, tc.Prefix)

			if len(trimmed) != int(expectedRunes) {
				t.Logf("%q does not have the expected number of random runes: %d", got, expectedRunes)
				t.Fail()
			}

			if tc.Acceptable == nil {
				tc.Acceptable = randname.PrefixedDefaultAcceptable
			}

			for _, r := range trimmed {
				if !strings.ContainsRune(string(tc.Acceptable), r) {
					t.Logf("%q contains unacceptable rune: %s", got, string(r))
					t.Fail()
				}
			}
		})
	}
}
