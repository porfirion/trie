package trie

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
)

func TestTrie_Add(t *testing.T) {
	tests := []struct {
		prefix []byte
		val    interface{}
	}{
		{[]byte(`👨‍`), `👨‍`},
		{[]byte(`👨‍🔧`), `👨‍🔧`},
		{[]byte(`👨`), `👨`},                           // 0xF0, 0x9F, 0x91, 0xA8
		{[]byte{0xF0, 0x9F, 0x91, 0x10}, "modified"}, // modified last byte of first key
		{[]byte{0xF0, 0x9F, 0x91, 0x10}, "modified"}, // add the same
		{[]byte{0xF0, 0x9F, 0x91}, "short"},          // add the same
	}

	trie := &Trie{}
	for _, tt := range tests {
		trie.Add(tt.prefix, tt.val)
	}

	var expected = &Trie{[]byte{0xF0, 0x9F, 0x91}, "short", &[256]*Trie{
		0x10: {[]byte{0x10}, "modified", nil},
		0xA8: {[]byte{0xA8}, `👨`, &[256]*Trie{
			0xE2: {[]byte{0xE2, 0x80, 0x8D}, `👨‍`, &[256]*Trie{
				0xF0: {[]byte{0xF0, 0x9F, 0x94, 0xA7}, `👨‍🔧`, nil},
			}},
		}},
	}}
	if !reflect.DeepEqual(trie, expected) {
		t.Fatalf("Not equal:\nexpected\n%s\ngot\n%s\n", expected, trie)
	}
}

func BenchmarkTrie_Add(b *testing.B) {
	b.ReportAllocs()

	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	randomString := func() []byte {
		var b = make([]byte, rand.Intn(10)+10)
		for i := range b {
			b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
		}
		return b[:]
	}

	tr := &Trie{}
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		// do not count random string allocations
		str := randomString()
		b.StartTimer()
		tr.Add(str, struct{}{})
	}
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
