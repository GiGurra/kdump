package errh

func Unwrap[A any](a A, err error) A {
	if err != nil {
		panic(err)
	} else {
		return a
	}
}

func Ignore(err error) {
	if err != nil {
		panic(err)
	}
}
