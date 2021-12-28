package cliUtil

import (
	"github.com/gigurra/kdump/internal/refl"
	"github.com/urfave/cli/v2"
	"reflect"
)

func FindAllFlags(object interface{}) []cli.Flag {

	out := make([]cli.Flag, 0, 20)

	objectR := reflect.ValueOf(object)
	count := objectR.NumField()

	for i := 0; i < count; i++ {
		out = append(out, refl.GetPtrToFieldCopy(objectR, i).Interface().(cli.Flag))
	}

	return out
}
