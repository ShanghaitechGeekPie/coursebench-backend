package events

func Panic(v interface{}) {
	if v, ok := v.(*AttributedEvent); ok {
		panic(v)
	}
	if v, ok := v.(error); ok {
		panic(v)
	}
	panic(v)
}
