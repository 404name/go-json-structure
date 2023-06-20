package json

//json解析器的七种类型
const (
	type_null = iota // null、false、true、数字、字符串、数组、对象
	type_bool
	type_number
	type_string
	type_array
	type_object
)

//输出格式
const (
	output_string string = ""
	output_json   string = "json"
	output_yaml   string = "yaml"
	output_xml    string = "xml"
)

//解析后的返回值
const (
	parse_ok = iota
	parse_expect_value
	parse_invalid_value
	parse_root_not_singular
	parse_number_too_big
	parse_miss_quotation_mark
	parse_invalid_string_escape
	parse_invalid_string_char
	parse_invalid_unicode_hex
	parse_invalid_unicode_surrogate
	parse_miss_comma_or_square_bracket
	parse_miss_key
	parse_miss_colon
	parse_miss_comma_or_curly_bracket
	parse_invalid_filepath // 解析路径出问题
)
