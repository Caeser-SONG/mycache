package cache

import "fmt"

type ByteView struct {
	// 为了可以保存任何值的value.
	b []byte
}

func (v ByteView) Len() int {
	return len(v.b)
}
func (v ByteView) String() string {
	return string(v.b)
}
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}
func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
func main() {
	fmt.Println("vim-go")
}
