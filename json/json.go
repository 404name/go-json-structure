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
func Stringify(value *JSONObject, outputType string) string {
	if value == nil {
		return "null"
	}
	var c easyContext
	easyStringifyValue(&c, value, outputType)
	return string(c.json)
}

func Save(value *JSONObject, filePath string, outputType string) error {
	err := ioutil.WriteFile(filePath, []byte(Stringify(value, outputType)), 0644)
	if err != nil {
		return err
	}
	return nil
}

func NewDemo() *JSONObject {
	var data string = `
{
	"code":"200",
	"msg":"ok",
	"data": {
		"date": "2023年6月19日",
		"city": "上海,北京广州",
		"details": [
			{"temp":"20℃/30℃","weather":"晴转多云","name":"上海","pm":"80","wind":"1级"},
			{"temp":"15℃/24℃","weather":"晴","name":"北京","pm":"98","wind":"3级"},
			{"temp":"26℃/32℃","weather":"多云","name":"广州","pm":"30","wind":"2级"}
		]
	}
}
	`
	var json JSONObject
	Parse(&json, data)
	return &json
}
