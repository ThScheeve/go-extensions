// Copyright 2025 Thom Scheeve. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package deep

import (
	"log"
	"reflect"
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
