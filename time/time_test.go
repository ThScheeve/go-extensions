// Copyright 2025 Thom Scheeve. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package time

import (
	"testing"
)

func TestDatePartString(t *testing.T) {
	if got, want := Month.String(), "month"; got != want {
		t.Errorf("month part of date = %q; want %q", got, want)
	}
	if got, want := DatePart(-1).String(), "%!DatePart(-1)"; got != want {
		t.Errorf("not a part of date = %q; want %q", got, want)
	}
}
