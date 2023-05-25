package export

import (
	"fmt"
	"testing"

	"github.com/porfirion/trie"
)

type T = trie.Trie[string]

var stringsTrieExample = &T{Prefix: []byte{0xF0, 0x9F, 0x91}, Value: ptr("short"), Children: &[256]*T{
	0x10: {Prefix: []byte{0x10}, Value: ptr("modified")},
	0xA8: {Prefix: []byte{0xA8}, Value: ptr("nokey"), Children: &[256]*T{
		0xE2: {Prefix: []byte{0xE2, 0x80, 0x8D}, Value: ptr("withsep"), Children: &[256]*T{
			0xF0: {Prefix: []byte{0xF0, 0x9F, 0x94, 0xA7}, Value: ptr("withkey")},
		}},
	}},
}}

func ExampleExport() {
	var res = Export(stringsTrieExample, ExportSettings{
		Padding:   "    ",
		TrieAlias: "StrTrie", // says to replace type Trie with alias (can be defined like type T = trie.Trie)
	})
	fmt.Print(res)
	// Output:
	// {Prefix: []byte{0xF0, 0x9F, 0x91}, Value: ptr[string]("short"), Children: &[256]*StrTrie{
	//     0x10: {Prefix: []byte{0x10}, Value: ptr[string]("modified")},
	//     0xA8: {Prefix: []byte{0xA8}, Value: ptr[string]("nokey"), Children: &[256]*StrTrie{
	//         0xE2: {Prefix: []byte{0xE2, 0x80, 0x8D}, Value: ptr[string]("withsep"), Children: &[256]*StrTrie{
	//             0xF0: {Prefix: []byte{0xF0, 0x9F, 0x94, 0xA7}, Value: ptr[string]("withkey")},
	//         }},
	//     }},
	// }}
}

var uint8Trie = trie.BuildFromMap(map[string]uint8{
	"one": 1,
	"two": 2,
})

func ExampleExport_basic() {
	fmt.Println(Export(uint8Trie, ExportSettings{TrieAlias: "U8Trie", Padding: "  "}))
	// Output:
	// {Children: &[256]*U8Trie{
	//   0x6F: {Prefix: []byte{0x6F, 0x6E, 0x65}, Value: ptr[uint8](1)},
	//   0x74: {Prefix: []byte{0x74, 0x77, 0x6F}, Value: ptr[uint8](2)},
	// }}
}

var anyTrieExample = trie.BuildFromMap(map[string]interface{}{
	"float":       31.7,
	"float.round": 32.0,
	"int":         16,
	"bool":        true,
	"uint":        uint(15),
	"uint64":      uint64(21),
	"uint32":      uint32(20),
	"string":      "teststr",
	"bytes":       [...]byte{},
})

func ExampleExport_withDifferentTypes() {
	var res = Export(anyTrieExample, ExportSettings{TrieAlias: "AnyTrie", Padding: "    "})

	fmt.Print(res)
	// Output:
	// {Children: &[256]*AnyTrie{
	//     0x62: {Prefix: []byte{0x62}, Children: &[256]*AnyTrie{
	//         0x6F: {Prefix: []byte{0x6F, 0x6F, 0x6C}, Value: ptr[any](true)},
	//         0x79: {Prefix: []byte{0x79, 0x74, 0x65, 0x73}, Value: ptr[any](([0]uint8)([]))},
	//     }},
	//     0x66: {Prefix: []byte{0x66, 0x6C, 0x6F, 0x61, 0x74}, Value: ptr[any]((float64)(31.7)), Children: &[256]*AnyTrie{
	//         0x2E: {Prefix: []byte{0x2E, 0x72, 0x6F, 0x75, 0x6E, 0x64}, Value: ptr[any]((float64)(32))},
	//     }},
	//     0x69: {Prefix: []byte{0x69, 0x6E, 0x74}, Value: ptr[any](16)},
	//     0x73: {Prefix: []byte{0x73, 0x74, 0x72, 0x69, 0x6E, 0x67}, Value: ptr[any]("teststr")},
	//     0x75: {Prefix: []byte{0x75, 0x69, 0x6E, 0x74}, Value: ptr[any]((uint)(15)), Children: &[256]*AnyTrie{
	//         0x33: {Prefix: []byte{0x33, 0x32}, Value: ptr[any]((uint32)(20))},
	//         0x36: {Prefix: []byte{0x36, 0x34}, Value: ptr[any]((uint64)(21))},
	//     }},
	// }}
}

// BenchmarkExport-4   	  114370	     10211 ns/op	   11754 B/op	     117 allocs/op
// BenchmarkExport-4   	  101808	     10075 ns/op	   11714 B/op	     113 allocs/op - Grow() in exportBytes
// BenchmarkExport-4   	  119307	     10167 ns/op	   12322 B/op	     109 allocs/op - Grow() in export itself
// BenchmarkExport-4   	  114774	     10012 ns/op	   12242 B/op	     105 allocs/op - prealloc for children
// BenchmarkExport-4   	  128707	      9210 ns/op	   11810 B/op	      92 allocs/op - prealloc for fields
// BenchmarkExport-4   	  129270	      9116 ns/op	   11682 B/op	      84 allocs/op - replace sprintf with concat
// BenchmarkExport-4   	  124080	      8747 ns/op	   11602 B/op	      79 allocs/op - concat for prefix
// BenchmarkExport-4   	  123142	      9854 ns/op	   11522 B/op	      74 allocs/op - concat for value
// BenchmarkExport-4   	  116490	      9519 ns/op	   11715 B/op	      72 allocs/op - bytesRep for formatBytes
// BenchmarkExport-4   	  189608	      6083 ns/op	   11201 B/op	      59 allocs/op - bytesRep for ind & concat child
// BenchmarkExport-4   	  225652	      5198 ns/op	   11120 B/op	      54 allocs/op - concat for string literal ("%s")
// BenchmarkExport-4   	  299611	      3911 ns/op	    4768 B/op	      42 allocs/op - create example before loop %)
// BenchmarkExport-8   	  129058	     13517 ns/op	    5632 B/op	      62 allocs/op - generics :(
func BenchmarkExport(b *testing.B) {
	b.ReportAllocs()

	example := &T{Prefix: []byte{0xF0, 0x9F, 0x91}, Value: ptr("short"), Children: &[256]*T{
		0x10: {Prefix: []byte{0x10}, Value: ptr("modified")},
		0xA8: {Prefix: []byte{0xA8}, Value: ptr("nokey"), Children: &[256]*T{
			0xE2: {Prefix: []byte{0xE2, 0x80, 0x8D}, Value: ptr("withsep"), Children: &[256]*T{
				0xF0: {Prefix: []byte{0xF0, 0x9F, 0x94, 0xA7}, Value: ptr("withkey")},
			}},
		}},
	}}
	settings := ExportSettings{
		Padding:   "    ",
		TrieAlias: "T",
	}

	for i := 0; i < b.N; i++ {
		_ = Export(example, settings)
	}
}

func ptr[T any](v T) *T {
	return &v
}
