// Copyright 2025 Thom Scheeve. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package deep

import (
	"log"
	"reflect"
	"unsafe"
)

// Checks for identity
func Is(x, y any) bool {
	if x == nil || y == nil {
		log.Default().Printf("is: nil")
		return false
	}
	v1 := reflect.ValueOf(x)
	v2 := reflect.ValueOf(y)
	if v1.Type() != v2.Type() {
		log.Default().Printf("is: %v != %v", v1.Type(), v2.Type())
		return false
	}
	switch v1.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Pointer, reflect.Slice, reflect.UnsafePointer:
		if v1.IsNil() || v2.IsNil() {
			// Can't check identity of a nil
			log.Default().Printf("is: nil %v", v1.Kind())
			return false
		}
		log.Default().Printf("is: %v == %v", v1.UnsafePointer(), v2.UnsafePointer())
		return v1.UnsafePointer() == v2.UnsafePointer()
	default:
		// by default deep copy, thus different identities
		log.Default().Printf("is: default")
		return false
	}
}

func DeepClone[T any](x T) T {
	v := reflect.ValueOf(x)
	if !v.IsValid() {
		log.Default().Printf("deepClone: nil")
		return x
	}
	return deepValueClone(v).Interface().(T)
}

func deepValueClone(v reflect.Value) reflect.Value {
	log.Default().Printf("deepValueClone: %v %v", v.Kind(), v.Type())
	switch v.Kind() {
	case reflect.Array:
		val := reflect.New(v.Type()).Elem()
		for i := 0; i < v.Len(); i++ {
			val.Index(i).Set(deepValueClone(v.Index(i)))
		}
		return val
	case reflect.Interface:
		if v.IsNil() {
			log.Default().Printf("deepValueClone: nil %v", v.Kind())
			return v
		}
		val := reflect.New(v.Elem().Type()).Elem()
		val.Set(deepValueClone(v.Elem()))
		return val
	case reflect.Map:
		if v.IsNil() {
			log.Default().Printf("deepValueClone: nil %v", v.Kind())
			return v
		}
		val := reflect.MakeMapWithSize(v.Type(), v.Len())
		for _, k := range v.MapKeys() {
			val.SetMapIndex(deepValueClone(k), deepValueClone(v.MapIndex(k)))
		}
		return val
	case reflect.Pointer:
		if v.IsNil() {
			log.Default().Printf("deepValueClone: nil %v", v.Kind())
			return v
		}
		val := reflect.New(v.Type().Elem())
		val.Elem().Set(deepValueClone(v.Elem()))
		return val
	case reflect.Slice:
		if v.IsNil() {
			log.Default().Printf("deepValueClone: nil %v", v.Kind())
			return v
		}
		val := reflect.MakeSlice(v.Type(), v.Len(), v.Cap())
		for i := 0; i < v.Len(); i++ {
			val.Index(i).Set(deepValueClone(v.Index(i)))
		}
		return val
	case reflect.Struct:
		val := reflect.New(v.Type()).Elem()
		val.Set(v)
		for i := 0; i < v.NumField(); i++ {
			f := val.Field(i)
			f = reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).Elem()
			f.Set(deepValueClone(f))
		}
		return val
	default:
		log.Default().Printf("deepValueClone: default")
		return v
	}
}
