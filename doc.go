// Sparse radix tree implementation in golang (without any dependencies).
//
// Trie maps `[]byte` prefixes to `interface{}` values.
// It can be used for efficient search of substrings or bulk prefix checking.
//
// Trie is created by
//   t := &Trie{}
//   t.Add("foo", "bar")
//   t.Add("buzz", 123)
// or by convenience constructors:
//   t := BuildFromMap(map[string]interface{}{
//     "foo": "bar",
//     "buzz": 123,
//   })
// or
//  t := BuildFromList([]struct{Key string; Value interface{}}{
//    {Key: "foo", Value: "bar"},
//    {Key: "buzz", Value: 123},
//  })
// Two common operations are:
//   value, ok := t.Get(key) // returns associated value
//   value, prefixLen, ok := t.SearchIn(input) // searches the longest prefix of input, that stored in trie
// but also handy method
//   prefix, ok := t.SearchPrefix(input) // returns longest prefix without associated value
package trie
