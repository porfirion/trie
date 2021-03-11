package export

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/porfirion/trie"
)

// Exports Trie as go code (compatible with gofmt).
func Export(t *trie.Trie, settings ExportSettings, currentPadding string) string {
	if t == nil {
		return currentPadding + "nil"
	}

	builder := &strings.Builder{}

	trieName := settings.TrieAlias
	if trieName == "" {
		trieName = "Trie"

		if len(settings.PackagePrefix) > 0 {
			trieName = settings.PackagePrefix + "." + trieName
		}
	}

	builder.WriteString("{") // start object >>>>>>>>>>>>

	var fields = []string{}

	if len(t.Prefix) > 0 {
		fields = append(fields, fmt.Sprintf("Prefix: []byte{%s}", encodeBytes(t.Prefix))) // prefix
	}
	if t.Value != nil {
		fields = append(fields, fmt.Sprintf("Value: %s", exportValue(t.Value))) // value
	}
	if t.Children != nil {
		var children []string
		for i, c := range t.Children {
			if c != nil {
				children = append(children, fmt.Sprintf(
					"%s%s0x%X: %s,\n",
					currentPadding,
					settings.Padding,
					i,
					Export(c, settings, currentPadding+settings.Padding),
				))
			}
		}

		fields = append(fields, fmt.Sprintf(
			"Children: &[256]*%s{\n%s%s}",
			trieName,
			strings.Join(children, ""),
			currentPadding,
		)) // begin children
	}
	builder.WriteString(strings.Join(fields, ", "))
	builder.WriteString("}") // <<<<<<<<<<<< end object

	return builder.String()
}

type ExportSettings struct {
	PackagePrefix string
	Padding       string
	TrieAlias     string
}

type Exportable interface {
	Export() string
}

func exportValue(v interface{}) string {
	if v == nil {
		return `nil`
	}
	switch val := v.(type) {
	case Exportable:
		return val.Export()
	case string:
		return fmt.Sprintf(`"%s"`, val)
		//return fmt.Sprintf(`"%s" /*%s*/`, val, stringToBytes(val))
	case int:
		return strconv.FormatInt(int64(val), 10)
	case int32:
		return strconv.FormatInt(int64(val), 10)
	case int64:
		return strconv.FormatInt(val, 10)
	case uint:
		return strconv.FormatUint(uint64(val), 10)
	case uint32:
		return strconv.FormatUint(uint64(val), 10)
	case uint64:
		return strconv.FormatUint(uint64(val), 10)
	case float32:
		return strconv.FormatFloat(float64(val), 'f', -1, 64)
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(val)
	default:
		return fmt.Sprintf(`"%+v"`, v)
	}
}

func encodeBytes(bts []byte) string {
	b := &strings.Builder{}
	for i := range bts {
		b.WriteString(fmt.Sprintf("0x%X", bts[i]))
		if i < len(bts)-1 {
			b.WriteString(", ")
		}
	}
	return b.String()
}
