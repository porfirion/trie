package trie

// BuildFromMap may be useful for var declaration
func BuildFromMap(inputs map[string]ValueType) *Trie {
	t := &Trie{}
	for key, value := range inputs {
		t.Put([]byte(key), value)
	}
	return t
}

// BuildFromList can be used to create Trie with arbitrary bytes slice as key (not valid strings, etc)
func BuildFromList(inputs []struct {
	Key   []byte
	Value ValueType
}) *Trie {
	t := &Trie{}
	for i := range inputs {
		t.Put(inputs[i].Key, inputs[i].Value)
	}
	return t
}

// BuildPrefixesOnly used to create just searching prefixes without any data
func BuildPrefixesOnly(strs ...string) *Trie {
	type dummy struct{}

	t := &Trie{}

	for i := range strs {
		t.Put([]byte(strs[i]), dummy{})
	}

	return t
}
