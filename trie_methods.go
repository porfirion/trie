package trie

import "fmt"

func (t *Trie) ScanString(str string) (value interface{}, prefixLen int, ok bool) {
	return t.Scan([]byte(str))
}

// Scan searches the longest matching key (prefix) in input bytes.
// If input contains any key as prefix - return associated value, prefix length, true
func (t *Trie) Scan(input []byte) (value interface{}, prefixLen int, ok bool) {
	ind := 0
	for ind < len(t.Prefix) && ind < len(input) && t.Prefix[ind] == input[ind] {
		ind++
	}

	if ind < len(t.Prefix) {
		// prefix didn't match It is not this trie or it's child
		return nil, 0, false
	}

	if ind < len(input) && t.Children != nil && t.Children[input[ind]] != nil {
		// continue matching children with next bytes from input. Greedy!
		value, prefixLen, ok = t.Children[input[ind]].Scan(input[ind:])
	}

	if ok {
		// we found something in children!
		prefixLen += len(t.Prefix) // our prefix should be added to children's

		return value, prefixLen, ok
	}

	if t.Value != nil {
		// take our value
		return t.Value, len(t.Prefix), true
	}

	// we have no value. Explicitly return size 0 because we can have prefix, but it doesn't matter
	return nil, 0, false
}

func (t *Trie) GetString(key string) (interface{}, bool) {
	return t.Get([]byte(key))
}

// Get searches for exactly matching key in trie
func (t *Trie) Get(key []byte) (interface{}, bool) {
	ind := 0
	for ind < len(t.Prefix) && ind < len(key) && t.Prefix[ind] == key[ind] {
		ind++
	}

	if ind < len(t.Prefix) {
		// prefix didn't match
		// it is not this trie or it's child
		return nil, false
	}

	if ind < len(key) {
		// not all key bytes matched
		if t.Children != nil && t.Children[key[ind]] != nil {
			// continue matching children with next bytes of key.
			return t.Children[key[ind]].Get(key[ind:])
		}

		// we have no child with such prefix
		return nil, false
	}

	if t.Value == nil {
		// all key matched, but current trie has no value (assuming we have some children with values)
		return nil, false
	}

	return t.Value, true
}

// SplitKeysOnly splits string into slice of keys
func (t *Trie) SplitKeysOnly(str string) (res []string, err error) {
	var (
		ind = 0
		bts = []byte(str)
	)
	for ind < len(str) {
		_, size, ok := t.Scan(bts[ind:])
		if !ok {
			return res, fmt.Errorf("not a key: %s", str[ind:])
		}
		res = append(res, str[ind:ind+size])
		ind += size
	}

	return res, nil
}
