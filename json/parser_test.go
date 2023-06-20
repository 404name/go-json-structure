package json

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

var mainRet = 0
var testCount = 0
var testPass = 0

//----- 测试对象 -----

func TestParseObject(t *testing.T) {
	var e JSONObject
	expectEQInt(parse_ok, Parse(&e, " { } "))
	expectEQInt(type_object, e.vType)

	var v JSONObject
	expectEQInt(parse_ok, Parse(&v,
		" { "+
			"\"n\" : null , "+
			"\"f\" : false , "+
			"\"t\" : true , "+
			"\"i\" : 123 , "+
			"\"s\" : \"abc\", "+
			"\"a\" : [ 1, 2, 3 ],"+
			"\"o\" : { \"1\" : 1, \"2\" : 2, \"3\" : 3 }"+
			" } "))

	assert.Equal(t, type_object, v.vType, "they should be equal")
	assert.Equal(t, type_null, v.obj["n"].vType, "they should be equal")
	assert.NotNil(t, v.obj["t"], "t应该能获取到对象")
	assert.Nil(t, v.obj["noKey"], "noKey应该为空")
	assert.Equal(t, "abc", string(v.obj["s"].value), "they should be equal")
	assert.Equal(t, 3, len(v.obj["a"].list), "they should be equal")
	assert.Equal(t, 3, len(v.obj["o"].obj), "they should be equal")
}

func TestMissKey(t *testing.T) {
	ParseExpectValue(parse_miss_key, "{:1,", type_null)
	ParseExpectValue(parse_miss_key, "{1:1,", type_null)
	ParseExpectValue(parse_miss_key, "{true:1,", type_null)
	ParseExpectValue(parse_miss_key, "{false:1,", type_null)
	ParseExpectValue(parse_miss_key, "{null:1,", type_null)
	ParseExpectValue(parse_miss_key, "{[]:1,", type_null)
	ParseExpectValue(parse_miss_key, "{{}:1,", type_null)
	ParseExpectValue(parse_miss_key, "{\"a\":1,", type_null)
}

func TestParseMissColon(t *testing.T) {
	ParseExpectValue(parse_miss_colon, "{\"a\"}", type_null)
	ParseExpectValue(parse_miss_colon, "{\"a\",\"b\"}", type_null)
}

func TestParseMissCommaOrCurlyBracket(t *testing.T) {
	ParseExpectValue(parse_miss_comma_or_curly_bracket, "{\"a\":1", type_null)
	ParseExpectValue(parse_miss_comma_or_curly_bracket, "{\"a\":1]", type_null)
	ParseExpectValue(parse_miss_comma_or_curly_bracket, "{\"a\":1 \"b\"", type_null)
	ParseExpectValue(parse_miss_comma_or_curly_bracket, "{\"a\":{}", type_null)
}

//----- 测试数组 -----

func TestParseArray(t *testing.T) {
	var v JSONObject
	expectEQInt(parse_ok, Parse(&v, "[ null , false , true , 123 , \"abc\" ]"))
	expectEQInt(type_array, v.vType)
	expectEQInt(type_null, v.list[0].vType)
	expectEQInt(type_bool, v.list[1].vType)
	expectEQInt(type_bool, v.list[2].vType)
	expectEQInt(type_number, v.list[3].vType)
	expectEQString("abc", string(v.list[4].value))

	var v2 JSONObject
	expectEQInt(parse_ok, Parse(&v2, "[ [ ] , [ 0 ] , [ 0 , 1 ] , [ 0 , 1 , 2 ] ]"))
	expectEQInt(type_array, v2.list[0].vType)
	expectEQInt(3, len(v2.list[3].list))
	expectEQInt(2, int(v2.list[3].list[2].num))
}

func TestParseMissCommaOrSquareBracket(t *testing.T) {
	ParseExpectValue(parse_miss_comma_or_square_bracket, "[1", type_null)
	ParseExpectValue(parse_miss_comma_or_square_bracket, "[1}", type_null)
	ParseExpectValue(parse_miss_comma_or_square_bracket, "[1 2", type_null)
	ParseExpectValue(parse_miss_comma_or_square_bracket, "[[]", type_null)
}

//----- 测试字符串 -----

func TestParseString(t *testing.T) {
	testString("", "\"\"")
	testString("Hello", "\"Hello\"")
	testString("Hello\nWorld", "\"Hello\\nWorld\"")
	testString("\" \\ / \b \f \n \r \t", "\"\\\" \\\\ \\/ \\b \\f \\n \\r \\t\"")
	testString("\x24", "\"\\u0024\"")
	testString("\xC2\xA2", "\"\\u00A2\"")
	testString("\xF0\x9D\x84\x9E", "\"\\uD834\\uDD1E\"")
	testString("\xF0\x9D\x84\x9E", "\"\\ud834\\udd1e\"")
}

