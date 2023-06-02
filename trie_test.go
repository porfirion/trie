package trie

import (
	"bytes"
	"math/rand"
	"reflect"
	"testing"
	"unicode/utf8"
)

func TestTrie_Put(t *testing.T) {
	tests := []struct {
		prefix []byte
		val    string
	}{
		{[]byte(`👨‍`), `👨‍`},
		{[]byte(`👨‍🔧`), `👨‍🔧`},
		{[]byte(`👨`), `👨`},                           // 0xF0, 0x9F, 0x91, 0xA8
		{[]byte{0xF0, 0x9F, 0x91, 0x10}, "modified"}, // modified last byte of first key
		{[]byte{0xF0, 0x9F, 0x91, 0x10}, "modified"}, // add the same
		{[]byte{0xF0, 0x9F, 0x91}, "short"},          // add the same
	}

	tr := &Trie[string]{}
	for _, tt := range tests {
		tr.Put(tt.prefix, tt.val)
	}

	var expected = &Trie[string]{Prefix: []byte{0xF0, 0x9F, 0x91}, Value: ptr("short"), Children: &[256]*Trie[string]{
		0x10: {Prefix: []byte{0x10}, Value: ptr("modified")},
		0xA8: {Prefix: []byte{0xA8}, Value: ptr(`👨`), Children: &[256]*Trie[string]{
			0xE2: {Prefix: []byte{0xE2, 0x80, 0x8D}, Value: ptr(`👨‍`), Children: &[256]*Trie[string]{
				0xF0: {Prefix: []byte{0xF0, 0x9F, 0x94, 0xA7}, Value: ptr(`👨‍🔧`)},
			}},
		}},
	}}

	if !reflect.DeepEqual(tr, expected) {
		t.Fatalf("Not equal:\nexpected\n%s\ngot\n%s\n", expected, tr)
	}
}

func TestTrie_Put__Empty(t *testing.T) {
	tr := &Trie[string]{}
	// insert something before empty
	tr.PutString("foo", "bar")
	tr.Put(nil, "universal")

	if raw, ok := tr.Get(nil); !ok || raw != "universal" {
		t.Error("can't get value with zero prefix")
	}
	if raw, ok := tr.GetByString("foo"); !ok || raw != "bar" {
		t.Error("can't get foo")
	}

	tr = &Trie[string]{}
	// insert empty before others
	tr.Put(nil, "universal")
	tr.PutString("foo", "bar")
	if raw, ok := tr.Get(nil); !ok || raw != "universal" {
		t.Error("can't get value with zero prefix")
	}
	if raw, ok := tr.GetByString("foo"); !ok || raw != "bar" {
		t.Error("can't get foo")
	}

	// replace one empty with another
	tr = &Trie[string]{}
	tr.Put(nil, "universal")
	prev := tr.Put(nil, "universal2")

	if prev == "" {
		t.Error("there was previous value and it should be returned")
	}

	if raw, ok := tr.Get(nil); !ok || raw != "universal2" {
		t.Error("can't get value with zero prefix")
	}
}

// it's ok to add funcs if they implement Equaler
func TestTrie_AddFuncTwice(t *testing.T) {
	tr := &Trie[func()]{}
	tr.PutString("foo", func() {})
	tr.PutString("foo", func() {})
}

type T = Trie[string]

// inputs := []string{"⏰", "✈️", "🆎", "🎟️", "◼"}
var (
	tr = T{Prefix: []byte{0xF0, 0x9F, 0x91}, Children: &[256]*T{
		0xA8: {Prefix: []byte{0xA8}, Value: ptr("👨") /*F0 9F 91 A8*/, Children: &[256]*T{
			0xE2: {Prefix: []byte{0xE2, 0x80, 0x8D}, Children: &[256]*T{
				0xE2: {Prefix: []byte{0xE2}, Children: &[256]*T{
					0x9A: {Prefix: []byte{0x9A}, Children: &[256]*T{
						0x95: {Prefix: []byte{0x95, 0xEF, 0xB8, 0x8F}, Value: ptr("👨‍⚕️") /*F0 9F 91 A8 E2 80 8D E2 9A 95 EF B8 8F*/},
						0x96: {Prefix: []byte{0x96, 0xEF, 0xB8, 0x8F}, Value: ptr("👨‍⚖️") /*F0 9F 91 A8 E2 80 8D E2 9A 96 EF B8 8F*/},
					}},
					0x9C: {Prefix: []byte{0x9C, 0x88, 0xEF, 0xB8, 0x8F}, Value: ptr("👨‍✈️") /*F0 9F 91 A8 E2 80 8D E2 9C 88 EF B8 8F*/},
					0x9D: {Prefix: []byte{0x9D, 0xA4, 0xEF, 0xB8, 0x8F, 0xE2, 0x80, 0x8D, 0xF0, 0x9F}, Children: &[256]*T{
						0x91: {Prefix: []byte{0x91, 0xA8}, Value: ptr("👨‍❤️‍👨") /*F0 9F 91 A8 E2 80 8D E2 9D A4 EF B8 8F E2 80 8D F0 9F 91 A8*/},
						0x92: {Prefix: []byte{0x92, 0x8B, 0xE2, 0x80, 0x8D, 0xF0, 0x9F, 0x91, 0xA8}, Value: ptr("👨‍❤️‍💋‍👨") /*F0 9F 91 A8 E2 80 8D E2 9D A4 EF B8 8F E2 80 8D F0 9F 92 8B E2 80 8D F0 9F 91 A8*/},
					}},
				}},
			}},
		}},
	}}
)

