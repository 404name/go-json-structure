package service

import (
	"errors"
	"strconv"

	"github.com/404name/go-json-structure/json"
)

var StaticJsonPath = "./json.txt"
var JSON json.JSONObject

var (
	Invalid_Keys          string = "错误的取值"
	Invalid_Index_Keys    string = "数组下标指定错误"
	Invalid_Delete_Action string = "请不要删除主节点"
	Invalid_Update_Action string = "请不要直接修改主节点"
	Invalid_Json          string = "非法的json结构或者数值"
)

func init() {
	err := json.ParseFile(&JSON, StaticJsonPath)
	if err != nil {
		JSON = *json.NewDemo()
		json.Save(&JSON, StaticJsonPath)
	}
}

// 获取当前节点
func getCurJson(keys []string) (*json.JSONObject, error) {
	return getJson(keys, len(keys))
}

// 获取当前节点上一个
func getPreJson(keys []string) (*json.JSONObject, error) {
	return getJson(keys, len(keys)-1)
}

func getJson(keys []string, deep int) (*json.JSONObject, error) {
	var curJson *json.JSONObject
	curJson = &JSON
	println(curJson)
	for i, key := range keys {
		if i == deep {
			break
		}
		if key == "" {
			continue
		}
		if curJson == nil {
			return nil, errors.New(Invalid_Keys)
		} else if curJson.IsArray() {
			idx, err := strconv.Atoi(key)
			if err != nil {
				return nil, errors.New(Invalid_Index_Keys)
			} else {
				curJson = curJson.GetIndex(idx)
			}
		} else if curJson.IsObject() {
			curJson = curJson.Get(key)
		}
	}
	return curJson, nil
}

func Demo() string {
	JSON = *json.NewDemo()
	json.Save(&JSON, StaticJsonPath)
	return json.Stringify(&JSON)
}

func GetValue(keys []string) (string, error) {
	curJson, err := getCurJson(keys)
	if err != nil {
		return "", err
	}
	return json.Stringify(curJson), nil
}

func DeleteValue(keys []string) (string, error) {
	preJson, err := getPreJson(keys)
	if err != nil {
		return "", err
	}
	if len(keys) < 1 {
		return "", errors.New(Invalid_Delete_Action)
	}
	key := keys[len(keys)-1]
	// 删除应该只对于对象和数组而言
	if preJson.IsArray() {
		idx, err := strconv.Atoi(key)
		if err != nil {
			return "", errors.New(Invalid_Index_Keys)
		}
		preJson.RemoveIndex(idx)
	} else if preJson.IsObject() {
		preJson.Remove(key)
	}

	// 更新json并且返回
	json.Save(&JSON, StaticJsonPath)
	return json.Stringify(&JSON), nil
}

func UpdateOrInsertValue(keys []string, value string, update bool) (string, error) {
	preJson, err := getPreJson(keys)
	if err != nil {
		return "", err
	}
	if len(keys) < 1 {
		return "", errors.New(Invalid_Update_Action)
	}
	key := keys[len(keys)-1]
	// 修改应该只对于对象属性和数组下标而言

	var v json.JSONObject
	if res := json.Parse(&v, value); res != 0 {
		return "", errors.New(Invalid_Json)
	}

	if preJson.IsArray() {
		idx, err := strconv.Atoi(key)
		if err != nil {
			return "", errors.New(Invalid_Index_Keys)
		}
		if setError := preJson.UpdateOrInsertIndex(idx, &v, update); setError != nil {
			return "", setError
		}
	} else if preJson.IsObject() {
		if setError := preJson.UpdateOrInsert(key, &v, update); setError != nil {
			return "", setError
		}
	}

	// 更新json并且返回
	json.Save(&JSON, StaticJsonPath)
	return json.Stringify(&JSON), nil

}
