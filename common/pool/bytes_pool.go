package pool

import (
	"bytes"
	"sync"
)

var (
	dataChunkSizeClasses = []int{
		1 << 10,
		2 << 10,
		4 << 10,
		8 << 10,
		16 << 10,
		32 << 10,
		64 << 10,
	}
	dataChunkPools = [...]sync.Pool{
		{New: func() interface{} {
			b := make([]byte, 1<<10)
			return &b
		}},
		{New: func() interface{} {
			b := make([]byte, 2<<10)
			return &b
		}},
		{New: func() interface{} {
			b := make([]byte, 4<<10)
			return &b
		}},
		{New: func() interface{} {
			b := make([]byte, 8<<10)
			return &b
		}},
		{New: func() interface{} {
			b := make([]byte, 16<<10)
			return &b
		}},
		{New: func() interface{} {
			b := make([]byte, 32<<10)
			return &b
		}},
		{New: func() interface{} {
			b := make([]byte, 64<<10)
			return &b
		}},
	}
)

func GetDataBufferChunk(size int64) *[]byte {
	i := 0
	for ; i < len(dataChunkSizeClasses); i++ {
		if size <= int64(dataChunkSizeClasses[i]) {
			return dataChunkPools[i].Get().(*[]byte)
		}
	}
	return nil
}
func PutDataBufferChunk(p *[]byte) {

	if p == nil {
		return
	}

	for i, n := range dataChunkSizeClasses {
		if len(*p) == n {
			dataChunkPools[i].Put(p)
			return
		}
	}
}

//---
var bufPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

func Get() *bytes.Buffer {
	b, _ := bufPool.Get().(*bytes.Buffer)
	b.Reset()
	return b
}

func Put(b *bytes.Buffer) {
	if cap(b.Bytes()) > 64<<10 {
		return
	}
	bufPool.Put(b)
}
