package main

import (
	"fmt"
	"reflect"
	"slices"
)

func main() {
	l := []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "стол"}
	expected := map[string][]string{
		"пятак":  {"пятак", "пятка", "тяпка"},
		"листок": {"листок", "слиток", "столик"},
	}
	if !reflect.DeepEqual(FindAnagrams(l), expected) {
		panic("incorrect result")
	}
	fmt.Println("all good")
}

type anagram [33]int8

// FindAnagrams time complexity: O(n * m * log n)
func FindAnagrams(words []string) map[string][]string {
	freq := make(map[anagram]int)
	anagrams := make(map[string]anagram)
	firstFound := make(map[anagram]string)
	for _, word := range words {
		tmp := getAnagram(word)
		anagrams[word] = tmp
		freq[tmp]++
		if _, ok := firstFound[tmp]; !ok {
			firstFound[tmp] = word
		}
	}
	slices.Sort(words)
	out := make(map[string][]string)
	for _, word := range words {
		if freq[anagrams[word]] == 1 {
			continue
		}

		out[firstFound[anagrams[word]]] = append(out[firstFound[anagrams[word]]], word)
	}
	return out
}

func getAnagram(s string) (out anagram) {
	for _, ch := range s {
		out[ch-'а']++
	}
	return
}
