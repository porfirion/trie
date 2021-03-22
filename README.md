#Trie - compact and efficient radix tree (Patricia trie) implementation in go

Efficient implementation with zero allocation for read operations (Get and Search) and 1 or 2 allocations per Add

Current implementation (due to lack of generics) uses interface{} to store values. But it's type defined as an alias, and you can easily copy source file and replace alias with any other nil'able type (pointer or other interface) and get a definitely typed implementation:

```go
type ValueType = interface{}
```

##Installation

    go get github.com/porfirion/trie

##Usage

Trie can be used in different ways:

1. Primarily I created it for searching emojis :smile: in text (in Telegram messages). There are about 3,3k emojis in current standard (https://www.unicode.org/emoji/charts-13.0/emoji-counts.html) and checking them one by one is very costly. For this purpose I added export package: you can generate source code for trie with all available emojis and compile it in your program. 

2. You can use it as map, where key is a slice of arbitrary bytes (`map[[]byte]interface{}` which is not possible in language because slices are not comparable and can't be used as keys).

3. You can use this trie to check for any string prefixes (possibly storing some payload for each prefix). See example below.

4. You can build some router using this Trie. For this purpose I added `SubTrie(mask []byte)` method that returns sub trie with all entries, prefixed by specified mask, and `GetAll(mask []byte)` that returns all entries containing specified mask. See example below.

Also, this implementation supports zero-length prefix (`[]byte{}` or `nil`). Value associated with this prefix can be used as fallback when no other entries found. Or it can serve as universal prefix for all entries.

Limitation: `nil` can't be stored as value, because node containing nil value considered empty.

##Examples

Search prefixes:
```go
package main

import (
    "fmt"
    "strings"
    "github.com/porfirion/trie"
)

// Also can be created with
//    prefixes := &trie.Trie{}
//    prefixes.PutString("one", 1)
//    prefixes.PutString("two", 2)
//    prefixes.PutString("three", 3)
//    prefixes.PutString("", 0)
//
var prefixes = trie.BuildFromMap(map[string]interface{}{
    "one":   1,
    "two":   2,
    "three": 3,
    "":      0,
})

func Example() {
    var inputs = []string{
        "twoApple",
        "oneBanana",
        "Carrot",
    }

    for _, inp := range inputs {
        if val, prefixLen, ok := prefixes.SearchPrefixInString(inp); ok {
            fmt.Println(strings.Repeat(inp[prefixLen:], val.(int)))
        }
    }

    // Output:
    //AppleApple
    //Banana
}
```

Use as router:
```go
package main

import (
    "fmt"
    "github.com/porfirion/trie"
)

var routes = trie.BuildFromMap(map[string]interface{}{
    "":                "root", // as universal prefix
    "/api/user":       "user",
    "/api/user/list":  "usersList",
    "/api/group":      "group",
    "/api/group/list": "groupsList",
    "/api/articles/":  "articles",
})

func Example_routing() {
    var inputs = []string{
        "/api/user/list",
        "/api/user/li",
        "/api/group",
        "/api/unknown",
    }

    for _, inp := range inputs {
        exact, ok := routes.GetByString(inp)
        route := routes.GetAllByString(inp)
        if ok {
            fmt.Printf("%-17s:\thandler %-10s\t(route %v)\n", inp, exact, route)
        } else {
            fmt.Printf("%-17s:\thandler not found\t(route %v)\n", inp, route)
        }
    }

    // Output:
    // /api/user/list   :	handler usersList 	(route [root user usersList])
    // /api/user/li     :	handler not found	(route [root user])
    // /api/group       :	handler group     	(route [root group])
    // /api/unknown     :	handler not found	(route [root])
}
```

##Notes

I didn't implement `Delete`/`Remove` operation as I assumed only checking for predefined list of prefixes. I you find it useful and really need it - please, open issue, I'll add it.