package json

import (
	"fmt"
	"strconv"
	"strings"
)

type easyContext struct {
	json []byte
}

// 去掉json前面的空格
func easyParseWhitespace(c *easyContext) {
	if len(c.json) == 0 {
		return
	}
	b := c.json[0]
	for b == ' ' || b == '\t' || b == '\n' || b == '\r' {
		c.json = c.json[1:]
		if len(c.json) == 0 {
			break
		}
		b = c.json[0]
	}
}

func easyParseValue(c *easyContext, value *JSONObject) int {
	if len(c.json) == 0 {
		return parse_expect_value
	}
	//easyParseWhitespace(c)
	switch c.json[0] {
	case 'n':
		return easyParseLiteral(c, value, "null", type_null)
	case 't':
		return easyParseLiteral(c, value, "true", type_bool)
	case 'f':
		return easyParseLiteral(c, value, "false", type_bool)
	case '"':
		return easyParseString(c, value)
	case '[':
		return easyParseArray(c, value)
	case '{':
		return easyParseObject(c, value)
	default:
		return easyParseNum(c, value)
	}
}

func easyParseObject(c *easyContext, value *JSONObject) int {
	var result int
	c.json = c.json[1:]
	easyParseWhitespace(c)

	//todo: 如果走进了这个if，意味着传入了一个空对象，EasyValue中的值是空的
	if c.json[0] == '}' {
		c.json = c.json[1:]
		value.vType = type_object
		return parse_ok
	}
	for {
		easyParseWhitespace(c)
		var k string
		var v JSONObject

		/* parse key */
		if len(c.json) < 1 || c.json[0] != '"' {
			result = parse_miss_key
			break
		}
		//todo 把字符串保存起来
		if result, k = easyParseStringRaw(c); result != parse_ok {
			break
		}

		/* parse : */
		easyParseWhitespace(c)
		if c.json[0] != ':' {
			result = parse_miss_colon
			break
		}
		c.json = c.json[1:]

		/* parse value */
		easyParseWhitespace(c)
		if result = easyParseValue(c, &v); result != parse_ok {
			break
		}
		if value.obj == nil {
			value.obj = make(map[string]*JSONObject)
		}
		value.obj[k] = &v
		easyParseWhitespace(c)
		if len(c.json) < 1 {
			return parse_miss_comma_or_curly_bracket
		} else {
			if c.json[0] == ',' {
				c.json = c.json[1:]
			} else if c.json[0] == '}' {
				c.json = c.json[1:]
				value.vType = type_object
				return parse_ok
			} else {
				return parse_miss_comma_or_curly_bracket
			}
		}
	}
	return result
}

func easyParseArray(c *easyContext, value *JSONObject) int {
	c.json = c.json[1:]
	easyParseWhitespace(c)

	//todo: 如果走进了这个if，意味着传入了一个空数组，EasyValue中的值是空的
	if c.json[0] == ']' {
		c.json = c.json[1:]
		value.vType = type_array
		return parse_ok
	}
	for {
		easyParseWhitespace(c)
		var v JSONObject
		if result := easyParseValue(c, &v); result != parse_ok {
			return result
		}

		value.list = append(value.list, &v)
		easyParseWhitespace(c)

		if len(c.json) < 1 {
			return parse_miss_comma_or_square_bracket
		} else {
			if c.json[0] == ',' {
				c.json = c.json[1:]
			} else if c.json[0] == ']' {
				c.json = c.json[1:]
				value.vType = type_array
				return parse_ok
			} else {
				return parse_miss_comma_or_square_bracket
			}
		}
	}
}

