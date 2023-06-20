package json

import (
	"testing"
)

func TestSave(t *testing.T) {
	var v JSONObject
	Parse(&v,
		`
{
	"hello": "131233",
	"test": [1,2,3]
}
	`)
	err := Save(&v, "./json.txt", output_json)
	if err != nil {
		t.Errorf("Failed to save file: %v", err)
	}
}

func TestParseFile(t *testing.T) {
	var v JSONObject
	err := ParseFile(&v, "./json.txt")
	if err != nil {
		t.Errorf("Failed to parse file: %v", err)
	}
}

func TestParse(t *testing.T) {
	type args struct {
		value *JSONObject
		json  string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "TestParse",
			args: args{
				value: &JSONObject{},
				json:  "{\"hello\": \"world\"}",
			},
			want: "world",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Parse(tt.args.value, tt.args.json); got != 0 {
				t.Errorf("Parse() Err")
			}
			v := tt.args.value.Get("hello").GetString()
			if v != tt.want {
				t.Errorf("Parse() = %v, want %v", v, tt.want)
			}
		})
	}
}

func TestStringify(t *testing.T) {
	type args struct {
		value *JSONObject
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "TestParse",
			args: args{
				value: &JSONObject{vType: type_object, obj: map[string]*JSONObject{"hello": &JSONObject{
					vType: type_string,
					value: []byte("world")}},
				},
			},
			want: "{\"hello\":\"world\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Stringify(tt.args.value, ""); got != tt.want {
				t.Errorf("Stringify() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewDemo(t *testing.T) {
	tests := []struct {
		name string
		want *JSONObject
	}{
		{
			"test",
			&JSONObject{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDemo(); got == nil {
				t.Errorf("NewDemo() 不应该为nil")
			}
		})
	}
}
