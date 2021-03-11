package prefix

import (
	"fmt"
)

func ExampleFromStrings() {
	prefixes := FromStrings([]string{
		"tiny_",
		"normal_",
		"large_",
	})

	var myString = "large_banana"

	if prefix, ok := prefixes.Check(myString); ok {
		fmt.Printf("prefix \"%s\" found\n", prefix)
	} else {
		fmt.Println("no prefix found")
	}
	// Output:
	//prefix "large_" found
}
