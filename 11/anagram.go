package main

import (
	"fmt"
	"sort"
	"strings"
)

type sliceRune []rune

func (s sliceRune) Len() int {
	return len(s)
}
func (s sliceRune) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s sliceRune) Less(i, j int) bool {
	return s[i] < s[j]
}

func sortString(s string) string {
	r := []rune(s)
	sort.Sort(sliceRune(r))
	return string(r)
}

func findAnagram(s []string) map[string][]string {
	anagramMap := make(map[string][]string)
	firstWord := make(map[string]string)
	resMap := make(map[string][]string)
	var temp string
	for _, e := range s {
		temp = sortString(strings.ToLower(e))
		if _, ok := anagramMap[temp]; !ok {
			firstWord[temp] = e
		}
		anagramMap[temp] = append(anagramMap[temp], e)
	}
	for k, v := range anagramMap {
		if len(v) > 1 {
			sort.Strings(v)
			resMap[firstWord[k]] = v
		}
	}
	return resMap
}

func main() {
	s := []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "стол"}

	fmt.Println(findAnagram(s))

}
