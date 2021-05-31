package internal

import "testing"

func TestSlice(t *testing.T) {
	c := []int{1, 2, 3}
	t.Log(len(c), cap(c))
	c = c[0:0] // 底层数组没有改变
	t.Log(len(c), cap(c))
	c = c[:0:0] // 清空底层数组，即第三层的容量置零
	t.Log(len(c), cap(c))
}
