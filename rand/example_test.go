// Copyright 2025 Thom Scheeve. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rand_test

import (
	"fmt"

	"github.com/thscheeve/go-extensions/rand"
)

func ExampleIntmn() {
	fmt.Println(rand.Intmn(1, 10))
	fmt.Println(rand.Intmn(1, 10))
	fmt.Println(rand.Intmn(1, 10))
}
