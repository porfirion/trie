package trie

// Because []byte can't be map key we only can have strings as index
func BuildFromMap(inputs map[string]interface{}) *Trie {
	t := &Trie{}
	for key, value := range inputs {
		t.Add([]byte(key), value)
	}
	return t
}

func BuildFromList(inputs []struct {
	Key   []byte
	Value interface{}
}) *Trie {
	t := &Trie{}
	for i := range inputs {
		t.Add(inputs[i].Key, inputs[i].Value)
	}
	return t
}

// Useful for just searching prefixes without any data
func BuildPrefixesOnly(strs ...string) *Trie {
	t := &Trie{}

	for i := range strs {
		t.Add([]byte(strs[i]), struct{}{})
	}

	return t
}