func TestParseMissingQuetationMark(t *testing.T) {
	ParseExpectValue(parse_miss_quotation_mark, "\"", type_null)
	ParseExpectValue(parse_miss_quotation_mark, "\"abc", type_null)
}

func TestParseInvalidStringEscape(t *testing.T) {
	ParseExpectValue(parse_invalid_string_escape, "\"\\v\"", type_null)
	ParseExpectValue(parse_invalid_string_escape, "\"\\'\"", type_null)
	ParseExpectValue(parse_invalid_string_escape, "\"\\0\"", type_null)
	ParseExpectValue(parse_invalid_string_escape, "\"\\x12\"", type_null)
}

func TestParseInvalidStringChar(t *testing.T) {
	ParseExpectValue(parse_invalid_string_char, "\"\x01\"", type_null)
	ParseExpectValue(parse_invalid_string_char, "\"\x1F\"", type_null)
}

func TestParseInvalidUnicodeHex(t *testing.T) {
	ParseExpectValue(parse_invalid_unicode_hex, "\"\\u\"", type_null)
	ParseExpectValue(parse_invalid_unicode_hex, "\"\\u0\"", type_null)
	ParseExpectValue(parse_invalid_unicode_hex, "\"\\u01\"", type_null)
	ParseExpectValue(parse_invalid_unicode_hex, "\"\\u012\"", type_null)
	ParseExpectValue(parse_invalid_unicode_hex, "\"\\u/000\"", type_null)
	ParseExpectValue(parse_invalid_unicode_hex, "\"\\uG000\"", type_null)
	ParseExpectValue(parse_invalid_unicode_hex, "\"\\u0/00\"", type_null)
	ParseExpectValue(parse_invalid_unicode_hex, "\"\\u0G00\"", type_null)
	ParseExpectValue(parse_invalid_unicode_hex, "\"\\u00/0\"", type_null)
	ParseExpectValue(parse_invalid_unicode_hex, "\"\\u00G0\"", type_null)
	ParseExpectValue(parse_invalid_unicode_hex, "\"\\u000/\"", type_null)
	ParseExpectValue(parse_invalid_unicode_hex, "\"\\u000G\"", type_null)
	ParseExpectValue(parse_invalid_unicode_hex, "\"\\u 123\"", type_null)
}

func TestParseInvalidUnicodeSurrogate(t *testing.T) {
	ParseExpectValue(parse_invalid_unicode_surrogate, "\"\\uD800\"", type_null)
	ParseExpectValue(parse_invalid_unicode_surrogate, "\"\\uDBFF\"", type_null)
	ParseExpectValue(parse_invalid_unicode_surrogate, "\"\\uD800\\\\\"", type_null)
	ParseExpectValue(parse_invalid_unicode_surrogate, "\"\\uD800\\uDBFF\"", type_null)
	ParseExpectValue(parse_invalid_unicode_surrogate, "\"\\uD800\\uE000\"", type_null)
}

func testString(expect string, json string) {
	var v JSONObject
	expectEQInt(parse_ok, Parse(&v, json))
	expectEQInt(type_string, v.Type())
	expectEQString(expect, string(v.value))
}

//----- 测试数字 -----

func TestParseNum(t *testing.T) {
	testNumber(123, "123")
	testNumber(0.0, "0")
	testNumber(0.0, "-0")
	testNumber(0.0, "-0.0")
	testNumber(1.0, "1")
	testNumber(-1.0, "-1")
	testNumber(1.5, "1.5")
	testNumber(-1.5, "-1.5")
	testNumber(3.1416, "3.1416")
	testNumber(1e10, "1E10")
	testNumber(1e10, "1e10")
	testNumber(1e+10, "1E+10")
	testNumber(1e-10, "1E-10")
	testNumber(-1e10, "-1E10")
	testNumber(-1e10, "-1e10")
	testNumber(-1e+10, "-1E+10")
	testNumber(-1e-10, "-1E-10")
	testNumber(1.234e+10, "1.234E+10")
	testNumber(1.234e-10, "1.234E-10")
	testNumber(0.0, "1e-10000") /* must underflow */
	testNumber(4.9406564584124654e-324, "4.9406564584124654E-324")
	testNumber(1.7976931348623157e-308, "1.7976931348623157e-308")
}

func TestParseNumToBig(t *testing.T) {
	ParseExpectValue(parse_number_too_big, "1e309", type_null)
	ParseExpectValue(parse_number_too_big, "-1e309", type_null)
}

func testNumber(num float64, json string) {
	var v JSONObject
	expectEQInt(parse_ok, Parse(&v, json))
	expectEQInt(type_number, v.Type())
	expectEQFloat(num, v.num)
}

//----- 测试null/true/false -----

func TestParseNull(t *testing.T) {
	ParseExpectValue(parse_ok, "null", type_null)
}

