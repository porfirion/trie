package trie

import (
	"fmt"
)

type Equaler interface {
	Equal(v interface{}) bool
}

// Sparse radix trie. Create it just as &Trie{} and add required data.
// Also there are some convenience constructors (for example for one line initialization)
// Makes zero allocation on Get and Scan operations
type Trie struct {
	Prefix   []byte
	Value    interface{}
	Children *[256]*Trie
}

// Add adds new entry into trie with specified prefix
//
// WARNING! nil shouldn't be stored as value: you wouldn't be able to find it nor by Scan, nor by Get, nor by Iterate
// If you don't need any value (you need only prefixes) - you can use struct{}{}. See checker.FromStrings
//
// WARNING! Add will panic when called with already existing prefix and different value
// For such situations (when the same value is added twice and it's not comparable by default operator (==))
// you can implement Equaler, that checks if it is the same value and no panic needed.
func (t *Trie) Add(newPrefix []byte, val interface{}) {
	//fmt.Printf("Adding %X (%v)\n", newPrefix, val)
	if len(newPrefix) == 0 {
		panic("zero length newPrefix")
	}
	if len(t.Prefix) == 0 && t.Children == nil {
		// just empty node. Set it's value
		if t.Value != nil {
			panic("there is already another value")
		}
		t.Prefix = newPrefix
		t.Value = val
	}

	var curPrefix = t.Prefix
	var ind int
	for ind < len(curPrefix) && ind < len(newPrefix) && curPrefix[ind] == newPrefix[ind] {
		// stop ind the end of arrays or when values diverged
		ind++
	}

	if len(curPrefix) == ind {
		// reached curPrefix end
		if len(newPrefix) == ind {
			// reached newPrefix end
			// equal bytes slices!
			if t.Value == nil {
				// just setting new value
				t.Value = val
			} else if t.Value != val {
				if eq, ok := t.Value.(Equaler); ok && eq.Equal(val) {
					// it's ok. It's the same
				} else {
					panic(fmt.Errorf("there is already another value:\n\t%s\n\t%s", t.Value, val))
				}
			}
		} else {
			// len(newPrefix) > ind - newPrefix longer than existing
			// rest of newPrefix would be added into proper child
			t.getChildOrCreate(newPrefix[ind]).Add(newPrefix[ind:], val)
		}
	} else {
		// len(curPrefix) > ind
		//new newPrefix shorted than existing or they diverged.
		//Split current prefix into parts (common part [0:ind] and the rest [ind:])

		// copy existing Trie into newChild, but take rest of prefix
		var newChild = &Trie{
			Prefix:   curPrefix[ind:], // take only diverging part
			Value:    t.Value,
			Children: t.Children,
		}

		// reset current Trie and add previous as newChild
		t.Prefix = curPrefix[:ind] // common part (in worst case - it would be empty slice)
		t.Value = nil              // no value - it's prefix only
		t.Children = &[256]*Trie{} // it would have a child anyway
		t.Children[newChild.Prefix[0]] = newChild

		// what to do with new value?
		if len(newPrefix) == ind {
			// newPrefix equals common part! Current Trie becomes value
			t.Value = val
		} else {
			// newPrefix longer than common part. Rest of newPrefix would be set into proper child
			t.getChildOrCreate(newPrefix[ind]).Add(newPrefix[ind:], val)
		}
	}
}

func (t *Trie) getChildOrCreate(ind byte) *Trie {
	if t.Children == nil {
		t.Children = &[256]*Trie{}
		t.Children[ind] = &Trie{}
	} else if t.Children[ind] == nil {
		t.Children[ind] = &Trie{}
	}
	return t.Children[ind]
}
