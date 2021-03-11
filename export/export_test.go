package export

import (
	"fmt"

	"github.com/porfirion/trie"
)

type T = trie.Trie

func ExampleExport() {
	example := &T{Prefix: []byte{0xF0, 0x9F, 0x91}, Value: "short", Children: &[256]*T{
		0x10: {Prefix: []byte{0x10}, Value: "modified"},
		0xA8: {Prefix: []byte{0xA8}, Value: "nokey", Children: &[256]*T{
			0xE2: {Prefix: []byte{0xE2, 0x80, 0x8D}, Value: "withsep", Children: &[256]*T{
				0xF0: {Prefix: []byte{0xF0, 0x9F, 0x94, 0xA7}, Value: "withkey"},
			}},
		}},
	}}
	var res = Export(example, ExportSettings{
		Padding:   "    ",
		TrieAlias: "T", // says to replace type Trie with alias (can be defined like type T = trie.Trie)
	}, "")
	fmt.Print(res)
	// Output:
	// {Prefix: []byte{0xF0, 0x9F, 0x91}, Value: "short", Children: &[256]*T{
	//     0x10: {Prefix: []byte{0x10}, Value: "modified"},
	//     0xA8: {Prefix: []byte{0xA8}, Value: "nokey", Children: &[256]*T{
	//         0xE2: {Prefix: []byte{0xE2, 0x80, 0x8D}, Value: "withsep", Children: &[256]*T{
	//             0xF0: {Prefix: []byte{0xF0, 0x9F, 0x94, 0xA7}, Value: "withkey"},
	//         }},
	//     }},
	// }}
}

func ExampleExport_withDifferentTypes() {
	exampleTypes := trie.FromMap(map[string]interface{}{
		"float":       31.7,
		"float.round": 32.0,
		"int":         16,
		"bool":        true,
	})
	var res = Export(exampleTypes, ExportSettings{Padding: "    "}, "")

	fmt.Print(res)
	// Output:
	// {Children: &[256]*Trie{
	//     0x62: {Prefix: []byte{0x62, 0x6F, 0x6F, 0x6C}, Value: true},
	//     0x66: {Prefix: []byte{0x66, 0x6C, 0x6F, 0x61, 0x74}, Value: 31.7, Children: &[256]*Trie{
	//         0x2E: {Prefix: []byte{0x2E, 0x72, 0x6F, 0x75, 0x6E, 0x64}, Value: 32},
	//     }},
	//     0x69: {Prefix: []byte{0x69, 0x6E, 0x74}, Value: 16},
	// }}
}
