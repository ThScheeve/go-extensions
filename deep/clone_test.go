// Copyright 2025 Thom Scheeve. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package deep

import (
	"reflect"
	"runtime"
	"testing"
)

func ptr[T any](x T) *T {
	return &x
}

func assertEqual(t *testing.T, x, y any) {
	if !reflect.DeepEqual(x, y) {
		t.Errorf("%#v != %#v", x, y)
	}
}

func assertNotEqual(t *testing.T, x, y any) {
	if reflect.DeepEqual(x, y) {
		t.Errorf("%#v == %#v", x, y)
	}
}

func assertIs(t *testing.T, x, y any) {
	if !Is(x, y) {
		t.Errorf("%#v is not %#v", x, y)
	}
}

func assertIsNot(t *testing.T, x, y any) {
	if Is(x, y) {
		t.Errorf("unexpectedly identical %#v", x)
	}
}

var primitiveValues = [...]any{nil, false, true, 42, uint8(2 * 100), 3.14, 1i, "hello", "hello\u1234"}

var isTests = []func(t *testing.T){
	testIs_any,
	testIs_array,
	testIs_chan,
	testIs_func,
	testIs_map,
	testIs_pointer,
	testIs_slice,
}

func testIs_any(t *testing.T) {
	var x any
	var y = x // non-initialized
	assertEqual(t, y, x)
	assertIsNot(t, x, y)

	for _, v := range primitiveValues {
		x = v
		y = x // deep copy as per Go implementation
		assertEqual(t, y, x)
		assertIsNot(t, x, y)
		y = "42"
		assertNotEqual(t, y, x)
		assertIsNot(t, x, y)
		assertEqual(t, x, v)
	}
}

func testIs_array(t *testing.T) {
	var x [2]int
	var y = x // non-initialized
	assertEqual(t, y, x)
	assertIsNot(t, x, y)

	x = [2]int{1, 2}
	y = x // deep copy as per Go implementation
	assertEqual(t, y, x)
	assertIsNot(t, x, y)
	y[0] = 42
	assertNotEqual(t, y, x)
	assertIsNot(t, x, y)
	assertEqual(t, x[0], 1)

	x = [2]int{1, 2}
	y = [2]int{1, 2} // equal but different identity (objects)
	assertEqual(t, y, x)
	assertIsNot(t, x, y)
}

func testIs_chan(t *testing.T) {
	var x chan int
	var y = x // non-initialized
	assertEqual(t, y, x)
	assertIsNot(t, x, y)

	x = make(chan int)
	y = x // shallow copy
	assertEqual(t, y, x)
	assertIs(t, x, y)
	f := func(ch chan int, n int) {
		ch <- n
		close(ch)
	}
	go f(y, 42)
	assertEqual(t, y, x)
	assertIs(t, x, y)
	assertEqual(t, <-x, 42)

	x = make(chan int)
	y = make(chan int) // equal but different identity (objects)
	// assertEqual(t, y, x)
	assertIsNot(t, x, y)
}

func testIs_func(t *testing.T) {
	var x func(x int) int
	var y = x // non-initialized
	assertEqual(t, y, x)
	assertIsNot(t, x, y)

	x = func(x int) int { return x }
	y = x // shallow copy
	// assertEqual(t, y, x) // false per Go implementation [reflect.DeepEqual]
	assertEqual(t, y(42), x(42))
	assertIs(t, x, y)
	y = func(x int) int { return x + x }
	assertIsNot(t, x, y)
}

func testIs_map(t *testing.T) {
	var x map[int]int
	var y = x // non-initialized
	assertEqual(t, y, x)
	assertIsNot(t, x, y)

	x = make(map[int]int)
	y = x // shallow copy
	assertEqual(t, y, x)
	assertIs(t, x, y)
	y[0] = 42
	assertEqual(t, y, x)
	assertIs(t, x, y)
	assertEqual(t, x[0], 42)

	x = map[int]int{1: 1, 2: 2, 3: 3}
	y = map[int]int{1: 1, 2: 2, 3: 3} // equal but different identity (objects)
	assertEqual(t, y, x)
	assertIsNot(t, x, y)
}

func testIs_pointer(t *testing.T) {
	var x *int
	var y = x // non-initialized
	assertEqual(t, y, x)
	assertIsNot(t, x, y)

	x = new(int)
	y = x // shallow copy
	assertEqual(t, y, x)
	assertIs(t, x, y)
	*y = 42
	assertEqual(t, y, x)
	assertIs(t, x, y)
	assertEqual(t, *x, 42)

	x = ptr(42)
	y = ptr(42) // equal but different identity (objects)
	assertEqual(t, y, x)
	assertIsNot(t, x, y)
}

func testIs_slice(t *testing.T) {
	var x []int
	var y = x // non-initialized
	assertEqual(t, y, x)
	assertIsNot(t, x, y)

	x = make([]int, 3)
	y = x // shallow copy
	assertEqual(t, y, x)
	assertIs(t, x, y)
	y[0] = 42
	assertEqual(t, y, x)
	assertIs(t, x, y)
	assertEqual(t, x[0], 42)

	x = []int{1, 2, 3}
	y = []int{1, 2, 3} // equal but different identity (objects)
	assertEqual(t, y, x)
	assertIsNot(t, x, y)
}

func TestIs(t *testing.T) {
	for _, tfn := range isTests {
		t.Run(runtime.FuncForPC(reflect.ValueOf(tfn).Pointer()).Name(), tfn)
	}
}

func testDeepClone_any(t *testing.T) {
	var x any
	var y = DeepClone(x)
	assertEqual(t, y, x)
	assertIsNot(t, x, y)

	for _, v := range primitiveValues {
		x = v
		y = DeepClone(x)
		assertEqual(t, y, x)
		assertIsNot(t, x, y)
		y = "42"
		assertNotEqual(t, y, x)
		assertEqual(t, x, v)
	}
}

var deepCloneTests = []func(t *testing.T){
	testDeepClone_any,
}

func TestDeepClone(t *testing.T) {
	for _, tfn := range deepCloneTests {
		t.Run(runtime.FuncForPC(reflect.ValueOf(tfn).Pointer()).Name(), tfn)
	}
}
