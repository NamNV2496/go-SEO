package utils

import (
	"encoding/json"
	"reflect"
)

func Copy(dst any, src any) {
	if src == nil {
		return
	}
	jsonData, _ := json.Marshal(src)
	json.Unmarshal(jsonData, dst)
}

func CopySlice(dst any, src any) {
	srcVal := reflect.ValueOf(src)
	if srcVal.Kind() != reflect.Slice || srcVal.Len() == 0 {
		return
	}
	for i := 0; i < srcVal.Len(); i++ {
		var elem any
		Copy(elem, srcVal.Index(i).Interface())
		dstVal := reflect.ValueOf(dst)
		dstVal.Set(reflect.Append(dstVal, reflect.ValueOf(elem)))
	}
}
