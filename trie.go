package trie

// Trie implements sparse radix (Patricia) trie.
// Makes zero allocation on Get and SearchPrefixIn operations and two allocations per Put
//
// Create it just as &Trie{} and add required data.
// Also there are some convenience constructors (for example for initialization from map[prefix]value)
type Trie[T any] struct {
	Prefix   []byte
	Value    *T
	Children *[256]*Trie[T]
}

// PutString is a convenience method for Put()
func (t *Trie[T]) PutString(prefix string, value T) {
	t.Put([]byte(prefix), value)
}

// Put adds new entry into trie or replaces existing with specified prefix.
// Prefix can have zero length - associated value would be added into root Trie
//
// There were some variants of control of replace. Handling it in place results in dependency on ValueType
// (interfaces and pointers are handled differently). Also there are some incomparable types, that cause panic when compared :(.
// Store separate function like OnReplaceCallback - requires passing it to children. And if you update it in parent -
// it wouldn't be updated in children automatically. Another problem - node doesn't know it's whole prefix, since
// it doesn't know it's parent!
// With current realization caller of Put knows whole prefix and we shouldn't collect it inside. IMHO, much easier.
func (t *Trie[T]) Put(newPrefix []byte, val T) (oldValue T) {
	var curPrefix = t.Prefix
	var ind int
	for ind < len(curPrefix) && ind < len(newPrefix) && curPrefix[ind] == newPrefix[ind] {
		// find common part of current and new prefixes
		ind++
	}

	if ind == len(curPrefix) {
		// reached curPrefix end
		if ind == len(newPrefix) {
			// complete match  of prefixes
			// put or replace current value
			// also for case of zero length prefix
			if t.Value != nil {
				oldValue = *t.Value
			}
			t.Value = &val
		} else {
			// ind < len(newPrefix) - newPrefix longer than existing
			if len(t.Prefix) == 0 && t.Children == nil && t.Value == nil {
				// case for empty Trie and first insertion
				// insert prefix and value into Trie itself
				t.Prefix = newPrefix
				t.Value = &val
			} else {
				// our trie is not empty (we already have value or children)
				// rest of newPrefix would be added into proper child
				oldValue = t.getChildOrCreate(newPrefix[ind]).Put(newPrefix[ind:], val)
			}
		}
	} else {
		// ind < len(curPrefix)

		// cur ****|**
		// new ****|
		// OR
		// cur ****|**
		// new ****|*

		// newPrefix shorted than existing or they diverged.
		// Split current prefix into parts (common part [0:ind] and the rest [ind:])
		// and place all current fields into newChild
		var newChild = &Trie[T]{
			Prefix:   curPrefix[ind:], // take only diverging part of prefix
			Value:    t.Value,
			Children: t.Children,
		}

		// reset current Trie and add newChild
		t.Prefix = curPrefix[:ind]    // common part (in worst case - it would be empty slice)
		t.Value = nil                 // no value - it's prefix only
		t.Children = &[256]*Trie[T]{} // it would have a child anyway
		t.Children[newChild.Prefix[0]] = newChild

		// what to do with new value?
		if ind == len(newPrefix) {
			// newPrefix equals common part! Current Trie becomes value
			t.Value = &val
		} else {
			// newPrefix longer than common part. Rest of newPrefix would be set into proper child
			oldValue = t.getChildOrCreate(newPrefix[ind]).Put(newPrefix[ind:], val)
		}
	}

	return
}

func (t *Trie[T]) getChildOrCreate(ind byte) *Trie[T] {
	if t.Children == nil {
		t.Children = &[256]*Trie[T]{}
		t.Children[ind] = &Trie[T]{}
	} else if t.Children[ind] == nil {
		t.Children[ind] = &Trie[T]{}
	}
	return t.Children[ind]
}

// GetByString is a convenience method for Get
func (t *Trie[T]) GetByString(key string) (T, bool) {
	return t.Get([]byte(key))
}

// Get searches for exactly matching key in trie
func (t *Trie[T]) Get(key []byte) (res T, found bool) {
	ind := 0
	for ind < len(t.Prefix) && ind < len(key) && t.Prefix[ind] == key[ind] {
		ind++
	}

	if ind < len(t.Prefix) {
		// prefix didn't match
		// it is not this trie or it's child
		return res, false
	}

	if ind < len(key) {
		// not all key bytes matched
		if t.Children != nil && t.Children[key[ind]] != nil {
			// continue matching children with next bytes of key.
			return t.Children[key[ind]].Get(key[ind:])
		}

		// we have no child with such prefix
		return res, false
	}

	if t.Value == nil {
		return res, false
	}

	return *t.Value, true
}

// TakePrefix returns only found prefix without corresponding value.
// Just for convenience.
func (t *Trie[T]) TakePrefix(str string) (prefix string, ok bool) {
	_, length, ok := t.SearchPrefixInString(str)
	if ok {
		return str[:length], true
	}

	return "", false
}

// SearchPrefixInString is a convenience method for SearchPrefixIn
func (t *Trie[T]) SearchPrefixInString(str string) (value T, prefixLen int, ok bool) {
	return t.SearchPrefixIn([]byte(str))
}

