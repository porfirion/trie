package trie

import (
	"math/rand"
	"reflect"
	"strings"
	"testing"
	"unicode/utf8"
)

type T = Trie

var (
	tr = T{Prefix: []byte{0xF0, 0x9F, 0x91}, Children: &[256]*T{
		0xA8: {Prefix: []byte{0xA8}, Value: "👨" /*F0 9F 91 A8*/, Children: &[256]*T{
			0xE2: {Prefix: []byte{0xE2, 0x80, 0x8D}, Children: &[256]*T{
				0xE2: {Prefix: []byte{0xE2}, Children: &[256]*T{
					0x9A: {Prefix: []byte{0x9A}, Children: &[256]*T{
						0x95: {Prefix: []byte{0x95, 0xEF, 0xB8, 0x8F}, Value: "👨‍⚕️" /*F0 9F 91 A8 E2 80 8D E2 9A 95 EF B8 8F*/},
						0x96: {Prefix: []byte{0x96, 0xEF, 0xB8, 0x8F}, Value: "👨‍⚖️" /*F0 9F 91 A8 E2 80 8D E2 9A 96 EF B8 8F*/},
					}},
					0x9C: {Prefix: []byte{0x9C, 0x88, 0xEF, 0xB8, 0x8F}, Value: "👨‍✈️" /*F0 9F 91 A8 E2 80 8D E2 9C 88 EF B8 8F*/},
					0x9D: {Prefix: []byte{0x9D, 0xA4, 0xEF, 0xB8, 0x8F, 0xE2, 0x80, 0x8D, 0xF0, 0x9F}, Children: &[256]*T{
						0x91: {Prefix: []byte{0x91, 0xA8}, Value: "👨‍❤️‍👨" /*F0 9F 91 A8 E2 80 8D E2 9D A4 EF B8 8F E2 80 8D F0 9F 91 A8*/},
						0x92: {Prefix: []byte{0x92, 0x8B, 0xE2, 0x80, 0x8D, 0xF0, 0x9F, 0x91, 0xA8}, Value: "👨‍❤️‍💋‍👨" /*F0 9F 91 A8 E2 80 8D E2 9D A4 EF B8 8F E2 80 8D F0 9F 92 8B E2 80 8D F0 9F 91 A8*/},
					}},
				}},
			}},
		}},
	}}
	joined string
)

func init() {
	tr.Iterate(func(prefix []byte, value interface{}) {
		joined += string(prefix)
	})
}

func TestTrie_Scan(t *testing.T) {
	var str string
	var expected = []interface{}{}

	var runes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ|!?*+-/*{}[]()_^")

	tr.Iterate(func(prefix []byte, value interface{}) {
		expected = append(expected, string(prefix))
		str += string(prefix)

		// fill some random runes between keys
		for i := 0; i < rand.Intn(10); i++ {
			r := runes[rand.Int63()%int64(len(runes))]
			expected = append(expected, r)
			str += string(r)
		}
	})

	var found = make([]interface{}, 0)
	var ind = 0
	for ind < len(str) {
		_, size, ok := tr.SearchIn([]byte(str[ind:]))
		if ok {
			if ind+size > len(str) {
				t.Fatalf("index out of bounds: %d", ind+size)
			}
			found = append(found, str[ind:ind+size])
		} else {
			// Current sequence of bytes doesn't match any key.
			// Assuming there is a rune in the beginning
			var r rune
			r, size = utf8.DecodeRuneInString(str[ind:])
			if r == utf8.RuneError {
				t.Fatalf("it is neither emoji, nor rune: %s", str[ind:])
			}
			found = append(found, r)
		}

		ind += size
	}

	if !reflect.DeepEqual(found, expected) {
		t.Errorf("got:\n%v\nexpected:\n%v\n", found, expected)
	}
}

func TestTrie_GetString(t *testing.T) {
	var available []string
	tr.Iterate(func(prefix []byte, value interface{}) {
		available = append(available, string(prefix))
	})
	for i := range available {
		if _, ok := tr.GetString(available[i]); !ok {
			t.Errorf("Existing key \"%s\" was not found\n", available[i])
		}
	}
	for i := range available {
		key := available[i] + "some junk"
		if _, ok := tr.GetString(key); ok {

			t.Errorf("Not exsiting key \"%s\" was found\n", key)
		}
	}
}

func TestTrie_SplitKeysOnly(t1 *testing.T) {
	inputs := []string{"⏰", "✈️", "🆎", "🎟️", "◼"}

	tr := &Trie{}

	for _, em := range inputs {
		tr.Add([]byte(em), em)
	}

	str := strings.Join(inputs, "")
	results, err := tr.SplitKeysOnly(str)
	if err != nil {
		t1.Error(err)
	}

	if !reflect.DeepEqual(inputs, results) {
		t1.Errorf("Inputs != results: %s != %s", inputs, results)
	}
}

func BenchmarkTrie_Get(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, ok := tr.Get([]byte("👨‍❤️‍💋‍👨"))
		if !ok {
			b.Fail()
		}
	}
}

func BenchmarkTrie_Scan(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _, ok := tr.SearchIn([]byte("👨‍❤️‍💋‍👨"))
		if !ok {
			b.Fail()
		}
	}
}

func BenchmarkTrie_SplitKeysOnly(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := tr.SplitKeysOnly(joined)
		if err != nil {
			b.Error(err)
		}
	}
}
