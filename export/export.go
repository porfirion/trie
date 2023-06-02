package export

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/porfirion/trie"
)

//go:generate go run gen.go

type Exportable interface {
	Export() string
}

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
func Export[T any](t *trie.Trie[T], settings ExportSettings) string {
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

func exportGenericType[T any]() string {
	var v T
	res := reflect.TypeOf(&v).Elem().String()
	if res == "interface {}" {
		return "[any]"
	} else {
		return "[" + res + "]"
	}
}

func exportValue[T any](v *T) (res string) {
	if v == nil {
		return `nil`
	}
	defer func() {
		res = "ptr" + exportGenericType[T]() + "(" + res + ")"
	}()
	defer func() {
		if reflect.ValueOf(v).Elem().Type().Kind() != reflect.Interface {
			//    │          │      └ T type itself
			//    │          └ value T
			//    └ pointer to T

			//
			return
		}
		tp := reflect.ValueOf(v).Elem().Elem().Type()
		switch tp.Kind() {
		case reflect.String, reflect.Bool, reflect.Int:
			return
		default:
			res = "(" + tp.String() + ")" + "(" + res + ")"
			return
		}
	}()
	switch val := any(*v).(type) {
	case Exportable:
		return val.Export()
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
		// There can be wrong formatting.
		// If so - you should just implement Exportable interface
		return fmt.Sprintf(`%+v`, *v)
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
