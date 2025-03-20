// Copyright 2025 Thom Scheeve. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rand

import "math/rand"

// The interval constants represent the different types of intervals that can be used
// when generating random numbers. Definitions are according to ISO 80000-2:2019(en).
type interval uint8

const (
	closed    interval = iota // Closed interval from a included to b included; [a,b]  (2-7.2)
	leftOpen                  // Left half-open interval from a excluded to b included; (a,b]  (2-7.3)
	rightOpen                 // Right half-open interval from a included to b excluded; [a,b)  (2-7.4)
	open                      // Open interval from a excluded to b excluded; (a,b)  (2-7.5)
)

func intmn(m, n int, interval interval) int {
	switch interval {
	case closed:
		return m + rand.Intn(n-m+1)
	case leftOpen:
		return m + rand.Intn(n-m) + 1
	case rightOpen:
		return m + rand.Intn(n-m)
	case open:
		return m + rand.Intn(n-m-1) + 1
	default:
		panic("invalid interval")
	}
}

// Intmn returns, as an int, a pseudo-random number in the half-open interval [m,n)
// from the default [math/rand.Source].
// It panics if n <= m.
func Intmn(m, n int) int {
	if n <= m {
		panic("invalid arguments to Intmn")
	}
	return intmn(m, n, rightOpen)
}

// IntmnClosed returns, as an int, a pseudo-random number in the closed interval [m,n]
// from the default [math/rand.Source].
// It panics if n < m.
func IntmnClosed(m, n int) int {
	if n < m {
		panic("invalid arguments to IntmnClosed")
	}
	return intmn(m, n, closed)
}

// IntmnLeftOpen returns, as an int, a pseudo-random number in the half-open interval (m,n]
// from the default [math/rand.Source].
// It panics if n <= m.
func IntmnLeftOpen(m, n int) int {
	if n <= m {
		panic("invalid arguments to IntmnLeftOpen")
	}
	return intmn(m, n, leftOpen)
}

// Function alias for [Intmn].
func IntmnRightOpen(m, n int) int {
	if n <= m {
		panic("invalid arguments to IntmnRightOpen")
	}
	return Intmn(m, n)
}

// IntmnOpen returns, as an int, a pseudo-random number in the open interval (m,n)
// from the default [math/rand.Source].
// It panics if n <= m.
func IntmnOpen(m, n int) int {
	if n <= m {
		panic("invalid arguments to IntmnOpen")
	}
	return intmn(m, n, open)
}
