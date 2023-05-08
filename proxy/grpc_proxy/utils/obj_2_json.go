package utils

import "encoding/json"

func Obj2Json(s interface{}) string {
	bts, _ := json.Marshal(s)
	return string(bts)
}