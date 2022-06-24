package goco

import "reflect"

type Stream struct {
	elements []interface{}
}

func NewStream(data interface{}) *Stream {
	var elements []interface{}
	dv := reflect.Indirect(reflect.ValueOf(data))
	switch dv.Kind() {
	case reflect.Slice, reflect.Array:
		value := reflect.ValueOf(data)
		length := value.Len()
		elements = make([]interface{}, 0, length)
		for i := 0; i < length; i++ {
			elements = append(elements, value.Index(i).Interface())
		}
	default:
		panic("not supported type")
	}

	s := &Stream{
		elements: elements,
	}

	return s
}

func (s *Stream) Filter(filterFunc func(elem interface{}) bool) *Stream {
	dst := make([]interface{}, 0, len(s.elements))
	for i := 0; i < len(s.elements); i++ {
		elem := s.elements[i]
		if filterFunc(elem) {
			dst = append(dst, elem)
		}
	}
	return s
}

func (s *Stream) Foreach() *Stream {
	return nil
}

func (s *Stream) Map(mapFunc func(elem interface{}) interface{}) *Stream {
	dst := make([]interface{}, 0, len(s.elements))
	for i := 0; i < len(s.elements); i++ {
		elem := mapFunc(s.elements[i])
		dst = append(dst, elem)
	}
	return s
}

func (s *Stream) Sort() *Stream {
	return nil
}

func (s *Stream) Distinct() *Stream {
	return nil
}

// -------------------------------------------

func (s *Stream) Sum() int {
	return len(s.elements)
}

func (s *Stream) ToSlice(slice interface{}) {
	sliceValue := reflect.ValueOf(slice)
	if reflect.ValueOf(slice).Kind() != reflect.Ptr {
		panic("using unaddressable value")
	} else {
		sliceValue = sliceValue.Elem()
	}

	if sliceValue.Kind() != reflect.Slice {
		panic("needs a pointer to a slice")
	}

	elemType := sliceValue.Type().Elem()
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
			sliceValue.Set(reflect.Append(sliceValue, elem.Elem().Addr()))
		} else {
			sliceValue.Set(reflect.Append(sliceValue, elem))
		}
	}

	for i := 0; i < len(s.elements); i++ {
		elem := reflect.ValueOf(s.elements[i])
		setElemFunc(newElemFunc(elem))
	}

}

//func (s *Stream) Collect(result interface{}) {
//
//}
