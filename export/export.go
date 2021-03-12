package export

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/porfirion/trie"
)

var (
	// string representation of byte (0x01..0xFF)
	bytesRep [256]string
)

func init() {
	for i := 0; i < 256; i++ {
		bytesRep[byte(i)] = fmt.Sprintf("0x%X", i)
	}
}

// Exports Trie as go code (compatible with gofmt).
func Export(t *trie.Trie, settings ExportSettings) string {
	if t == nil {
		return settings.CurrentPadding + "nil"
	}

	builder := &strings.Builder{}
	builder.Grow(200)

	builder.WriteString("{") // start object >>>>>>>>>>>>

	var fields = make([]string, 0, 4)

	if len(t.Prefix) > 0 {
		fields = append(fields, "Prefix: []byte{"+encodeBytes(t.Prefix)+"}") // prefix
	}
	if t.Value != nil {
		fields = append(fields, "Value: "+exportValue(t.Value)) // value
	}
	if t.Children != nil {
		var children = make([]string, 0, 256)

		for i, c := range t.Children {
			if c != nil {
				children = append(children,
					settings.CurrentPadding+
						settings.Padding+
						bytesRep[byte(i)]+": "+
						Export(c, settings.ForChild())+
						",\n",
				)
			}
		}

		fields = append(fields,
			"Children: &[256]*"+settings.GetTrieName()+"{\n"+
				strings.Join(children, "")+
				settings.CurrentPadding+"}",
		)
	}
	builder.WriteString(strings.Join(fields, ", "))
	builder.WriteString("}") // <<<<<<<<<<<< end object

	return builder.String()
}

type ExportSettings struct {
	TrieAlias      string
	PackagePrefix  string
	Padding        string
	CurrentPadding string
}

func (settings ExportSettings) GetTrieName() string {
	if settings.TrieAlias != "" {
		return settings.TrieAlias
	} else if len(settings.PackagePrefix) > 0 {
		return settings.PackagePrefix + ".Trie"
	} else {
		return "Trie"
	}
}

func (settings ExportSettings) ForChild() ExportSettings {
	settings.CurrentPadding += settings.Padding
	return settings
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
	case *string:
		return `"` + *val + `"` // + fmt.Sprintf(`/*%s*/`, stringToBytes(val))
	case string:
		return `"` + val + `"` // + fmt.Sprintf(`/*%s*/`, stringToBytes(val))
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
	if len(bts) == 0 {
		return ""
	}
	b := &strings.Builder{}
	b.Grow(len(bts) * 6)

	for i := range bts {
		if i < len(bts)-1 {
			b.WriteString(bytesRep[bts[i]] + ", ")
		} else {
			b.WriteString(bytesRep[bts[i]])
		}
	}
	return b.String()

}
