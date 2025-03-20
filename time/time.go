// Copyright 2025 Thom Scheeve. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package time

import (
	"fmt"
	"time"
)

// A DatePart represents a part of a date.
type DatePart int

const (
	Year DatePart = iota
	Month
	Day
	Hour
	Minute
	Second
)

// String returns the English name of the part of a date.
func (p DatePart) String() string {
	if Year <= p && p <= Second {
		return longDatePartNames[p]
	}
	return fmt.Sprintf("%%!DatePart(%d)", p)
}

// Floor returns the result of rounding t down to the nearest part of t.
// If p is not a valid part of a date, Floor returns t unchanged.
//
// Floor operates on the presentation form of the time; it does not
// operate on the time as an absolute duration since the zero time.
// Thus, Floor(t, Hour) will not return a time with a non-zero minute.
func Floor(t time.Time, p DatePart) time.Time {
	year, month, day := t.Date()
	hour, min, sec := t.Clock()
	switch p {
	case Year:
		month, day, hour, min, sec = 1, 1, 0, 0, 0
	case Month:
		day, hour, min, sec = 1, 0, 0, 0
	case Day:
		hour, min, sec = 0, 0, 0
	case Hour:
		min, sec = 0, 0
	case Minute:
		sec = 0
	case Second:
		// do nothing as nsec will always be 0
	default:
		return t
	}
	return time.Date(year, month, day, hour, min, sec, 0, t.Location())
}
