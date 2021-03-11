package trie

// Because []byte can't be map key we only can have strings as index
func FromMap(inputs map[string]interface{}) *Trie{
	t := &Trie{}
	for key, value := range inputs {
		t.Add([]byte(key), value)
	}
	return t
}

func FromList(inputs []struct{Key []byte; Value interface{}}) *Trie {
	t := &Trie{}
	for i := range inputs {
		t.Add(inputs[i].Key, inputs[i].Value)
	}
	return t
}
