package trie

import (
	"fmt"
	"strings"
)

func (t Trie) String() string {
	return strings.Join(getStrings(&t), "\n")
}

func getStrings(t *Trie) []string {
	var resStrings = []string{
		fmt.Sprintf("[%s] %v", stringToBytes(string(t.Prefix)), t.Value),
	}

	if t.Children != nil {
		for ind, c := range t.Children {
			if c != nil {
				var childStrings = getStrings(c)
				resStrings = append(resStrings, fmt.Sprintf("├─%X─ %s", ind, childStrings[0]))
				resStrings = append(resStrings, addPrefix("│     ", childStrings[1:])...)
			}
		}
	}

	return resStrings
}

func addPrefix(prefix string, strs []string) []string {
	for ind, str := range strs {
		strs[ind] = prefix + str
	}
	return strs
}


func stringToBytes(val string) string {
	var bts = make([]string, 0, len(val))
	for i := 0; i < len(val); i++ {
		str := fmt.Sprintf("%X", val[i])
		bts = append(bts, str)
	}
	return strings.Join(bts, " ")
}
