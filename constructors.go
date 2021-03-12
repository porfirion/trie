package trie

func BuildFromMap(inputs map[string]ValueType) *Trie {
	t := &Trie{}
	for key, value := range inputs {
		t.Put([]byte(key), value)
	}
	return t
}

// []byte can't be map key. So we can use list of structs
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

// Useful for just searching prefixes without any data
func BuildPrefixesOnly(strs ...string) *Trie {
	t := &Trie{}

	for i := range strs {
		t.Put([]byte(strs[i]), struct{}{})
	}

	return t
}
