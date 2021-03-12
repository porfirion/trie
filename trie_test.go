package trie

import (
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

	tr := &Trie{}
	for _, tt := range tests {
		tr.Add(tt.prefix, tt.val)
	}

	var expected = &Trie{[]byte{0xF0, 0x9F, 0x91}, "short", &[256]*Trie{
		0x10: {[]byte{0x10}, "modified", nil},
		0xA8: {[]byte{0xA8}, `👨`, &[256]*Trie{
			0xE2: {[]byte{0xE2, 0x80, 0x8D}, `👨‍`, &[256]*Trie{
				0xF0: {[]byte{0xF0, 0x9F, 0x94, 0xA7}, `👨‍🔧`, nil},
			}},
		}},
	}}
	if !reflect.DeepEqual(tr, expected) {
		t.Fatalf("Not equal:\nexpected\n%s\ngot\n%s\n", expected, tr)
	}
}

type EqFunc func()

func (f EqFunc) Equal(v interface{}) bool {
	return true
}

var _ Equaler = EqFunc(nil)

// it's ok to add funcs if they implement Equaler
func TestTrie_AddFuncTwice(t *testing.T) {
	tr := &Trie{}
	tr.AddString("foo", EqFunc(func() {}))
	tr.AddString("foo", EqFunc(func() {}))
}

func BenchmarkTrie_Add(b *testing.B) {
	b.ReportAllocs()

	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	randomString := func() []byte {
		var b = make([]byte, 32)
		for i := range b {
			b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
		}
		return b
	}

	tr := &Trie{}
	for i := 0; i < b.N; i++ {
		// one allocation for random string generation
		tr.Add(randomString(), struct{}{})

		// This variant shows only 1 allocation,
		// but stopping and starting timer is very slow - benchmark can last for 30 seconds!
		//b.StopTimer()
		//str := randomString()
		//b.StartTimer()
		//tr.Add(str, struct{}{})
	}
}
