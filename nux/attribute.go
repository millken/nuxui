// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"unsafe"

	"github.com/nuxui/nuxui/log"
)

type Attr map[string]any

func MergeAttrs(attrs ...Attr) Attr {
	l := len(attrs)
	if l == 0 {
		return Attr{}
	}
	if l == 1 {
		return attrs[0]
	}
	attr := attrs[0]
	i := 0
	for _, atr := range attrs {
		if i == 0 {
			i++
			continue
		}
		for k, v := range atr {
			attr[k] = v
		}
	}
	return attr
}

func (me Attr) Has(key string) bool {
	_, ok := me[key]
	return ok
}

func (me Attr) Set(key string, value any) {
	me[key] = value
}

func (me Attr) Get(key string, defaultValue any) any {
	if attr, ok := me[key]; ok {
		return attr
	}
	return defaultValue
}

func (me Attr) GetAttr(key string, defaultValue Attr) Attr {
	ret := defaultValue
	if attr, ok := me[key]; ok {
		switch t := attr.(type) {
		case Attr:
			ret = t
		case map[string]any:
			ret = *(*Attr)(unsafe.Pointer(&t))
		default:
			log.E("nuxui", "Error: unsupport convert %T to Attr, use default value instead", t)

		}
	}
	return ret
}

func (me Attr) GetString(key string, defaultValue string) string {
	if attr, ok := me[key]; ok {
		switch t := attr.(type) {
		case string:
			return t
		case rune, byte, uint, uint16, uint32, uint64, int, int8, int16, int64, float32, float64, bool: //byte = uint8, rune = int32
			return fmt.Sprintf("%s", t)
		default:

			log.E("nuxui", "Error: unsupport convert %s %T:%s to string, use default value instead", key, t, t)
		}
	}
	return defaultValue
}

func (me Attr) GetBool(key string, defaultValue bool) bool {
	if attr, ok := me[key]; ok {
		switch t := attr.(type) {
		case string:
			if "true" == t {
				return true
			} else if "false" == t {
				return false
			}
			log.E("nuxui", "Error: unsupport convert %s %T:%s to bool, use default value instead", key, t, t)
		case bool:
			return t
		default:
			log.E("nuxui", "Error: unsupport convert %s %T:%s to bool, use default value instead", key, t, t)
		}
	}
	return defaultValue
}

func (me Attr) GetColor(key string, defaultValue Color) (result Color) {
	result = defaultValue
	if attr, ok := me[key]; ok {
		switch t := attr.(type) {
		case string:
			t = strings.TrimSpace(t)
			if !strings.HasPrefix(t, "#") {
				log.E("nuxui", "Error: color value of %s should start with #, use default value instead", key)
				result = defaultValue
			} else {
				t = strings.TrimPrefix(t, "#")
				if strlen(t) == 6 {
					t = "FF" + t
				}
				if strlen(t) != 8 {
					log.E("nuxui", "Error: cannot convert %s %T:%s to Color, use default value instead", key, t, t)
					result = defaultValue
				}
				v, e := strconv.ParseUint(t, 16, 32)
				if e != nil {
					log.E("nuxui", "Error: cannot convert %s %T:%s to Color, use default value instead", key, t, t)
					result = defaultValue
				} else {
					result = Color(v)
				}
			}

		case int:
			result = Color(t)
		case int64:
			result = Color(t)
		case int32:
			result = Color(t)
		case uint64:
			result = Color(t)
		case uint32:
			result = Color(t)
		default:
			log.E("nuxui", "Error: unsupport convert %s %T:%s to Color, use default value instead", key, t, t)
			result = defaultValue
		}
	}

	if result > math.MaxUint32 {
		log.E("nuxui", "Error: %s %d overflow to Color, use default value instead", key, result)
		result = defaultValue
	}

	return result
}

