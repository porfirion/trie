#Trie - compact and efficient radix tree implementation in go

Efficient implementation with zero allocation for read operations (Get and Search) and 1 or 2 allocation per Add

Current implementation (due to lack of generics) uses interface{} to store values. But it stored in alias and you can easily copy source file and replace it with any other nillable type (pointer or other interface).

    type ValueType = interface{}


