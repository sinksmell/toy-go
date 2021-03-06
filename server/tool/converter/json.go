package converter

import (
	"fmt"
	"go/format"
	"math"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/sinksmell/toy-go/util"

	"github.com/buger/jsonparser"
)

// JsonConverter convert json to go struct 
type JsonConverter struct {
	stk       util.Stack
	sb        strings.Builder
	root      *Node
	omitEmpty bool
}

// JsonConvertOption build option of json converter 
type JsonConvertOption interface {
	apply(converter *JsonConverter)
}

type funcJsonConvertOpt func(converter *JsonConverter)

func (f funcJsonConvertOpt) apply(c *JsonConverter) {
	f(c)
}

// NewJsonConverter new json converter instance
func NewJsonConverter(opt ...JsonConvertOption) *JsonConverter {
	c := newDefaultJsonConverter()
	for _, p := range opt {
		p.apply(c)
	}
	return c
}

// WithOmitEmptyTag if omitEmpty is true, generate result will add omitempty tag for json
func WithOmitEmptyTag(omitEmpty bool) JsonConvertOption {
	return funcJsonConvertOpt(func(c *JsonConverter) {
		c.omitEmpty = omitEmpty
	})
}

func newDefaultJsonConverter() *JsonConverter {
	return &JsonConverter{
		stk:       util.Stack{},
		sb:        strings.Builder{},
		root:      nil,
		omitEmpty: true,
	}
}

// Node node of json tree
type Node struct {
	jsonDataType jsonparser.ValueType
	goDataType   string
	key          string
	child        []*Node
}

func (j *JsonConverter) buildJsonTree(data []byte) {
	j.root = buildJsonTree(j.root, data, "")
}

func buildJsonTree(root *Node, data []byte, parentKey string) *Node {
	if root == nil {
		root = &Node{}
	}

	var (
		nodeValue []byte
	)

	root.key = parentKey
	if parentKey == "" {
		root.key = "AutoGenerated"
		nodeValue, root.jsonDataType, _, _ = jsonparser.Get(data)
	} else {
		nodeValue, root.jsonDataType, _, _ = jsonparser.Get(data, parentKey)
		root.goDataType = parse2GoType(string(nodeValue), root.jsonDataType)
	}

	jsonparser.ObjectEach(nodeValue, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		var c *Node
		if dataType == jsonparser.Array {
			c = &Node{
				jsonDataType: jsonparser.Array,
				key:          string(key),
				child:        nil,
			}

			var cc = &Node{}

			// use first element
			var found = false
			jsonparser.ArrayEach(value, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
				if found {
					return
				}
				cc = buildJsonTree(nil, value, "")
				found = true
			})

			c.child = append(c.child, cc.child...)

		} else {
			c = buildJsonTree(nil, nodeValue, string(key))
		}
		root.child = append(root.child, c)
		return nil
	})

	return root
}

func (j *JsonConverter) GenStruct(jsonStr string) string {
	j.buildJsonTree([]byte(jsonStr))
	j.dfsTravelJsonTree(j.root)
	data, _ := format.Source([]byte(j.sb.String()))
	return string(data)
}

func (j *JsonConverter) dfsTravelJsonTree(root *Node) {
	if root == nil {
		return
	}

	var (
		shouldComplete = false
	)

	// visit root
	switch root.jsonDataType {
	case jsonparser.Object:
		if root.key == "AutoGenerated" {
			j.sb.WriteString("type ")
		}
		j.sb.WriteString(fmt.Sprintf("%s struct {\n", parse2GoFiledName(root.key)))
		if root.key == "AutoGenerated" {
			j.stk.Push("}\n")
		} else {
			j.stk.Push(fmt.Sprintf("} %s\n", genJsonTag(root.key, j.omitEmpty)))
		}
		shouldComplete = true
	case jsonparser.Array:
		j.sb.WriteString(fmt.Sprintf("%s [] struct {\n", parse2GoFiledName(root.key)))
		j.stk.Push(fmt.Sprintf("} %s\n", genJsonTag(root.key, j.omitEmpty)))
		shouldComplete = true
	default:
		j.sb.WriteString(fmt.Sprintf("%v %v %v\n", parse2GoFiledName(root.key),
			root.goDataType, genJsonTag(root.key, j.omitEmpty)))
	}

	// visit child
	if len(root.child) != 0 {
		for _, node := range root.child {
			j.dfsTravelJsonTree(node)
		}
	}

	// after visit
	if shouldComplete {
		if s, ok := j.stk.Pop(); ok {
			j.sb.WriteString(s.(string))
		}
	}
}

func genJsonTag(key string, omitEmpty bool) string {
	var sb strings.Builder
	sb.WriteString("`json:")
	sb.WriteString(`"`)
	sb.WriteString(key)

	if omitEmpty {
		sb.WriteString(",")
		sb.WriteString("omitempty")
	}

	sb.WriteString(`"`)
	sb.WriteString("`")

	return sb.String()
}

func parse2GoType(value string, tp jsonparser.ValueType) string {
	switch tp {
	case jsonparser.Number:
		var (
			n int64
			e error
		)

		if _, e = fmt.Sscanf(value, "%d", &n); e == nil {
			if fmt.Sprintf("%v", n) != value {
				return "float64"
			}
			if n > math.MinInt32 && n < math.MaxInt32 {
				return "int32"
			}
			return "int64"
		}

		return "float64"
	case jsonparser.String:
		if _, err := time.Parse(time.RFC3339, value); err == nil {
			return "time.Time"
		}

		return "string"
	case jsonparser.Boolean:
		return "bool"
	}

	return "interface{}"
}

func parse2GoFiledName(key string) string {
	return ToCamelCase(key)
}

// ToCamelCase ??????????????????
func ToCamelCase(str string) string {
	str = strings.TrimSpace(str)
	if utf8.RuneCountInString(str) < 2 {
		return str
	}

	if strings.Contains(str, "_") {
		return snakesToCamel(str)
	}

	var buff strings.Builder
	var temp string
	for _, r := range str {
		c := string(r)
		if c != " " {
			if temp == " " {
				c = strings.ToUpper(c)
			}
			_, _ = buff.WriteString(c)
		}
		temp = c
	}

	res := buff.String()
	b := []byte(res)
	if len(b) != 0 {
		if b[0] >= byte('a') {
			b[0] -= byte('a' - 'A')
		}
	}

	return string(b)
}

func snakesToCamel(str string) (result string) {
	part := strings.Split(str, "_")
	sb := strings.Builder{}
	for _, p := range part {
		p = strings.ToLower(p)
		p = strings.ToUpper(string(p[0:1])) + p[1:]
		sb.WriteString(p)
	}

	return sb.String()
}