func (me Attr) GetUint32(key string, defaultValue uint32) (result uint32) {
	result = defaultValue
	if attr, ok := me[key]; ok {
		switch t := attr.(type) {
		case string:
			v, e := strconv.ParseFloat(t, 64)
			if e != nil {
				log.E("nuxui", "Error: cannot convert %s %T:%s to uint32, use default value instead", key, t, t)
				result = defaultValue
			} else {
				result = uint32(v)
			}
		case int:
			result = uint32(t)
		case int64:
			result = uint32(t)
		case int32:
			result = uint32(t)
		case uint64:
			result = uint32(t)
		case uint32:
			result = uint32(t)
		default:
			log.E("nuxui", "Error: unsupport convert %s %T:%s to uint32, use default value instead", key, t, t)
			result = defaultValue
		}
	}

	if result > math.MaxUint32 {
		log.E("nuxui", "Error: %s %d overflow to uint32, use default value instead", key, result)
		result = defaultValue
	}

	return result
}

func (me Attr) GetInt32(key string, defaultValue int32) int32 {
	if attr, ok := me[key]; ok {
		switch t := attr.(type) {
		case string:
			v, e := strconv.ParseFloat(t, 64)
			if e != nil {
				log.E("nuxui", "Error: cannot convert %s %T:%s to int32, use default value instead", key, t, t)
				return defaultValue
			}
			return int32(v)
		case int:
			return int32(t)
		case int64:
			return int32(t)
		case int32:
			return int32(t)
		case uint64:
			return int32(t)
		case uint32:
			return int32(t)
		default:
			log.E("nuxui", "Error: unsupport convert %s %T:%s to int32, use default value instead", key, t, t)
		}
	}
	return defaultValue
}

func (me Attr) GetInt(key string, defaultValue int) int {
	if attr, ok := me[key]; ok {
		switch t := attr.(type) {
		case string:
			v, e := strconv.ParseFloat(t, 64)
			if e != nil {
				log.E("nuxui", "Error: cannot convert %s %T:%s to int32, use default value instead", key, t, t)
				return defaultValue
			}
			return int(v)
		case int:
			return t
		case int64:
			return int(t)
		case int32:
			return int(t)
		case uint64:
			return int(t)
		case uint32:
			return int(t)
		default:
			log.E("nuxui", "Error: unsupport convert %s %T:%s to int32, use default value instead", key, t, t)
		}
	}
	return defaultValue
}

func (me Attr) GetInt64(key string, defaultValue int64) int64 {
	if attr, ok := me[key]; ok {
		switch t := attr.(type) {
		case string:
			v, e := strconv.ParseFloat(t, 64)
			if e != nil {
				log.E("nuxui", "Error: cannot convert %s %T:%s to int64, use default value instead", key, t, t)
				return defaultValue
			}

			return int64(v)
		case int:
			return int64(t)
		case int64:
			return int64(t)
		case int32:
			return int64(t)
		case uint64:
			return int64(t)
		case uint32:
			return int64(t)
		default:
			log.E("nuxui", "Error: unsupport convert %s %T:%s to int64, use default value instead", key, t, t)
		}
	}
	return defaultValue
}

func (me Attr) GetFloat32(key string, defaultValue float32) float32 {
	if attr, ok := me[key]; ok {
		switch t := attr.(type) {
		case string:
			v, e := strconv.ParseFloat(t, 64)
			if e != nil {
				log.E("nuxui", "Error: cannot convert %s %T:%s to float32, use default value instead", key, t, t)
				return defaultValue
			}

			return float32(v)
		case float64:
			return float32(t)
		case float32:
			return float32(t)
		case int:
			return float32(t)
		case int64:
			return float32(t)
		case int32:
			return float32(t)
		case uint64:
			return float32(t)
		case uint32:
			return float32(t)
		default:
			log.E("nuxui", "Error: unsupport convert %s %T:%s to float32, use default value instead", key, t, t)
		}
	}
	return defaultValue
}

func (me Attr) GetFloat64(key string, defaultValue float64) float64 {
	if attr, ok := me[key]; ok {
		switch t := attr.(type) {
		case string:
			v, e := strconv.ParseFloat(t, 64)
			if e != nil {
				log.E("nuxui", "Error: cannot convert %s %T:%s to float64, use default value instead", key, t, t)
				return defaultValue
			}

			return float64(v)
		case float64:
			return float64(t)
		case float32:
			return float64(t)
		case int:
			return float64(t)
		case int64:
			return float64(t)
		case int32:
			return float64(t)
		case uint64:
			return float64(t)
		case uint32:
			return float64(t)
		default:
			log.E("nuxui", "Error: unsupport convert %s %T:%s to float64, use default value instead", key, t, t)
		}
	}
	return defaultValue
}