func easyParseStringRaw(c *easyContext) (int, string) {
	var result []byte
	c.json = c.json[1:] //去掉第一个"
	i := 0
	for i < len(c.json) {
		switch c.json[i] {
		case '"':
			c.json = c.json[i+1:]
			return parse_ok, string(result)
		case '\\':
			i++
			switch c.json[i] {
			case '"':
				result = append(result, '"')
			case '\\':
				result = append(result, '\\')
			case '/':
				result = append(result, '/')
			case 'b':
				result = append(result, '\b')
			case 'f':
				result = append(result, '\f')
			case 'n':
				result = append(result, '\n')
			case 'r':
				result = append(result, '\r')
			case 't':
				result = append(result, '\t')
			case 'u':
				ok, u := easyParseHex4(c, i)
				if !ok {
					return parse_invalid_unicode_hex, ""
				}
				i += 4 //指针往后移动4格

				//判断是否要继续解析，因为有可能还有一个低代理项（也就是有连续两个/uxxxx/uxxxx）
				if u >= 0xD800 && u <= 0xDBFF {
					i++
					if c.json[i] != '\\' {
						return parse_invalid_unicode_surrogate, ""
					}
					i++
					if c.json[i] != 'u' {
						return parse_invalid_unicode_surrogate, ""
					}
					i++
					ok2, u2 := easyParseHex4(c, i-1)
					if !ok2 {
						return parse_invalid_unicode_hex, ""
					}
					if u2 < 0xDC00 || u2 > 0xDFFF {
						return parse_invalid_unicode_surrogate, ""
					}
					u = (((u - 0xD800) << 10) | (u2 - 0xDC00)) + 0x10000 //计算unicode码点
					i += 3
				}
				result = append(result, easyParseUtf8(u)...)

			default:
				return parse_invalid_string_escape, ""
			}
		default:
			if c.json[i] < 0x20 {
				return parse_invalid_string_char, ""
			}
			result = append(result, c.json[i])
		}
		i++
	}
	return parse_miss_quotation_mark, ""
}

func easyParseString(c *easyContext, value *JSONObject) int {
	result, str := easyParseStringRaw(c)
	//其实不用判断也行
	if result != parse_ok {
		return result
	} else {
		value.vType = type_string
		value.value = []byte(str)
		return result
	}
}

func easyParseUtf8(u int64) []byte {
	var tail []byte
	if u <= 0x7f {
		tail = append(tail, byte(u))
	} else if u <= 0x7FF {
		tail = append(tail, byte(0xc0|((u>>6)&0xff)))
		tail = append(tail, byte(0x80|(u&0x3f)))
	} else if u <= 0xFFFF {
		tail = append(tail, byte(0xe0|((u>>12)&0xff)))
		tail = append(tail, byte(0x80|((u>>6)&0x3f)))
		tail = append(tail, byte(0x80|(u&0x3f)))
	} else if u <= 0x10FFFF {
		tail = append(tail, byte(0xf0|((u>>18)&0xff)))
		tail = append(tail, byte(0x80|((u>>12)&0x3f)))
		tail = append(tail, byte(0x80|((u>>6)&0x3f)))
		tail = append(tail, byte(0x80|u&0x3f))
	}
	return tail
}

func easyParseHex4(c *easyContext, index int) (bool, int64) {
	part := c.json[index+1 : index+5]
	i, err := strconv.ParseInt(string(part), 16, 64)
	if err != nil {
		return false, 0
	}
	return true, i
}

