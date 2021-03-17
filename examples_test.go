package trie

import "fmt"

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
		"red":    "\033[1;31m%s\033[0m",
		"blue":   "\033[1;36m%s\033[0m",
		"green":  "\033[1;32m%s\033[0m",
		"yellow": "\033[1;33m%s\033[0m",
	})
	inputs := []string{
		"green_apple",
		"yellow_banana",
		"red_strawberry",
		"blue_whale",
		"noprefixnocolor",
	}
	for _, inp := range inputs {
		format := "%s"

		if raw, prefixLen, ok := tr.SearchPrefixInString(inp); ok {
			format = raw.(string)
			inp = inp[prefixLen+1:]
		}

		fmt.Printf(format+"\n", inp)
	}
	// Output:
	// [1;32mapple[0m
	// [1;33mbanana[0m
	// [1;31mstrawberry[0m
	// [1;36mwhale[0m
	// noprefixnocolor
}
