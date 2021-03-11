package trie

// Iterate calls callback for each value stored in trie
//
// Not thread safe.
//
// Also prefix's underlying array would change on every call - so you can't rely on it after callback finishes
// (e.g. you should not pass it to another goroutine without copying)
//
// It seems like the only possible iteration order is by key (prefix):
//     0x1, 0x1 0x1, 0x1 0x2, 0x1 0x3, 0x2, 0x2 0x1, 0x2 0x2, etc...
// But it's not guarantied. You shouldn't rely on it!
func (t *Trie) Iterate(callback func(prefix []byte, value interface{})) {
	t.iterate(make([]byte, 0, 1024), callback)
}

func (t *Trie) iterate(prefix []byte, callback func([]byte, interface{})) {
	curPrefix := append(prefix[:len(prefix):len(prefix)], t.Prefix...)
	if t.Value != nil {
		callback(curPrefix, t.Value)
	}
	if t.Children != nil {
		for i := range t.Children {
			if t.Children[i] != nil {
				t.Children[i].iterate(curPrefix, callback)
			}
		}
	}
}