func easyParseNum(c *easyContext, value *JSONObject) int {
	if c.json[0] == '0' {
		//下面这个if判断是用来过滤0后面只能跟.eE三种字符，否则不表示数字，但是在解释数组的时候这样的判断出现了麻烦
		//if len(c.json) > 1 && !(c.json[1] == '.' || c.json[1] == 'e' || c.json[1] == 'list') {
		//	return parse_root_not_singular
		//}
		if len(c.json) > 1 && (c.json[1] == 'x' || (c.json[1] >= '0' && c.json[1] <= '9')) {
			return parse_root_not_singular
		}
	}

	if c.json[0] == '+' || c.json[0] == '.' || c.json[len(c.json)-1] == '.' || startWithLetter(c.json[0]) {
		return parse_invalid_value
	}

	//下面这段出问题，用","来划分数字的结束不靠谱，还是一个个数吧
	////在把字节数组转为数字之前要判断后面是否有","，否则解析数组json串的时候会有问题
	//index := strings.Index(string(c.json), ",")
	//convStr := string(c.json)
	//if index >= 0 && index <= len(c.json) {
	//	convStr = string(c.json[0:index])
	//}
	//convStr = strings.TrimSpace(convStr)
	index := 0
	for index < len(c.json) {
		if isDigit(c.json[index]) {
			index++
		} else {
			break
		}
	}
	convStr := string(c.json[0:index])
	f, err := strconv.ParseFloat(convStr, 64)
	if err != nil {
		//todo 错误判断处理是否有更好的方法
		if strings.Contains(err.Error(), strconv.ErrRange.Error()) {
			return parse_number_too_big
		}
		return parse_invalid_value
	}
	c.json = c.json[len(convStr):]
	value.num = f
	value.vType = type_number
	return parse_ok
}

func isDigit(b byte) bool {
	if (b >= '0' && b <= '9') || b == 'e' || b == 'E' || b == '-' || b == '+' || b == '.' {
		return true
	} else {
		return false
	}
}

func startWithLetter(b byte) bool {
	if b >= '0' && b <= '9' || b == '-' {
		return false
	} else {
		return true
	}
}

func easyParseLiteral(c *easyContext, value *JSONObject, target string, valueType int) int {
	ts := []byte(target)
	tsLen := len(ts)

	if len(c.json) < tsLen {
		return parse_invalid_value
	}
	for i := 0; i < tsLen; i++ {
		if c.json[i] != ts[i] {
			return parse_invalid_value
		}
	}
	c.json = c.json[tsLen:]
	value.vType = valueType
	value.value = ts
	return parse_ok
}

func easyStringifyValue(c *easyContext, value *JSONObject) {
	c.json = getByteSliceJudgeByType(value)
}

func stringifyObject(value *JSONObject) []byte {
	var result []byte
	result = append(result, '{')
	count := 0
	for k, v := range value.obj {
		result = append(result, stringifyString([]byte(k))...)
		result = append(result, ':')
		result = append(result, getByteSliceJudgeByType(v)...)

		if count < len(value.obj)-1 {
			result = append(result, ',')
		}
		count++
	}
	result = append(result, '}')
	return result
}

func stringifyArray(value *JSONObject) []byte {
	var result []byte
	result = append(result, '[')
	for i, v := range value.list {
		result = append(result, getByteSliceJudgeByType(v)...)
		if i != len(value.list)-1 {
			result = append(result, ',')
		}
	}
	result = append(result, ']')
	return result
}

func getByteSliceJudgeByType(value *JSONObject) []byte {
	var result []byte
	switch value.vType {
	case type_null:
		result = append(result, "null"...)
	case type_bool:
		result = append(result, stringifyString(value.value)...)
	case type_number:
		result = append(result, []byte(fmt.Sprintf("%.17g", value.num))...)
	case type_string:
		result = append(result, stringifyString(value.value)...)
	case type_array:
		result = append(result, stringifyArray(value)...)
	case type_object:
		result = append(result, stringifyObject(value)...)
	}
	return result
}

func stringifyString(str []byte) []byte {
	var result []byte
	result = append(result, '"')
	for _, b := range str {
		switch b {
		case '"':
			result = append(result, '\\', '"')
		case '\\':
			result = append(result, '\\', '\\')
		case '\n':
			result = append(result, '\\', 'n')
		case '\b':
			result = append(result, '\\', 'b')
		case '\f':
			result = append(result, '\\', 'f')
		case '\r':
			result = append(result, '\\', 'r')
		case '\t':
			result = append(result, '\\', 't')
		default:
			if b < 0x20 {
				result = append(result, fmt.Sprintf("\\u%04X", b)...)
			} else {
				result = append(result, b)
			}
		}
	}
	result = append(result, '"')
	return result
}