func (me Attr) GetArray(key string, defaultValue []any) []any {
	if attr, ok := me[key]; ok {
		switch t := attr.(type) {
		case []any:
			return t
		default:
			log.E("nuxui", "Error: unsupport convert %s %T:%s to array, use default value instead", key, t, t)
		}
	}
	return defaultValue
}

func (me Attr) GetStringArray(key string, defaultValue []string) []string {
	if attr, ok := me[key]; ok {
		switch t := attr.(type) {
		case []string:
			return t
		case []any:
			ret := make([]string, len(t))
			for i, any := range t {
				if str, ok := any.(string); ok {
					ret[i] = str
				} else {
					log.E("nuxui", "Error: unsupport convert %s %T:%s to []string, use default value instead", key, t, t)
					return defaultValue
				}
			}
			return ret
		default:
			log.E("nuxui", "Error: unsupport convert %s %T:%s to []string, use default value instead", key, t, t)
		}
	}
	return defaultValue
}

func (me Attr) Merge(attr ...Attr) {
	for _, a := range attr {
		for k, v := range a {
			me[k] = v
		}
	}
}

func (me Attr) String() string {
	return attrToString(me, 0)
}

func attrToString(attr Attr, depth int) string {
	var ret string
	var space string

	for i := 0; i != depth; i++ {
		space += "  "
	}

	ret += space + "{\n"

	depth++
	space += "  "
	for k, v := range attr {
		if arr, ok := v.([]any); ok {
			for _, a := range arr {
				if child, ok := a.(Attr); ok {
					ret += attrToString(child, depth)
				} else {
					ret += fmt.Sprintf("%s%s: %s\n", space, k, v)
				}
			}
		} else {
			ret += fmt.Sprintf("%s%s: %s\n", space, k, v)
		}
	}
	ret += space + "}"

	return ret
}

// Dimension px
type Dimension int32

func NumeralToDimension(numeral float64) Dimension {
	return 0
}

func (me Attr) GetDimen(key string, defaultValue string) Dimen {
	if attr, ok := me[key]; ok {
		switch t := attr.(type) {
		case string:
			ret, e := ParseDimen(t)
			if e != nil {
				log.E("nuxui", "Error: Can't convert %s to Dimen, use default value instead, %s", t, e.Error())
				ret, e = ParseDimen(defaultValue)
				if e != nil {
					log.E("nuxui", "Error: Can't convert default value %s to Dimen, %s", defaultValue, e.Error())
					goto end
				}
				return ret
			}
			return ret
		case int:
			return Dimen(t)
		default:
			log.E("nuxui", "Error: unsupport convert %T to Dimen, use default value instead", t)
		}
	} else {
		ret, e := ParseDimen(defaultValue)
		if e != nil {
			log.E("nuxui", "Error: Can't convert default value %s to Dimen, %s", defaultValue, e.Error())
			goto end
		}
		return ret
	}

end:
	ret := ADimen(0, Auto)
	return ret
}

///////////////////////////////////// Parse /////////////////////////////////////

// ParseAttr parse all value to string
// print error if has error and return empty Attr
func ParseAttr(attrtext string) Attr {
	ret, err := ParseAttrWitthError(attrtext)
	if err != nil {
		log.E("nuxui", err.Error())
		return Attr{}
	}
	return ret
}

// TODO:: error
func ParseAttrWitthError(attrtext string) (Attr, error) {
	return (&attr{}).parse(attrtext), nil
}

const EOF = -1

type attr struct {
	data []rune
	len  int
	pos  int
	r    rune
	i    Attr // import
}

func (me *attr) parse(template string) Attr {
	me.data = []rune(template)
	me.len = len(me.data)
	me.pos = -1
	me.r = EOF

	me.next()
	me.skipBlank()
	if me.r != '{' {
		log.Fatal("nuxui", "first element is not '{'.")
	}
	return me.nextStruct()
}

