package coll

import (
	"fmt"
	"github.com/thoas/go-funk"
	"log"
	"reflect"
)

func GroupBy(list interface{}, pivotFunc interface{}) interface{} {

	// input value must be a slice
	if !funk.IsCollection(list) {
		panic(fmt.Sprintf("%v must be a collection (slice or array)", list))
	}

	keys := funk.Map(list, pivotFunc)
	count := reflect.ValueOf(keys).Len()

	keysR := reflect.ValueOf(keys)
	valuesR := reflect.ValueOf(list)

	log.Printf("keysR: %#v", keysR)
	log.Printf("valuesR: %#v", valuesR)

	valueType := valuesR.Type().Elem()
	keyType := reflect.TypeOf(pivotFunc).Out(0)

	arrayType := reflect.SliceOf(valueType)
	resultType := reflect.MapOf(keyType, arrayType)
	result := reflect.MakeMapWithSize(resultType, 0)

	for i := 0; i < count; i++ {

	}

	log.Printf("keyType: %#v", keyType.Kind().String())
	log.Printf("result: %#v", result)

	return keys
	/*
		elementType := valuesR.Type().Elem()

		elemSlice := reflect.MakeSlice(reflect.SliceOf(elementType), 0, 10)

		// We begin with a map[interface][]Elem
		collectionType := reflect.MapOf(reflect.TypeOf((*interface{})(nil)), elemSlice.Type())

		// create a map from scratch
		collection := reflect.MakeMap(collectionType)

		log.Printf("ifc: %+v", collection.Interface())

		pivotfuncValue := reflect.ValueOf(pivotfunc)

		for i := 0; i < valuesR.Len(); i++ {
			elementValue := valuesR.Index(i)
			keyValue := pivotfuncValue.Call([]reflect.Value{elementValue})
		}
	*/
	/*
		for i := 0; i < value.Len(); i++ {
		instance := value.Index(i)
		var field reflect.Value

		if instance.Kind() == reflect.Ptr {
		field = instance.Elem().FieldByName(pivot)
		} else {
		field = instance.FieldByName(pivot)
		}

		collection.SetMapIndex(field, instance)
		}
	*/
	//	return collection.Interface()
}
