package ztools

import (
	"strings"
)

func CancelTargetStr(str, substr string) string {
	var (
		resStr string
		index  int
	)
	for {
		if index = strings.Index(str, substr); index == -1 {
			resStr += str
			break
		}
		resStr += str[:index]
		str = str[index+len(substr):]
	}
	return resStr
}
