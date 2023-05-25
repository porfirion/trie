package trie

import (
	"fmt"
	"strings"
)

// Also can be created with
//
//	prefixes := &trie.Trie{}
//	prefixes.PutString("one", 1)
//	prefixes.PutString("two", 2)
//	prefixes.PutString("three", 3)
//	prefixes.PutString("", 0)
var prefixes = BuildFromMap(map[string]int{
	"one":   1,
	"two":   2,
	"three": 3,
	"":      0,
})

func Example_prefixes() {
	var inputs = []string{
		"twoApple",
		"oneBanana",
		"Carrot",
	}

	for _, inp := range inputs {
		if val, prefixLen, ok := prefixes.SearchPrefixInString(inp); ok {
			fmt.Println(strings.Repeat(inp[prefixLen:], val))
		}
	}

	// Output:
	// AppleApple
	// Banana
	//
}
