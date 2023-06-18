package json

import "errors"

//json树的节点
type JSONObject struct {
	vType int
	num   float64
	value []byte
	list  []*JSONObject
	obj   map[string]*JSONObject
}

// TODO 添加error抛出（大工程。。。）

// 查对象key
func (obj *JSONObject) Type() int {
	return obj.vType
}

// 查对象key
func (obj *JSONObject) IsArray() bool {
	return obj.vType == type_array
}

// 查对象key
func (obj *JSONObject) IsObject() bool {
	return obj.vType == type_object
}

// ======= 增 ==========、
// 仅限对象
func (obj *JSONObject) Add(key string, value *JSONObject) bool {

	if obj.obj == nil {
		obj.obj = make(map[string]*JSONObject)
	}

	if obj.vType == type_object {
		obj.obj[key] = value
		return true
	}
	return false
}

// 仅限列表
func (obj *JSONObject) Append(jsonObj *JSONObject) bool {
	if obj.list == nil {
		obj.list = make([]*JSONObject, 0)
	}

	if obj.vType == type_array {
		obj.list = append(obj.list, jsonObj)
		return true
	}
	return false
}

// ======= 删 ==========
func (obj *JSONObject) Remove(key string) bool {
	if obj.vType == type_object {
		_, ok := obj.obj[key]
		if ok == true {
			delete(obj.obj, key)
		}
		return ok
	}
	return false
}

func (obj *JSONObject) RemoveIndex(idx int) bool {
	if obj.vType == type_array && idx >= 0 && idx < len(obj.list) {
		obj.list = append(obj.list[:idx], obj.list[idx+1:]...)
		return true
	}
	return false
}

// ====== 改 =========
func (obj *JSONObject) SetString(str string) {
	obj.vType = type_string
	obj.value = []byte(str)
}
func (obj *JSONObject) SetNumber(number float64) {
	obj.vType = type_number
	obj.num = number
}

func (obj *JSONObject) SetBool(b bool) {
	obj.vType = type_bool
	if b {
		obj.value = []byte("true")
	} else {
		obj.value = []byte("false")
	}
}

func (obj *JSONObject) UpdateOrInsert(key string, value *JSONObject, update bool) error {
	if obj.Type() == type_object {
		if _, ok := obj.obj[key]; ok && update {
			obj.obj[key] = value
			return nil
		} else if !update {
			obj.obj[key] = value
			return nil
		}
		return errors.New("对象没有该属性")
	}
	return errors.New("非对象不能设置属性")
}

func (obj *JSONObject) UpdateOrInsertIndex(idx int, value *JSONObject, update bool) error {
	if obj.Type() == type_array && idx >= 0 {

		if idx < len(obj.list) && update {
			obj.list[idx] = value
		} else if !update {
			obj.list = append(obj.list, value)
		}
		return nil
	}
	return errors.New("非数组不能设置属性")
}

// ====== 查 =========
// 查对象key
func (obj *JSONObject) Get(key string) *JSONObject {
	if obj.vType == type_object {
		if val, ok := obj.obj[key]; ok != true {
			return nil
		} else {
			return val
		}
	}
	return nil
}

func (obj *JSONObject) GetString() string {
	if obj.vType == type_string {
		return string(obj.value)
	}
	return ""
}

func (obj *JSONObject) GetNumber() float64 {
	if obj.vType == type_number {
		return obj.num
	}
	return 0
}

func (obj *JSONObject) GetBool() bool {
	if obj.vType == type_bool && string(obj.value) == "true" {
		return true
	}
	return false
}

func (obj *JSONObject) GetIndex(idx int) *JSONObject {
	if obj.vType == type_array && idx < len(obj.list) {
		return obj.list[idx]
	}
	return nil
}

// 获取值
func (obj *JSONObject) Value() interface{} {
	switch obj.vType {
	case type_null:
		return nil
	case type_bool:
		if string(obj.value) == "true" {
			return true
		} else {
			return false
		}
	case type_number:
		return obj.num
	case type_string:
		return string(obj.value)
	case type_array:
		return obj.list
	case type_object:
		return obj.obj
	}
	return nil
}
