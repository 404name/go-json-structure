package json

import (
	"reflect"
	"testing"
)

func TestJSONObject_Type(t *testing.T) {
	tests := []struct {
		name string
		obj  *JSONObject
		want int
	}{
		{
			"类型获取测试null", &JSONObject{vType: type_null}, type_null,
		},
		{
			"类型获取测试string", &JSONObject{vType: type_string}, type_string,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.obj.Type(); got != tt.want {
				t.Errorf("JSONObject.Type() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJSONObject_Add(t *testing.T) {
	type args struct {
		key   string
		value *JSONObject
	}
	tests := []struct {
		name string
		obj  *JSONObject
		args args
		want string
	}{
		{
			"对象添加节点测试",
			&JSONObject{vType: type_object},
			args{
				key: "key",
				value: &JSONObject{
					vType: type_string,
					value: []byte("test")}}, "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.obj.Add(tt.args.key, tt.args.value); got != true {
				t.Errorf("JSONObject.Add() = %v, want %v", got, true)
			}
			v := tt.obj.Get(tt.args.key).GetString()
			if v != tt.want {
				t.Errorf("JSONObject.Add() = %v, want %v", v, tt.want)
			}
		})
	}
}

func TestJSONObject_Append(t *testing.T) {
	type args struct {
		jsonObj *JSONObject
	}
	tests := []struct {
		name string
		obj  *JSONObject
		args args
		want string
	}{
		{
			"数组添加节点测试",
			&JSONObject{vType: type_array},
			args{
				jsonObj: &JSONObject{
					vType: type_string,
					value: []byte("test")}}, "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.obj.Append(tt.args.jsonObj); got != true {
				t.Errorf("JSONObject.Append() = %v, want %v", got, true)
			}
			v := tt.obj.GetIndex(0).GetString()
			if v != tt.want {
				t.Errorf("JSONObject.Append() = %v, want %v", v, tt.want)
			}
		})
	}
}

func TestJSONObject_Remove(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		obj  *JSONObject
		args args
		want bool
	}{
		{
			"对象删除节点测试",
			&JSONObject{vType: type_object, obj: map[string]*JSONObject{"key": &JSONObject{
				vType: type_string,
				value: []byte("test")}}},
			args{
				key: "key",
			}, true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.obj.Get("key") == nil {
				t.Errorf("初始化应该获取到key")
			}
			if got := tt.obj.Remove(tt.args.key); got != tt.want {
				t.Errorf("JSONObject.Remove() = %v, want %v", got, tt.want)
			}
			if tt.obj.Get("key") != nil {
				t.Errorf("删除后应该拿不到key")
			}
		})
	}
}

func TestJSONObject_RemoveIndex(t *testing.T) {
	type args struct {
		idx int
	}
	tests := []struct {
		name string
		obj  *JSONObject
		args args
		want bool
	}{
		{
			"数组删除节点测试",
			&JSONObject{vType: type_array, list: []*JSONObject{&JSONObject{
				vType: type_string,
				value: []byte("test")}}},
			args{
				idx: 0,
			}, true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.obj.GetIndex(0) == nil {
				t.Errorf("初始化应该获取到key")
			}
			if got := tt.obj.RemoveIndex(tt.args.idx); got != tt.want {
				t.Errorf("JSONObject.RemoveIndex() = %v, want %v", got, tt.want)
			}
			if tt.obj.GetIndex(0) != nil {
				t.Errorf("删除后应该拿不到key")
			}
		})
	}
}

func TestJSONObject_SetString(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		obj  *JSONObject
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.obj.SetString(tt.args.str)
		})
	}
}

func TestJSONObject_SetNumber(t *testing.T) {
	type args struct {
		number float64
	}
	tests := []struct {
		name string
		obj  *JSONObject
		args args
	}{
		{
			"设置数字测试",
			&JSONObject{vType: type_number, num: 10.00},
			args{
				number: 12.33,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.obj.SetNumber(tt.args.number)
			if tt.obj.GetNumber() != tt.args.number {
				t.Errorf("SetNumber_error")
			}
		})
	}
}

func TestJSONObject_SetBool(t *testing.T) {
	type args struct {
		b    bool
		want string
	}
	tests := []struct {
		name string
		obj  *JSONObject
		args args
	}{
		{
			"设置bool测试",
			&JSONObject{vType: type_bool, value: []byte("true")},
			args{
				b:    false,
				want: "false",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.obj.SetBool(tt.args.b)
		})
		if tt.obj.GetBool() != tt.args.b {
			t.Errorf("SetBool_error")
		}
	}
}

func TestJSONObject_Get(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		obj  *JSONObject
		args args
		want *JSONObject
	}{
		{
			"get_test",
			&JSONObject{vType: type_object, obj: map[string]*JSONObject{"key": &JSONObject{
				vType: type_string,
				value: []byte("test")}}},
			args{
				key: "key",
			},
			&JSONObject{
				vType: type_string,
				value: []byte("test")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.obj.Get(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JSONObject.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJSONObject_GetString(t *testing.T) {
	tests := []struct {
		name string
		obj  *JSONObject
		want string
	}{
		{
			"TestJSONObject_GetString",
			&JSONObject{vType: type_string, value: []byte("hello")},
			"hello",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.obj.GetString(); got != tt.want {
				t.Errorf("JSONObject.GetString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJSONObject_GetNumber(t *testing.T) {
	tests := []struct {
		name string
		obj  *JSONObject
		want float64
	}{
		{
			"TestJSONObject_GetNumber",
			&JSONObject{vType: type_number, num: 12.33},
			12.33,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.obj.GetNumber(); got != tt.want {
				t.Errorf("JSONObject.GetNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJSONObject_GetBool(t *testing.T) {
	tests := []struct {
		name string
		obj  *JSONObject
		want bool
	}{
		{
			"TestJSONObject_GetBool",
			&JSONObject{vType: type_bool, value: []byte("true")},
			true,
		},
		{
			"TestJSONObject_GetBool",
			&JSONObject{vType: type_bool, value: []byte("false")},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.obj.GetBool(); got != tt.want {
				t.Errorf("JSONObject.GetBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJSONObject_GetIndex(t *testing.T) {
	type args struct {
		idx int
	}
	tests := []struct {
		name string
		obj  *JSONObject
		args args
		want *JSONObject
	}{
		{
			"TestJSONObject_GetIndex",
			&JSONObject{vType: type_array, list: []*JSONObject{&JSONObject{
				vType: type_string,
				value: []byte("test")}}},
			args{
				idx: 0,
			},
			&JSONObject{
				vType: type_string,
				value: []byte("test")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.obj.GetIndex(tt.args.idx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JSONObject.GetIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJSONObject_Value(t *testing.T) {
	tests := []struct {
		name string
		obj  *JSONObject
		want interface{}
	}{
		{
			"TestJSONObject_Value",
			&JSONObject{vType: type_bool, value: []byte("true")},
			true,
		},
		{
			"TestJSONObject_GetNumber",
			&JSONObject{vType: type_number, num: 12.33},
			12.33,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.obj.Value(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JSONObject.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}
