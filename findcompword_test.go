package main

import (
	"testing"
)

func TestContains(t *testing.T) {
	cases := []struct {
		ss   []string
		s    string
		want bool
	}{
		{[]string{"a", "b"}, "a", true},
		{[]string{"a", "b"}, "bb", false},
		{[]string{"a", "c"}, "b", false},
		{[]string{"c"}, "c", true},
		{[]string{"a"}, "", false},
		{[]string{""}, "a", false},
		{[]string{""}, "", true},
		{[]string{}, "a", false},
		{[]string{}, "", false},
	}

	for _, c := range cases {
		got := contains(c.ss, c.s)
		if got != c.want {
			t.Errorf("contains(%q, %q) == %v, want %v", c.ss, c.s, got, c.want)
		}
	}
}

func TestCompound(t *testing.T) {
	cases := []struct {
		words []string
		word  string
		want  bool
	}{
		{[]string{"a"}, "a", false},
		{[]string{"a", "b"}, "a", false},
		{[]string{"a", "ab", "b"}, "ab", true},
		{[]string{"a", "ab", "c"}, "ab", false},
		// Compound-word assembled from 3 sub-words:
		{[]string{"a", "abc", "b", "c"}, "abc", true},
		{[]string{"a", "ab", "abc", "c"}, "abc", true},
		{[]string{"a", "ab", "abc", "ababc"}, "ababc", true},
	}

	for _, c := range cases {
		got := compound(c.words, c.word)
		if got != c.want {
			t.Errorf("compound(%q, %q) == %v, want %v", c.words, c.word, got, c.want)
		}
	}
}

func TestFindLongest(t *testing.T) {
	cases := []struct {
		words []string
		want  string
	}{
		{[]string{"a"}, ""},
		{[]string{"a", "b"}, ""},
		{[]string{"a", "ab", "b"}, "ab"},
		{[]string{"a", "ab", "c"}, ""},
		{[]string{"a", "abc", "b", "c"}, "abc"},
		{[]string{"a", "ab", "abc", "c"}, "abc"},
		{[]string{"a", "ab", "abc", "ababc", "c"}, "ababc"},
		// Two longest: "ab" and "bc", first in alphabetical order will be returned
		{[]string{"a", "ab", "b", "bc", "c"}, "ab"},
		// Following tests words with non-english alphabet
		{[]string{"世", "世界", "界"}, "世界"},
		{[]string{"a", "b", "a世界", "世", "世界", "界"}, "a世界"},
		{[]string{"a", "b", "a世界b", "世", "世界", "界"}, "a世界b"},
	}

	for _, c := range cases {
		got := findLongest(c.words)
		if got != c.want {
			t.Errorf("findLongest(%q) == %v, want %v", c.words, got, c.want)
		}
	}
}

func TestFindLongestParal(t *testing.T) {
	cases := []struct {
		words []string
		want  string
	}{
		{[]string{"a"}, ""},
		{[]string{"a", "b"}, ""},
		{[]string{"a", "ab", "b"}, "ab"},
		{[]string{"a", "ab", "c"}, ""},
		{[]string{"a", "abc", "b", "c"}, "abc"},
		{[]string{"a", "ab", "abc", "c"}, "abc"},
		{[]string{"a", "ab", "abc", "ababc", "c"}, "ababc"},
		// Following tests words with non-english alphabet
		{[]string{"世", "世界", "界"}, "世界"},
		{[]string{"a", "b", "a世界", "世", "世界", "界"}, "a世界"},
		{[]string{"a", "b", "a世界b", "世", "世界", "界"}, "a世界b"},
	}

	for _, c := range cases {
		got := findLongestParal(c.words)
		if got != c.want {
			t.Errorf("findLongestParal(%q) == %v, want %v", c.words, got, c.want)
		}
	}
}

var words, _ = readLines(*src)

func BenchmarkFindLongest(b *testing.B) {
	*parallel = false
	for i := 0; i < b.N; i++ {
		findLongest(words)
	}
}

func BenchmarkFindLongestParal(b *testing.B) {
	*parallel = true
	for i := 0; i < b.N; i++ {
		findLongestParal(words)
	}
}
