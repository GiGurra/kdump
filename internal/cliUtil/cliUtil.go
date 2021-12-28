package cliUtil

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"reflect"
)

func FindAllFlags(pointerToStruct interface{}) []cli.Flag {

	ptrR := reflect.ValueOf(pointerToStruct)

	if ptrR.Kind() != reflect.Ptr {
		panic(fmt.Sprintf("FindAllFlags needs pointer as input. Cannot use %#v", pointerToStruct))
	}

	objectR := ptrR.Elem()
	out := make([]cli.Flag, 0, 20)
	count := objectR.NumField()

	for i := 0; i < count; i++ {
		out = append(out, objectR.Field(i).Addr().Interface().(cli.Flag))
	}

	return out
}
