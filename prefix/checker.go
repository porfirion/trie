package prefix

import "github.com/porfirion/trie"

type Checker trie.Trie

func FromStrings(strs []string) Checker {
	t := &trie.Trie{}

	for i := range strs {
		t.Add([]byte(strs[i]), struct{}{})
	}

	return Checker(*t)
}

func (t Checker) Check(str string) (prefix string, ok bool) {
	bts := []byte(str)
	tr := trie.Trie(t)

	_, length, ok := tr.Scan(bts)
	if ok {
		return string(bts[:length]), true
	}

	return "", false
}
