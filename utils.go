package restle

import (
    "reflect"
)

func getTypeName(f interface{}) string {
    v := reflect.Indirect(reflect.ValueOf(f))
    name := v.Type().Name()
    return name
}

func createInstanceByType(t reflect.Type) interface{} {
    return reflect.New(t).Interface()
}

func createSliceByType(t reflect.Type) interface{} {
    return reflect.New(reflect.Type(reflect.SliceOf(t))).Interface()
}
