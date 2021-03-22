// Package trie contains radix tree implementation in golang (without any dependencies).
//
// Copyright 2021, Mikhail Vitsen (@porfirion)
// https://github.com/porfirion/trie
//
// Trie maps `[]byte` prefixes to `interface{}` values.
// It can be used for efficient search of substrings or bulk prefix checking.
//
// Trie is created by
//   t := &Trie{}
//   t.Put("foo", "bar")
//   t.Put("buzz", 123)
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
//   value, prefixLen, ok := t.SearchPrefixIn(input) // searches the longest prefix of input, that stored in trie
// but also handy method
//   prefix, ok := t.TakePrefix(input) // returns longest prefix without associated value
package trie