// SearchPrefixIn searches the longest matching prefix in input bytes.
// If input has prefix that matches any stored key
// returns associated value, prefix length, true OR nil, 0, false otherwise
func (t *Trie[T]) SearchPrefixIn(input []byte) (value T, prefixLen int, ok bool) {
	ind := 0
	for ind < len(t.Prefix) && ind < len(input) && t.Prefix[ind] == input[ind] {
		ind++
	}

	if ind < len(t.Prefix) {
		// prefix didn't match It is not this trie or it's child
		return value, 0, false
	}

	if ind < len(input) && t.Children != nil && t.Children[input[ind]] != nil {
		// continue matching children with next bytes from input. Greedy!
		value, prefixLen, ok = t.Children[input[ind]].SearchPrefixIn(input[ind:])
	}

	if ok {
		// we found something in children!
		prefixLen += len(t.Prefix) // our prefix should be added to children's

		return value, prefixLen, ok
	}

	if t.Value != nil {
		// take our value
		return *t.Value, len(t.Prefix), true
	}

	// we have no value. Explicitly return size 0 because we can have prefix, but it doesn't matter
	return value, 0, false
}

// Iterate calls callback for each value stored in trie
//
// Not thread safe.
//
// Also prefix's underlying array would change on every call - so you can't rely on it after callback finishes
// (e.g. you should not pass it to another goroutine without copying)
//
// It seems like the only possible iteration order is by key (prefix):
//
//	0x1, 0x1 0x1, 0x1 0x2, 0x1 0x3, 0x2, 0x2 0x1, 0x2 0x2, etc...
//
// But it's not guarantied. You shouldn't rely on it!
//
// Can be used in combination with SubTrie:
//
//	tr.SubTrie(mask).Iterate(func...)
func (t *Trie[T]) Iterate(callback func(prefix []byte, value T)) {
	t.iterate(make([]byte, 0, 1024), callback)
}

func (t *Trie[T]) iterate(prefix []byte, callback func([]byte, T)) {
	curPrefix := append(prefix[:len(prefix):len(prefix)], t.Prefix...)
	if t.Value != nil {
		callback(curPrefix, *t.Value)
	}
	if t.Children != nil {
		for i := range t.Children {
			if t.Children[i] != nil {
				t.Children[i].iterate(curPrefix, callback)
			}
		}
	}
}

// SubTrie returns new Trie with all values whose prefixes include mask
//
// keepPrefix indicates whether the new trie should keep original prefixes,
// or should contain only those parts, that are out of mask:
//
//	tr := {"": v0, "/user/": v1, "/user/list": v2, "/group/": v3}.
//
//	tr.SubTrie("/user", false)
//	-> {"/": v1, "/list": v2}
//
//	tr.SubTrie("/user", true)
//	-> {"/user/": v1, "/user/list": v2}
func (t *Trie[T]) SubTrie(mask []byte, keepPrefix bool) (subTrie *Trie[T], ok bool) {
	return t.subTrie(mask, keepPrefix, mask, 0)
}

func (t *Trie[T]) subTrie(mask []byte, keepPrefix bool, originalMask []byte, originalMaskInd int) (subTrie *Trie[T], ok bool) {
	var ind = 0
	for ind < len(mask) && ind < len(t.Prefix) && mask[ind] == t.Prefix[ind] {
		ind++
	}

	if ind == len(mask) {
		// complete match for mask!
		res := &Trie[T]{
			// Prefix to be filled
			Value:    t.Value,
			Children: t.Children,
		}

		if keepPrefix {
			res.Prefix = make([]byte, len(originalMask[:originalMaskInd])+len(t.Prefix))
			copy(res.Prefix, originalMask[:originalMaskInd])
			copy(res.Prefix[ind:], t.Prefix)
		} else {
			// copy the rest of t.Prefix (would be empty if t.Prefix match also complete)
			res.Prefix = t.Prefix[ind:]
		}

		return res, true
	} else {
		// ind < len(mask)
		if ind == len(t.Prefix) {
			// match with current t.Prefix is complete. Continue match with it's child
			if t.Children != nil && t.Children[mask[ind]] != nil {
				return t.Children[mask[ind]].subTrie(mask[ind:], keepPrefix, originalMask, originalMaskInd+ind)
			} else {
				// no such child to continue(
				return nil, false
			}
		} else {
			// mask and t.Prefix diverged
			return nil, false
		}
	}
}

// GetAllByString is a convenience method for GetAll
func (t *Trie[T]) GetAllByString(str string) []T {
	return t.GetAll([]byte(str))
}

// GetAll returns all Values whose prefixes are subsets of mask
//
//	tr := {"": v0, "/user/": v1, "/user/list": v2, "/group/": v3}
//
//	tr.GetAll("/user/list", false)
//	-> [v0, v1, v2]
func (t *Trie[T]) GetAll(mask []byte) []T {
	var ind = 0
	for ind < len(mask) && ind < len(t.Prefix) && mask[ind] == t.Prefix[ind] {
		ind++
	}
	if ind == len(t.Prefix) {
		var res = make([]T, 0, 1)
		if t.Value != nil {
			res = append(res, *t.Value)
		}

		if ind < len(mask) && t.Children != nil && t.Children[mask[ind]] != nil {
			return append(res, t.Children[mask[ind]].GetAll(mask[ind:])...)
		}

		return res
	} else {
		// doesn't match current prefix
		return nil
	}
}

// Count returns amount of values (non nil) stored in all nodes of trie.
func (t *Trie[T]) Count() int {
	if t == nil {
		return 0
	}
	var count = 0
	if t.Value != nil {
		count++
	}
	if t.Children != nil {
		for i := range t.Children {
			if t.Children[i] != nil {
				count += t.Children[i].Count()
			}
		}
	}
	return count
}
