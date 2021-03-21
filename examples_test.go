package trie

import (
	"fmt"
	"strings"
	"unicode"
)

func ExampleBuildPrefixesOnly() {
	prefixes := BuildPrefixesOnly("tiny_", "small_", "normal_", "large_")

	var myString = "large_banana"

	if prefix, ok := prefixes.TakePrefix(myString); ok {
		fmt.Printf("prefix \"%s\" found\n", prefix)
	} else {
		fmt.Println("no prefix found")
	}
	// Output:
	//prefix "large_" found
}

func ExampleTrie_String() {
	example := &T{Prefix: []byte{0xF0, 0x9F, 0x91}, Value: "short", Children: &[256]*T{
		0x10: {Prefix: []byte{0x10}, Value: "modified"},
		0xA8: {Prefix: []byte{0xA8}, Value: "nokey", Children: &[256]*T{
			0xE2: {Prefix: []byte{0xE2, 0x80, 0x8D}, Value: "withsep", Children: &[256]*T{
				0xF0: {Prefix: []byte{0xF0, 0x9F, 0x94, 0xA7}, Value: "withkey"},
			}},
		}},
	}}
	fmt.Println(example)
	// Output:
	// [F0 9F 91] short
	// ├─10─ [10] modified
	// ├─A8─ [A8] nokey
	// │     ├─E2─ [E2 80 8D] withsep
	// │     │     ├─F0─ [F0 9F 94 A7] withkey
}

func ExampleTrie_Iterate() {
	example := &T{Prefix: []byte{0xF0, 0x9F, 0x91}, Value: "short", Children: &[256]*T{
		0x10: {Prefix: []byte{0x10}, Value: "modified"},
		0xA8: {Prefix: []byte{0xA8}, Value: "nokey", Children: &[256]*T{
			0xE2: {Prefix: []byte{0xE2, 0x80, 0x8D}, Value: "withsep", Children: &[256]*T{
				0xF0: {Prefix: []byte{0xF0, 0x9F, 0x94, 0xA7}, Value: "withkey"},
			}},
		}},
	}}
	example.Iterate(func(prefix []byte, value interface{}) {
		fmt.Printf("[%v] %+v\n", prefix, value)
	})
	// Output:
	// [[240 159 145]] short
	// [[240 159 145 16]] modified
	// [[240 159 145 168]] nokey
	// [[240 159 145 168 226 128 141]] withsep
	// [[240 159 145 168 226 128 141 240 159 148 167]] withkey
}

func ExampleTrie_SearchPrefixInString() {
	tr := BuildFromMap(map[string]interface{}{
		"lower": func(s string) string { return strings.ToLower(s) },
		"upper": func(s string) string { return strings.ToUpper(s) },
		"snake": func(s string) string {
			var res = make([]rune, 0, len(s))
			for i, r := range s {
				if unicode.IsUpper(r) && i > 0 {
					res = append(res, '_', unicode.ToLower(r))
				} else {
					res = append(res, unicode.ToLower(r))
				}
			}
			return string(res)
		},
		"camel": func(s string) string {
			var res = make([]rune, 0, len(s))
			toUpper := false
			for _, r := range s {
				if r == '_' {
					toUpper = true
					continue
				} else if toUpper {
					res = append(res, unicode.ToUpper(r))
					toUpper = false
				} else {
					res = append(res, unicode.ToLower(r))
				}
			}
			return string(res)
		},
	})
	inputs := []string{
		"upperapple",
		"lowerBANANA",
		"cameltest_me",
		"snakeAnotherExample",
		"noprefixString",
	}
	for _, inp := range inputs {
		if raw, prefixLen, ok := tr.SearchPrefixInString(inp); ok {
			format := raw.(func(string) string)
			fmt.Println(format(inp[prefixLen:]))
		} else {
			fmt.Printf("no prefix found in \"%s\"\n", inp)
		}
	}

	// Output:
	//APPLE
	//banana
	//testMe
	//another_example
	//no prefix found in "noprefixString"
}
