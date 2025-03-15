package shellparser

import (
	"reflect"
	"testing"
)

func TestParseInput(t *testing.T) {

	t.Run("Parse should separate by whitespaces", func(t *testing.T) {
		input := "word0 word1 word2 word3"
		want := []string{"word0", "word1", "word2", "word3"}
		parser := NewParser()
		got, err := parser.Parse([]byte(input))

		assertNoError(t, err)
		assertParsedStrings(t, want, got)
	})

	t.Run("Parse should handle enclosed with single quotes as one string", func(t *testing.T) {
		input := "cat ~/Desktop/newfolder/'tmp file'"
		want := []string{"cat", "~/Desktop/newfolder/tmp file"}
		parser := NewParser()

		got, err := parser.Parse([]byte(input))
		assertNoError(t, err)
		assertParsedStrings(t, want, got)
	})
}

func TestParseRedirectionOperators(t *testing.T) {
	t.Run("Parse should handle > operator", func(t *testing.T) {
		table := []struct {
			input string
			want  []string
		}{
			{"echo Hello > file", []string{"echo", "Hello", "1>", "file"}},
			{"echo Hello>file", []string{"echo", "Hello", "1>", "file"}},
			{"echo Hello1 > file", []string{"echo", "Hello1", "1>", "file"}},
			{"echo Hello1 1> file", []string{"echo", "Hello1", "1>", "file"}},
			{"echo Hello >file", []string{"echo", "Hello", "1>", "file"}},
			{"echo Hello2 > file", []string{"echo", "Hello2", "1>", "file"}},
			{"echo Hello2>file", []string{"echo", "Hello2", "1>", "file"}},
			{"echo Hello 2> file", []string{"echo", "Hello", "2>", "file"}},
			{"command 1>log.txt", []string{"command", "1>", "log.txt"}},
			{"ls > output.txt", []string{"ls", "1>", "output.txt"}},
		}

		parser := NewParser()

		for _, entry := range table {
			t.Run(entry.input, func(t *testing.T) {
				got, err := parser.Parse([]byte(entry.input))
				assertNoError(t, err)
				assertParsedStrings(t, entry.want, got)
			})
		}
	})

	t.Run("Parse should handle < operator", func(t *testing.T) {
		table := []struct {
			input string
			want  []string
		}{
			{"cat < input.txt", []string{"cat", "0<", "input.txt"}},
			{"sort <data.txt", []string{"sort", "0<", "data.txt"}},
			{"command 2<input.txt", []string{"command", "0<", "input.txt"}},
			{"grep pattern <file.txt", []string{"grep", "pattern", "0<", "file.txt"}},
			{"awk '{print $1}' <data.txt", []string{"awk", "{print $1}", "0<", "data.txt"}},
		}

		parser := NewParser()

		for _, entry := range table {
			t.Run(entry.input, func(t *testing.T) {
				got, err := parser.Parse([]byte(entry.input))
				assertNoError(t, err)
				assertParsedStrings(t, entry.want, got)
			})
		}
	})

	t.Run("Parse should handle >> operator", func(t *testing.T) {
		table := []struct {
			input string
			want  []string
		}{
			{"echo Hello >> file", []string{"echo", "Hello", "1>>", "file"}},
			{"echo Hello>>file", []string{"echo", "Hello", "1>>", "file"}},
			{"command 2>>log.txt", []string{"command", "2>>", "log.txt"}},
			{"echo Append this >>output.txt", []string{"echo", "Append", "this", "1>>", "output.txt"}},
			{"ls -l >> dirs.txt", []string{"ls", "-l", "1>>", "dirs.txt"}},
		}

		parser := NewParser()

		for _, entry := range table {
			t.Run(entry.input, func(t *testing.T) {
				got, err := parser.Parse([]byte(entry.input))
				assertNoError(t, err)
				assertParsedStrings(t, entry.want, got)
			})
		}
	})

	t.Run("Should raise unexpected token error", func(t *testing.T) {
		table := []string{">", "1>", "2>", ">>", "1>>", "2>>", "<", "0<"}
		parser := NewParser()
		for _, entry := range table {
			t.Run(entry, func(t *testing.T) {
				got, err := parser.Parse([]byte(entry))
				if err == nil {
					t.Errorf("%s alone should raise an error", entry)
				}
				if len(got) > 0 {
					t.Errorf("got '%v' should be empty", got)
				}
			})
		}
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
