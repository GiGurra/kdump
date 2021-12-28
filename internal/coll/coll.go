package coll

import (
	"fmt"
	"github.com/thoas/go-funk"
	"reflect"
)

func HasKey(dict map[string]interface{}, key string) bool {
	if _, ok := dict[key]; ok {
		return true
	} else {
		return false
	}
}

func GroupBy(list interface{}, pivotFunc interface{}) interface{} {

	if !funk.IsCollection(list) {
		panic(fmt.Sprintf("%v must be a collection (slice or array)", list))
	}

	// cheating here by using funk.Map, to get funk validation :).
	// Could have used reflect to call the pivot function, but then
	// I would need to re-implement all nice checks in go-funk
	keys := funk.Map(list, pivotFunc)
	count := reflect.ValueOf(keys).Len()

	keysR := reflect.ValueOf(keys)
	valuesR := reflect.ValueOf(list)

	valueType := valuesR.Type().Elem()
	keyType := reflect.TypeOf(pivotFunc).Out(0)

	sliceType := reflect.SliceOf(valueType)
	resultType := reflect.MapOf(keyType, sliceType)
	resultR := reflect.MakeMapWithSize(resultType, 0)

	initGroupCap := 1 + count/3
	if initGroupCap > 10 {
		initGroupCap = 10
	}

	for i := 0; i < count; i++ {

		keyR := keysR.Index(i)
		valueR := valuesR.Index(i)

		groupR := resultR.MapIndex(keyR)

		if groupR == (reflect.Value{}) {
			groupR = reflect.MakeSlice(sliceType, 0, initGroupCap)
		}

		newGroup := reflect.Append(groupR, valueR)
		resultR.SetMapIndex(keyR, newGroup)
	}

	return resultR.Interface()
}