func TestParseTrue(t *testing.T) {
	ParseExpectValue(parse_ok, "true", type_bool)
}

func TestParseFalse(t *testing.T) {
	ParseExpectValue(parse_ok, "false", type_bool)
}

func TestParseExpectValue(t *testing.T) {
	ParseExpectValue(parse_expect_value, "", type_null)
	ParseExpectValue(parse_expect_value, " ", type_null)
}

func TestParseInvalidValue(t *testing.T) {
	ParseExpectValue(parse_invalid_value, "nul", type_null)
	ParseExpectValue(parse_invalid_value, "?", type_null)

	/* invalid number */
	ParseExpectValue(parse_invalid_value, "+0", type_null)
	ParseExpectValue(parse_invalid_value, "+1", type_null)
	ParseExpectValue(parse_invalid_value, ".123", type_null)
	ParseExpectValue(parse_invalid_value, "1.", type_null)
	ParseExpectValue(parse_invalid_value, "INF", type_null)
	ParseExpectValue(parse_invalid_value, "inf", type_null)
	ParseExpectValue(parse_invalid_value, "NAN", type_null)
	ParseExpectValue(parse_invalid_value, "nan", type_null) //fixed in judge if null
}

func TestParseRootNotSingular(t *testing.T) {
	ParseExpectValue(parse_root_not_singular, "null x", type_null)
	ParseExpectValue(parse_root_not_singular, "0123", type_null)
	ParseExpectValue(parse_root_not_singular, "0x0", type_null)
	ParseExpectValue(parse_root_not_singular, "0x123", type_null)
}

/*----- 测试生成器 -----*/

func TesteasyStringifyValue(t *testing.T) {
	testRoundTrip("null")
	testRoundTrip("true")
	testRoundTrip("false")
}

func TestEasyStringifyNumber(t *testing.T) {
	testRoundTrip("0")
	testRoundTrip("-0")
	testRoundTrip("1")
	testRoundTrip("-1")
	testRoundTrip("1.5")
	testRoundTrip("-1.5")
	testRoundTrip("3.25")
	testRoundTrip("1e+20")
	testRoundTrip("1.234e+20")
	testRoundTrip("1.234e-20")
	testRoundTrip("1.0000000000000002")      /* the smallest number > 1 */
	testRoundTrip("4.9406564584124654e-324") /* minimum denormal */
	testRoundTrip("-4.9406564584124654e-324")
	testRoundTrip("2.2250738585072009e-308") /* Max subnormal double */
	testRoundTrip("-2.2250738585072009e-308")
	testRoundTrip("2.2250738585072014e-308") /* Min normal positive double */
	testRoundTrip("-2.2250738585072014e-308")
	testRoundTrip("1.7976931348623157e+308") /* Max double */
	testRoundTrip("-1.7976931348623157e+308")
}

func TestStringifyString(t *testing.T) {
	testRoundTrip("\"\"")
	testRoundTrip("\"Hello\"")
	testRoundTrip("\"Hello\\nWorld\"")
	testRoundTrip("\"\\\" \\\\ / \\b \\f \\n \\r \\t\"")
	testRoundTrip("\"Hello\\u0000World\"")
}

func TestStringifyArray(t *testing.T) {
	testRoundTrip("[]")
	testRoundTrip("[null,false,true,123,\"abc\",[1,2,3]]")
}

func TestStringifyObject(t *testing.T) {
	testRoundTrip("{}")
	testRoundTrip("{\"n\":null,\"f\":false,\"t\":true,\"i\":123,\"s\":\"abc\",\"a\":[1,2,3],\"o\":{\"1\":1,\"2\":2,\"3\":3}}")

}

func testRoundTrip(json string) {
	var v JSONObject
	expectEQInt(parse_ok, Parse(&v, json))
	json2 := Stringify(&v, output_json)
	expectEQString(json, json2)
}

func expectEQ(equality bool, expect interface{}, actual interface{}, format string) {
	testCount++
	if equality {
		testPass++
	} else {
		log.Printf("expect: "+format+", actual: "+format+"\n", expect, actual)
	}
}

func expectEQInt(expect int, actual int) {
	expectEQ(expect == actual, expect, actual, "%d")
}

func expectEQFloat(expect float64, actual float64) {
	expectEQ(expect == actual, expect, actual, "%f")
}

func expectEQString(expect string, actual string) {
	expectEQ(expect == actual, expect, actual, "%s")
}

func ParseExpectValue(parseCode int, value string, valueType int) {
	var v JSONObject
	v.vType = type_bool
	expectEQInt(parseCode, Parse(&v, value))
	expectEQInt(valueType, v.Type())
}

func TestInfo(t *testing.T) {
	log.Printf("%d/%d (%f%%) passed, fail: %d\n",
		testPass, testCount, float32(testPass*100.0/testCount), testCount-testPass)
}
