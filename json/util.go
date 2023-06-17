package json

import (
	"errors"
	"fmt"
	"io/ioutil"
)

/* ===== 解析器-从文件解析 ===== */
func ParseFile(value *JSONObject, filePath string) error {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		// 地址问题
		return err
	}

	result := Parse(value, string(content))
	if result != 0 {
		return errors.New("解析错误-错误代码" + fmt.Sprint(result))
	}
	return nil
}

/* ===== 解析器-从string解析 ===== */
func Parse(value *JSONObject, json string) int {
	var c easyContext
	c.json = []byte(json)
	value.vType = type_null
	easyParseWhitespace(&c)
	result := easyParseValue(&c, value)
	if result == parse_ok {
		easyParseWhitespace(&c)
		if len(c.json) > 0 {
			return parse_root_not_singular
		}
	}
	return result
}

/* ===== 生成器 ===== */
func Stringify(value *JSONObject) string {
	var c easyContext
	easyStringifyValue(&c, value)
	return string(c.json)
}

func Save(value *JSONObject, filePath string) error {
	err := ioutil.WriteFile(filePath, []byte(Stringify(value)), 0644)
	if err != nil {
		return err
	}
	return nil
}
