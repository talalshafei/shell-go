package shellparser

import (
	"reflect"
	"testing"
)

func TestParseInput(t *testing.T) {

	t.Run("ParseInput should separate by whitespaces", func(t *testing.T) {
		input := "word0 word1 word2 word3"
		want := []string{"word0", "word1", "word2", "word3"}
		parser := NewParser()
		got, err := parser.Parse([]byte(input))

		assertNoError(t, err)
		assertParsedStrings(t, want, got)
	})

	t.Run("ParseInput should handle enclosed with single quotes as one string", func(t *testing.T) {
		input := "cat ~/Desktop/newfolder/'tmp file'"
		want := []string{"cat", "~/Desktop/newfolder/tmp file"}
		parser := NewParser()

		got, err := parser.Parse([]byte(input))
		assertNoError(t, err)
		assertParsedStrings(t, want, got)
	})

}

func assertParsedStrings(t testing.TB, want, got []string) {
	t.Helper()
	if !reflect.DeepEqual(want, got) {
		t.Errorf("Wanted %v, Got %v", want, got)
	}
}

func assertNoError(t testing.TB, err error) {
	if err != nil {
		t.Error(err.Error())
	}
}
