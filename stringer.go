package trie

import (
	"fmt"
	"strings"
)

const (
	formatAsBytes   = 0
	formatAsStrings = 1
)

func (t Trie) String() string {
	return strings.Join(getStrings(&t, formatAsBytes), "\n")
}

func getStrings(t *Trie, format int) []string {
	var resStrings []string
	if format == formatAsStrings {
		resStrings = append(resStrings, fmt.Sprintf("[%s] %v", string(t.Prefix), t.Value))
	} else {
		resStrings = append(resStrings, fmt.Sprintf("[%s] %v", bytesToString(t.Prefix), t.Value))
	}

	if t.Children != nil {
		for ind, c := range t.Children {
			if c != nil {
				var childStrings = getStrings(c, format)
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

func bytesToString(val []byte) string {
	var bts = make([]string, 0, len(val))
	for i := 0; i < len(val); i++ {
		str := fmt.Sprintf("%X", val[i])
		bts = append(bts, str)
	}
	return strings.Join(bts, " ")
}
