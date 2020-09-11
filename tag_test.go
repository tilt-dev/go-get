// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package get

var selectTagTestTags = []string{
	"go.r58",
	"go.r58.1",
	"go.r59",
	"go.r59.1",
	"go.r61",
	"go.r61.1",
	"go.weekly.2010-01-02",
	"go.weekly.2011-10-12",
	"go.weekly.2011-10-12.1",
	"go.weekly.2011-10-14",
	"go.weekly.2011-11-01",
	"go1",
	"go1.0.1",
	"go1.999",
	"go1.9.2",
	"go5",

	// these should be ignored:
	"release.r59",
	"release.r59.1",
	"release",
	"weekly.2011-10-12",
	"weekly.2011-10-12.1",
	"weekly",
	"foo",
	"bar",
	"go.f00",
	"go!r60",
	"go.1999-01-01",
	"go.2x",
	"go.20000000000000",
	"go.2.",
	"go.2.0",
	"go2x",
	"go20000000000000",
	"go2.",
	"go2.0",
}
