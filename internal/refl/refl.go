package refl

import "reflect"

func GetPtrToFieldCopy(object reflect.Value, index int) reflect.Value {
	f0R := object.Field(index)
	f0 := f0R.Interface()
	ptrF0 := reflect.New(reflect.TypeOf(f0)) // of type *T
	ptrF0.Elem().Set(f0R)
	return ptrF0
}
