package json

//json树的节点
type JSONObject struct {
	vType int
	num   float64
	value []byte
	list  []*JSONObject
	obj   map[string]*JSONObject
}

// 查对象key
func (obj *JSONObject) Type() int {
	return obj.vType
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

// ====== 查 =========
// 查对象key
func (obj *JSONObject) Get(key string) *JSONObject {
	if obj.vType == type_object {
		return obj.obj[key]
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
