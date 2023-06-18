package main

import "github.com/404name/go-json-structure/cmd"

func main() {
	cmd.Execute()
	// 	var v json.JSONObject

	// 	json.Parse(&v,
	// 		`
	// {
	// 	"hello": "131233",
	// 	"test": [1,2,3]
	// }
	// 	`)
	// 	println(v.Get("hello").GetString())
	// 	println(int(v.Get("test").GetIndex(1).GetNumber()))
	// 	println(json.Stringify(&v))
}
