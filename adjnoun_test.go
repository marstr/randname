package randname

import (
	"regexp"
	"testing"
)

var pascalForm = regexp.MustCompile(`^[A-Z][a-z]+[A-Z][a-z]+\d{2}$`)
var camelForm = regexp.MustCompile(`^[a-z]+[A-Z][a-z]+\d{2}$`)

func TestAdjNound_Generate(t *testing.T) {
	testCases := []struct {
		Format   AdjNounFormat
		Expected *regexp.Regexp
	}{
		{GenerateCamelCaseAdjNoun, camelForm},
		{GeneratePascalCaseAdjNoun, pascalForm},
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
