package errors

func Panic(v interface{}) {
	if v, ok := v.(UserError); ok {
		panic(v.Stacktrace())
	}
	if v, ok := v.(error); ok {
		panic(v.Error())
	}
	panic(v)
}
