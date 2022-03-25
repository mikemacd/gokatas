package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"unicode"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("No word list file supplied as an argument\n\n")
		os.Exit(0)
	}
	fileName := os.Args[1]

	file, err := os.Open(fileName) // For read access.
	if err != nil {
		fmt.Printf("error reading file: %v\n\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// the anagram dictionary that we will build will store all of the words that have the same set of
	// characters in the same 'bucket' using the letters stripped of accents and apostrophes as the key
	// eg:
	//	map[string][]string{
	//		"aars":[]string{"åsar"},
	//		"aepst":[]string{"paste", "pates", "peats", "septa", "spate", "tapes", "tepas"}
	//	}
	dictionary := map[string][]string{}

	longestWord := ""
	secondLongestWord := ""
	mostAnagrams := 0
	mostAnagramWord := ""
	secondMostAnagramWord := ""

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Grab a word from the input file and cast it to lowercase
		originalWord := strings.ToLower(scanner.Text())

		// strip any apostrophes
		word := strings.ReplaceAll(originalWord, "'", "")

		// map accented words to their unaccented equivalents
		b := make([]byte, len(word))
		t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
		_, _, e := t.Transform(b, []byte(word), true)
		if e != nil {
			panic(e)
		}
		word = strings.ToLower(strings.ReplaceAll(string(b), "\x00", ""))

		// sort the letters of the word to use it as the key to our dictionary
		sortedWord := SortString(word)
		if _, ok := dictionary[sortedWord]; !ok {
			dictionary[sortedWord] = []string{}
		}

		// add the original word to the appropriate anagram bucket
		dictionary[sortedWord] = append(dictionary[sortedWord], originalWord)

		// while we are looping through the file keep track of the longest word
		if len(sortedWord) > len(longestWord) {
			secondLongestWord = longestWord
			longestWord = sortedWord
		}

		// keep track of which bucket is the largest
		if len(dictionary[sortedWord]) > mostAnagrams {
			secondMostAnagramWord = mostAnagramWord
			mostAnagramWord = sortedWord
			mostAnagrams = len(dictionary[sortedWord])
		}

	}

	fmt.Printf("Longest Word anagrams: %v\n", dictionary[longestWord])
	fmt.Printf("2nd Longest Word anagrams: %v\n", dictionary[secondLongestWord])
	fmt.Println()
	fmt.Printf("Most Anagrams: %#v\n", dictionary[mostAnagramWord])
	fmt.Printf("2nd Most Anagrams: %#v\n", dictionary[secondMostAnagramWord])
	// fmt.Printf("\n\ndictionary:\n\n%#v\n\n", dictionary)

}

func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
}

func SortString(w string) string {
	s := strings.Split(w, "")
	sort.Strings(s)
	return strings.Join(s, "")
}
