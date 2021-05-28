// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/5/16

package filecontent

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileOrContent_Read(t *testing.T) {
	list := []struct {
		f    FileOrContent
		want string
	}{
		{
			f:    "../../test/file_or_content.txt",
			want: "11",
		},
		{
			f:    "file_or_content.txt",
			want: "file_or_content.txt",
		},
	}
	for _, tt := range list {
		got, err := tt.f.Read()
		if assert.NoError(t, err) {
			assert.Equal(t, tt.want, string(got))
		}
	}
}
