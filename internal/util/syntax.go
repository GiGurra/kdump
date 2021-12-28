package util

func OrPanic(first interface{}, err error) interface{} {
	if err != nil {
		panic(err)
	} else {
		return first
	}
}
