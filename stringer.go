package trie

import (
	"fmt"
	"strings"
)

type prefixFormat int

const (
	formatAsBytes   prefixFormat = 0
	formatAsStrings prefixFormat = 1
)

func (t Trie[T]) String() string {
	return strings.Join(t.toStrings(formatAsBytes), "\n")
}

func (t Trie[T]) toStrings(format prefixFormat) []string {
	var resStrings []string
	if format == formatAsStrings {
		resStrings = append(resStrings, fmt.Sprintf("[%s] %s", string(t.Prefix), t.valueToString(t.Value)))
	} else {
		resStrings = append(resStrings, fmt.Sprintf("[%s] %s", bytesToString(t.Prefix), t.valueToString(t.Value)))
	}

	if t.Children != nil {
		for ind, c := range t.Children {
			if c != nil {
				var childStrings = c.toStrings(format)
				resStrings = append(resStrings, fmt.Sprintf("├─%X─ %s", ind, childStrings[0]))
				resStrings = append(resStrings, addPrefix(childStrings[1:], "│     ")...)
			}
		}
	}

	return resStrings
}

func (t Trie[T]) valueToString(v *T) string {
	if v == nil {
		return "nil"
	} else {
		return fmt.Sprintf("%+v", *v)
	}
}

func addPrefix(strs []string, prefix string) []string {
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
