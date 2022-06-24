package iter

import "reflect"

// MapKeys get an array of the map keys.
func MapKeys(map0 interface{}, keys interface{}) {
	mapValue := reflect.Indirect(reflect.ValueOf(map0))
	if mapValue.Kind() != reflect.Map {
		panic("can't get keys from no-map")
	}

	keysValue := reflect.ValueOf(keys)
	if reflect.ValueOf(keys).Kind() != reflect.Ptr {
		panic("using unaddressable value")
	} else {
		keysValue = keysValue.Elem()
	}

	if keysValue.Kind() != reflect.Slice {
		panic("needs a pointer to a slice")
	}

	elemType := keysValue.Type().Elem()
	isPtr := false
	if elemType.Kind() == reflect.Ptr {
		isPtr = true
		if elemType = elemType.Elem(); elemType.Kind() == reflect.Ptr {
			panic("not supported **TYPE")
		}
	}

	// define new slice elem function
	newElemFunc := func(v reflect.Value) reflect.Value {
		if isPtr {
			elem := reflect.New(v.Type())
			elem.Elem().Set(v)
			return elem
		}

		return v
	}

	// define set slice elem function
	setElemFunc := func(elem reflect.Value) {
		if isPtr {
			keysValue.Set(reflect.Append(keysValue, elem.Elem().Addr()))
		} else {
			keysValue.Set(reflect.Append(keysValue, elem))
		}
	}

	for _, kv := range mapValue.MapKeys() {
		setElemFunc(newElemFunc(kv))
	}
}