func (me *attr) nextStruct() Attr {
	data := Attr{}

	for {
		me.next()
		me.skipBlank()
		if me.r == '}' || me.r == EOF {
			break
		}
		if me.r == ',' {
			continue
		}

		key := me.nextKey()
		me.next()
		me.skipBlank()
		value := me.nextValue()
		if debug_attr {
			if _, ok := data[key]; ok {
				log.E("nuxui", "the attribute key %s is already exist.", key)
			}
		}
		data[key] = value

		if me.i == nil && "import" == key {
			if i, ok := value.(Attr); ok {
				me.i = i
			}
		}
	}

	return data
}

func (me *attr) nextValue() any {
	switch me.r {
	case '"', '`':
		return me.nextString(me.r)
	case '{':
		return me.nextStruct()
	case '[':
		return me.nextArray()
	default:
		str := me.nextString(',')
		if me.i != nil {
			strs := strings.Split(str, ".")
			if len(strs) == 2 {
				if v, ok := me.i[strs[0]]; ok {
					return fmt.Sprintf("%s.%s", v, strs[1])
				}
			}
		}
		return str
	}
	return nil
}

func (me *attr) nextArray() []any {
	arr := make([]any, 0)
	for {
		me.next()
		me.skipBlank()
		if me.r == ']' {
			break
		}
		if me.r == ',' {
			continue
		}
		value := me.nextValue()
		arr = append(arr, value)
	}
	return arr
}

func (me *attr) skipBlank() {
	for {
		switch me.r {
		case ' ', '\r', '\n', '\t', '\f', '\b', 65279:
			me.next()
			continue
		case '/':
			me.skipComment()
			continue
		default:
			break
		}
		break
	}
}

func (me *attr) next() rune {
	me.pos++
	if me.pos >= me.len {
		me.r = EOF
	} else {
		me.r = me.data[me.pos]
	}
	return me.r
}

func (me *attr) previous() rune {
	me.pos--
	if me.pos >= me.len || me.pos < 0 {
		me.r = EOF
	} else {
		me.r = me.data[me.pos]
	}
	return me.r
}

func (me *attr) skipComment() {
	me.next()
	if me.r == '/' {
		for {
			me.next()
			if me.r == '\n' {
				me.next()
				return
			} else if me.r == EOF {
				return
			}
		}
	} else if me.r == '*' {
		me.next()
		for me.r != EOF {
			if me.r == '*' {
				me.next()
				if me.r == '/' {
					me.next()
					return
				}
				continue
			}
			me.next()
		}
	} else {
		log.Fatal("nuxui", "invalid comment")
	}
}

// obtain next key, and last rune is ':'
func (me *attr) nextKey() string {
	p := me.pos
	if me.r < 'A' || (me.r > 'Z' && me.r < 'a') || me.r > 'z' {
		log.Fatal("nuxui", "Invalid key name, the first letter must be [A-Za-z]")
	}
	me.next()
	for (me.r >= 'a' && me.r <= 'z') || (me.r >= 'A' && me.r <= 'Z') || (me.r >= '0' && me.r <= '9') || me.r == '_' {
		me.next()
	}
	key := string(me.data[p:me.pos])
	me.skipBlank()
	if me.r != ':' {
		log.Fatal("nuxui", "Invalid format after key name %s, missing ':'", key)
	}
	return key
}

func (me *attr) nextString(end rune) string {
	p := me.pos
	ret := ""
	quot := (end == '"' || end == '`')

	for {
		me.next()
		if quot {
			switch me.r {
			case '"', '`':
				if me.data[me.pos-1] != '\\' {
					ret = string(me.data[p+1 : me.pos])
					goto out
				} else {
					// TODO
				}

			case EOF:
				log.Fatal("nuxui", "unclosed string")
			}
		} else {
			switch me.r {
			case ',', '}', ']':
				if me.data[me.pos-1] != '\\' {
					ret = strings.TrimSpace(string(me.data[p:me.pos]))
					me.previous()
					goto out
				} else {
					// TODO
				}

			case EOF:
				log.Fatal("nuxui", "unclosed string")
			}
		}

	}

out:
	return ret
}

func strlen(text string) int {
	return len([]rune(text))
}