func TestTrie_SearchPrefixIn(t *testing.T) {
	var str string
	var expected = []any{}

	var runes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ|!?*+-/*{}[]()_^")

	tr.Iterate(func(prefix []byte, value string) {
		expected = append(expected, string(prefix))
		str += string(prefix)

		// fill some random runes between keys
		for i := 0; i < rand.Intn(10); i++ {
			r := runes[rand.Int63()%int64(len(runes))]
			expected = append(expected, r)
			str += string(r)
		}
	})

	var found = make([]any, 0)
	var ind = 0
	for ind < len(str) {
		_, size, ok := tr.SearchPrefixIn([]byte(str[ind:]))
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

func TestTrie_GetAll(t *testing.T) {
	tr := BuildFromMap(map[string]string{
		"":                "root",
		"/api/user":       "user",
		"/api/user/list":  "list",
		"/api/group/":     "group",
		"/api/group/list": "list",
	})

	var inputs = map[string][]string{
		"/api/user/list": {"root", "user", "list"},
		"/api/test":      {"root"},
		"/api/user/li":   {"root", "user"},
	}

	for key, expected := range inputs {
		res := tr.GetAll([]byte(key))

		typedRes := make([]string, len(res))
		for i := range typedRes {
			typedRes[i] = res[i]
		}

		if !reflect.DeepEqual(typedRes, expected) {
			t.Errorf("%s: got %v, expected %v", key, res, expected)
		}
	}
}

func TestTrie_SubTrie(t *testing.T) {
	tr := BuildPrefixesOnly(
		"",
		"/api/user",
		"/api/user/list",
		"/api/group/",
		"/api/group/list",
		"/api/articles/list",
		"/api/articles/raw",
	)

	//fmt.Println(strings.Join(tr.toStrings(formatAsStrings), "\n"))

	type args struct {
		selector   string
		keepPrefix bool
	}
	type results struct {
		ok         bool
		rootPrefix []byte
	}

	selectors := map[args]results{
		{"/api/group", false}:        {ok: true, rootPrefix: []byte("/")},
		{"/api/group/", false}:       {ok: true, rootPrefix: []byte("")},
		{"/test/", false}:            {ok: false},
		{"/api/group", true}:         {ok: true, rootPrefix: []byte("/api/group/")},
		{"/api/articles/test", true}: {ok: false},
	}

	for args, res := range selectors {
		subTrie, ok := tr.SubTrie([]byte(args.selector), args.keepPrefix)
		if ok != res.ok {
			t.Errorf(`wrong result for %v: got %t expected %t`, args, ok, res.ok)
		} else if ok && !bytes.Equal(res.rootPrefix, subTrie.Prefix) {
			t.Errorf("wrong prefix in root Trie for %v: got %s expected %s", args, subTrie.Prefix, res.rootPrefix)
		}
	}
}

func TestTrie_GetByString(t *testing.T) {
	tr := BuildFromMap(map[string]string{
		"":                       "root",
		"/api/user":              "user",
		"/api/user/list":         "users list",
		"/api/group/":            "group",
		"/api/group/list":        "groups list",
		"/api/articles/list":     "articles list",
		"/api/articles/raw/list": "raw articles list",
	})

	type result struct {
		Value string
		OK    bool
	}

	var inputs = map[string]result{
		"":                {"root", true},
		"/api/user/list":  {"users list", true},
		"/api/user/1":     {"", false},
		"/api/articles/":  {"", false},
		"/api/article":    {"", false},
		"/api/articles/1": {"", false},
	}

	for key, res := range inputs {
		v, ok := tr.GetByString(key)
		if !reflect.DeepEqual(res.Value, v) || ok != res.OK {
			t.Errorf("get %v expected %v", result{v, ok}, res)
		}
	}
}

func TestTrie_Count(t *testing.T) {
	sources := map[string]string{
		"":                       "root",
		"/api/user":              "user",
		"/api/user/list":         "users list",
		"/api/group/":            "group",
		"/api/group/list":        "groups list",
		"/api/articles/list":     "articles list",
		"/api/articles/raw/list": "raw articles list",
	}

	tr := BuildFromMap(sources)

	// "/api/articles/" is common prefix for 2 entries,
	// but it doesn't store a value and is not included into result

	if l := tr.Count(); l != len(sources) {
		t.Errorf("got %d expected %d", l, len(sources))
	}

	if (*Trie[struct{}])(nil).Count() != 0 {
		t.Errorf("uninitialized tree count should return 0")
	}
}

// BenchmarkTrie_Put-8   	 1000000	      1268 ns/op	     509 B/op	       2 allocs/op
func BenchmarkTrie_Put(b *testing.B) {
	b.ReportAllocs()

	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	randomString := func() []byte {
		var b = make([]byte, 32)
		for i := range b {
			b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
		}
		return b
	}

	tr := &Trie[struct{}]{}
	for i := 0; i < b.N; i++ {
		// one allocation for random string generation
		tr.Put(randomString(), struct{}{})

		// This variant shows only 1 allocation,
		// but stopping and starting timer is very slow - benchmark can last for 30 seconds or more!
		//b.StopTimer()
		//str := randomString()
		//b.StartTimer()
		//tr.Put(str, struct{}{})
	}
}

// BenchmarkTrie_Get-8   	24781945	        43.63 ns/op	       0 B/op	       0 allocs/op
func BenchmarkTrie_Get(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, ok := tr.Get([]byte("👨‍❤️‍💋‍👨"))
		if !ok {
			b.Fail()
		}
	}
}

// BenchmarkTrie_SearchPrefixIn-8   	23258272	        49.96 ns/op	       0 B/op	       0 allocs/op
func BenchmarkTrie_SearchPrefixIn(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _, ok := tr.SearchPrefixIn([]byte("👨‍❤️‍💋‍👨"))
		if !ok {
			b.Fail()
		}
	}
}

func ptr[T any](v T) *T {
	return &v
}
