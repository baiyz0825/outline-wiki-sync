package jsonutils

import (
	"encoding/json"
)

func ToJsonStr(data interface{}) string {
	marshal, _ := json.Marshal(data)
	return string(marshal)
}
