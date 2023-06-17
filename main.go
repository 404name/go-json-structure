package main

import json "github.com/404name/go-json-structure/json"

func main() {
	var v json.JSONObject

	json.Parse(&v,
		`
{
	"hello": "131233",
	"test": [1,2,3]
}
	`)
	println(v.Get("hello").GetString())
	println(int(v.Get("test").GetIndex(1).GetNumber()))
	println(json.Stringify(&v))
}
