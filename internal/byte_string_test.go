package internal

import "testing"

func BenchmarkStringByte(b *testing.B) {
	input := []byte(`hsudishgd11111111111111115445444444444444444444444444444444444444444444444`)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		t := String(input)
		_ = Bytes(t)
	}
}

func BenchmarkCommon(b *testing.B) {
	input := []byte(`hsudishgd11111111111111115445444444444444444444444444444444444444444444444`)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		t := string(input)
		_ = []byte(t)
	}
}
