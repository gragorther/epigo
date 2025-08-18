package reflectutil

import (
	"reflect"
	"runtime"
	"strings"
)

func GetFunctionName(i any) string {
	strs := strings.Split((runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()), ".")
	return strs[len(strs)-1]
}
