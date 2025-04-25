// Copyright 2025 Thom Scheeve. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package deep

import (
	"reflect"
	"runtime"
	"testing"
	"time"
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

func testDeepClone_array(t *testing.T) {
	var x [3]int
	var y = DeepClone(x)
	assertEqual(t, y, x)
	assertIsNot(t, x, y)

	x = [3]int{1, 2, 3}
	y = DeepClone(x)
	assertEqual(t, y, x)
	assertIsNot(t, x, y)
	y[0] = 42
	assertNotEqual(t, y, x)
	assertEqual(t, x[0], 1)
}

func testDeepClone_chan(t *testing.T) {
	var x chan int
	var y = DeepClone(x)
	assertEqual(t, y, x)
	assertIsNot(t, x, y)

	f := func(ch chan int, n int) {
		ch <- n
		close(ch)
	}

	// unbuffered
	x = make(chan int)
	go f(x, 42)
	time.Sleep(100 * time.Millisecond)
	y = DeepClone(x)
	assertNotEqual(t, y, x)
	assertIsNot(t, x, y)
	assertEqual(t, len(y), 0)
	assertEqual(t, cap(y), 0)

	// buffered
	x = make(chan int, 3)
	go f(x, 42)
	time.Sleep(100 * time.Millisecond)
	y = DeepClone(x)
	assertNotEqual(t, y, x)
	assertIsNot(t, x, y)
	assertEqual(t, len(y), 0)
	assertEqual(t, cap(y), 3)
}

func testDeepClone_interface(t *testing.T) {
	var x, y any

	t.Run("Array", func(t *testing.T) {
		x = [3]any{0, 0, 0}
		y = DeepClone(x)
		assertEqual(t, y, x)
		assertIsNot(t, x, y)
	})

	t.Run("Map", func(t *testing.T) {
		x = make(map[string]any)
		y = DeepClone(x)
		assertEqual(t, y, x)
		assertIsNot(t, x, y)
		y.(map[string]any)["foo"] = 42
		assertNotEqual(t, y, x)
		assertEqual(t, x.(map[string]any)["foo"], nil)
	})

	t.Run("Pointer", func(t *testing.T) {
		x = new(any)
		y = DeepClone(x)
		assertEqual(t, y, x)
		assertIsNot(t, x, y)
		*y.(*any) = 42
		assertNotEqual(t, y, x)
		assertEqual(t, *x.(*any), nil)
	})

	t.Run("Slice", func(t *testing.T) {
		x = make([]any, 3)
		y = DeepClone(x)
		assertEqual(t, y, x)
		assertIsNot(t, x, y)
		y.([]any)[0] = 42
		assertNotEqual(t, y, x)
		assertEqual(t, x.([]any)[0], nil)
	})
}

func testDeepClone_map(t *testing.T) {
	var x map[string]int
	var y = DeepClone(x)
	assertEqual(t, y, x)
	assertIsNot(t, x, y)

	x = map[string]int{"foo": 1, "bar": 2}
	y = DeepClone(x)
	assertEqual(t, y, x)
	assertIsNot(t, x, y)
	y["foo"] = 42
	assertNotEqual(t, y, x)
	assertEqual(t, x["foo"], 1)
}

func testDeepClone_pointer(t *testing.T) {
	var x *int
	var y = DeepClone(x)
	assertEqual(t, y, x)
	assertIsNot(t, x, y)

	x = new(int)
	y = DeepClone(x)
	assertEqual(t, y, x)
	assertIsNot(t, x, y)
	*y = 42
	assertNotEqual(t, y, x)
	assertEqual(t, *x, 0)
}

func testDeepClone_slice(t *testing.T) {
	var x []int
	var y = DeepClone(x)
	assertEqual(t, y, x)
	assertIsNot(t, x, y)

	x = make([]int, 3)
	y = DeepClone(x)
	assertEqual(t, y, x)
	assertIsNot(t, x, y)
	y[0] = 42
	assertNotEqual(t, y, x)
	assertEqual(t, x[0], 0)
}

func testDeepClone_struct(t *testing.T) {
	type S struct {
		Foo, Bar int
	}
	var x S
	var y = DeepClone(x)
	assertEqual(t, y, x)
	assertIsNot(t, x, y)

	x = S{1, 2}
	y = DeepClone(x)
	assertEqual(t, y, x)
	assertIsNot(t, x, y)
	y.Foo = 42
	assertNotEqual(t, y, x)
	assertEqual(t, x.Foo, 1)
}

func testDeepClone_struct_unexported_fields(t *testing.T) {
	type S struct {
		foo, bar int
	}
	var x S
	var y = DeepClone(x)
	assertEqual(t, y, x)
	assertIsNot(t, x, y)

	x = S{1, 2}
	y = DeepClone(x)
	assertEqual(t, y, x)
	assertIsNot(t, x, y)
	y.foo = 42
	assertNotEqual(t, y, x)
	assertEqual(t, x.foo, 1)
}

var deepCloneTests = []func(t *testing.T){
	testDeepClone_any,
	testDeepClone_array,
	testDeepClone_chan,
	testDeepClone_interface,
	testDeepClone_map,
	testDeepClone_pointer,
	testDeepClone_slice,
	testDeepClone_struct,
	testDeepClone_struct_unexported_fields,
}

func TestDeepClone(t *testing.T) {
	for _, tfn := range deepCloneTests {
		t.Run(runtime.FuncForPC(reflect.ValueOf(tfn).Pointer()).Name(), tfn)
	}
}
