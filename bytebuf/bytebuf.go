package bytebuf

import (
	"math/bits"
	"sync"
)

var defaultPool *BuffersPool

func init() {
	defaultPool = NewBuffersPool()
}

func Make(size int) []byte {
	return defaultPool.Get(size)
}

func Put(b []byte) {
	defaultPool.Put(b)
}

const poolsCount = 32

type BuffersPool struct {
	pools [poolsCount]sync.Pool
}

func NewBuffersPool() *BuffersPool {
	bp := new(BuffersPool)
	for i := range &bp.pools {
		size := nSize(i)
		bp.pools[i] = sync.Pool{New: func() interface{} {
			return make([]byte, size)
		}}
	}
	return bp
}

func nSize(n int) int {
	return int(1 << uint(n))
}

func lenN(l int) int {
	if l == 0 {
		return 0
	}
	return bits.Len(uint(l) - 1)
}

func (bp *BuffersPool) Get(size int) []byte {
	n := lenN(size)
	if n > poolsCount-1 {
		return make([]byte, size)
	}
	b := bp.pools[n].Get().([]byte)
	return b[:size]
}

func (bp *BuffersPool) Put(b []byte) {
	n := lenN(cap(b))
	if n > poolsCount-1 {
		return
	}
	if cap(b) < nSize(n) {
		if n == 0 {
			return
		}
		n--
	}
	bp.pools[n].Put(b)
}
