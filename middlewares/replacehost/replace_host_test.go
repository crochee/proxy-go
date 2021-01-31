// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/31

package replacehost

import (
	"strings"
	"testing"
)

func TestSplit(t *testing.T) {
	path := "/host"
	list := strings.SplitN(path, "/", 3)
	t.Log(list)
}
