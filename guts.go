/*
Guts attempts to break down a go data type in order to give it context within
a program. This includes trying to determine an Id, Kind, Type, and finally the
data itself. It also gets applied recursively within container structures such
as structs, maps, and slices
*/
package guts

import (
	"fmt"
	"reflect"
)

// Container housing basic information about a type
type Guts struct {
	Id         uintptr
	Kind, Type string
	Data       interface{}
}

// Break down a type and create a Gut
func Gut(v interface{}) (g Guts) {
	var val, pval reflect.Value
	defer func() { recover() }()

	val = reflect.ValueOf(v)
	pval = val
	for pval.Kind() == reflect.Ptr {
		pval = pval.Elem()
	}

	switch pval.Kind() {
	case reflect.Slice:
		g.Data = GutSlice(pval)
	case reflect.Map:
		g.Data = GutMap(pval)
	case reflect.Struct:
		g.Data = GutStruct(pval)
	default:
		g.Data = v
	}

	g.Kind = pval.Kind().String()
	g.Type = val.Type().String()
	g.Id = val.Pointer()
	return g
}

// Gut each field in a struct and rebuild as a map
func GutStruct(v reflect.Value) map[string]Guts {
	gm := make(map[string]Guts)
	vt := v.Type()

	for i := 0; i < v.NumField(); i++ {
		if vt.Field(i).PkgPath == "" {
			gm[vt.Field(i).Name] = Gut(v.Field(i).Interface())
		}
	}
	return gm
}

// Gut each pair in a map and rebuild with strings as keys
func GutMap(v reflect.Value) map[string]Guts {
	gm := make(map[string]Guts)

	for _, k := range v.MapKeys() {
		gm[fmt.Sprint(k.Interface())] = Gut(v.MapIndex(k).Interface())
	}
	return gm
}

// Gut each element of a slice
func GutSlice(v reflect.Value) []Guts {
	gs := []Guts{}

	for i := 0; i < v.Len(); i++ {
		gs = append(gs, Gut(v.Index(i).Interface()))
	}
	return gs
}
