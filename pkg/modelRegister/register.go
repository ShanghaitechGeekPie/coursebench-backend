package modelRegister

var registeredTypes = make([]interface{}, 0)

func Register(t interface{}) {
	registeredTypes = append(registeredTypes, t)
}

func GetRegisteredTypes() []interface{} {
	return registeredTypes
}
