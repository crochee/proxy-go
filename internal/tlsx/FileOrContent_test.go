// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/2/28

package tlsx

import "testing"

func TestFileOrContent_Read(t *testing.T) {
	list := []struct {
		f    FileOrContent
		want string
	}{
		{
			f:    "/obs/file/suv/test.txt",
			want: "11",
		},
		{
			f:    "/obs/file/suv/test1.txt",
			want: "11",
		},
		{
			f:    "/obs/file/suvsdfgha.23",
			want: "/obs/file/suvsdfgha.23",
		},
	}
	for _, tt := range list {
		if got, err := tt.f.Read(); err != nil {
			t.Error(err)
		} else {
			if string(got) != tt.want {
				t.Fatalf("got %s want %s", got, tt.want)
			}
		}
	}
}
