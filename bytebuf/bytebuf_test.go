package bytebuf

import (
	"github.com/stretchr/testify/require"
	"math/bits"
	"math/rand"
	"testing"
)

func TestLenNs(t *testing.T) {
	for l := 0; l < 1000; l++ {
		n := lenN(l)
		size := nSize(n)
		if size < l {
			t.Errorf("Not enough size for l=%d, n=%d, size=%d", l, n, size)
		}
	}
	for n := 0; n < poolsCount; n++ {
		size := nSize(n)
		lenN := lenN(size)
		if lenN != n {
			t.Errorf("For bucket %d size is %d, but for this size bucket is %d", n, size, lenN)
		}
	}
}

func Benchmark_SizeN(b *testing.B) {
	n := 0
	_ = n
	b.Run("simple 2x", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for l := 0; l < 1000000; l++ {
				if l == 0 {
					n = 0
					continue
				}
				n = bits.Len(uint(l) - 1)
			}
		}
	})
	b.Run("flatter 1.4x", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for l := 0; l < 1000000; l++ {
				if l < 4 {
					n = 0
				}
				bits := uint(bits.Len(uint(l)))
				secondBit := (uint(l) >> (bits - 2)) & 1
				n = int(2*bits+secondBit) - 5
			}
		}
	})
}
func Benchmark_LenS(b *testing.B) {
	n := 0
	_ = n
	b.Run("v1", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for l := 0; l < 2000; l++ {
				if l == 0 {
					n = 0
					continue
				}
				n = bits.Len(uint(l) - 1)
			}
		}
	})
	b.Run("v2", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for l := 0; l < 2000; l++ {
				n = bits.Len(uint(l))
			}
		}
	})
}

func Benchmark_Buffer(b *testing.B) {
	bp := NewBuffersPool()
	b.ResetTimer()
	var buf []byte
	for i := 0; i < b.N; i++ {
		size := rand.Intn(1000000)
		buf = bp.Get(size)
		bp.Put(buf)
	}
}

func TestBytesBuf(t *testing.T) {
	got := Make(1024)
	require.Len(t, got, 1024)
	require.Equal(t, 1024, cap(got))
	back := got[:5]
	Put(back)
	check := Make(1024)
	require.Len(t, check, 1024)
	require.Equal(t, 1024, cap(check))
}
