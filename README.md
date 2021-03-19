#Trie - compact and efficient radix tree (Patricia trie) implementation in go

Efficient implementation with zero allocation for read operations (Get and Search) and 1 or 2 allocations per Add

Current implementation (due to lack of generics) uses interface{} to store values. But it's type defined as an alias, and you can easily copy source file and replace alias with any other nil'able type (pointer or other interface).

    type ValueType = interface{}

Trie can be used in different ways:

1. Primarily I created it for searching emojis in text (in Telegram messages). There are about 3,3k emojis in current standard (https://www.unicode.org/emoji/charts-13.0/emoji-counts.html) and checking them one by one is very costly. For this purpose I added export package: you can generate static list of all available emojis and compile it in your program. 

2. You can use it as map, where key is a slice of bytes (`map[[]byte]interface{}` which is not possible in language because slices, used as keys, are not comparable).

3. You can use this trie to check for any string prefixes (possibly storing some payload for each prefix). See example below.

4. You can build some router using this Trie. For this purpose I added `SubTrie(mask []byte)` method that returns sub trie with all entries, prefixed by specified mask, and `GetAll(mask []byte)` that returns all entries containing specified mask. See example below.

Also this implementation supports zero-length prefix (`[]byte{}` or `nil`). Value associated with this prefix can be used as fallback when no other entries found. Or it can serve as universal prefix for all entries.