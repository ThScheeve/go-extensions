// Copyright 2025 Thom Scheeve. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package time_test

import (
	"fmt"
	gotime "time"

	"github.com/thscheeve/go-extensions/time"
)

func ExampleFloor() {
	t := gotime.Date(2012, 12, 7, 12, 15, 30, 918273645, gotime.UTC)
	floor := []time.DatePart{
		time.Second,
		time.Minute,
		time.Hour,
		time.Day,
		time.Month,
		time.Year,
	}

	for _, f := range floor {
		fmt.Printf("Floor(%v, %6s) = %s\n", t, f, time.Floor(t, f).Format("2006-01-02 15:04:05"))
	}
	// Output:
	// Floor(2012-12-07 12:15:30.918273645 +0000 UTC, second) = 2012-12-07 12:15:30
	// Floor(2012-12-07 12:15:30.918273645 +0000 UTC, minute) = 2012-12-07 12:15:00
	// Floor(2012-12-07 12:15:30.918273645 +0000 UTC,   hour) = 2012-12-07 12:00:00
	// Floor(2012-12-07 12:15:30.918273645 +0000 UTC,    day) = 2012-12-07 00:00:00
	// Floor(2012-12-07 12:15:30.918273645 +0000 UTC,  month) = 2012-12-01 00:00:00
	// Floor(2012-12-07 12:15:30.918273645 +0000 UTC,   year) = 2012-01-01 00:00:00
}
