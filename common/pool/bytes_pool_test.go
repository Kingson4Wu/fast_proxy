package pool

import (
    "bytes"
    "testing"
)

func TestGetPutDataBufferChunk(t *testing.T) {
    // request size within smallest class (1<<10)
    b := GetDataBufferChunk(100)
    if b == nil || len(*b) != 1<<10 {
        t.Fatalf("expected 1KB chunk, got len=%d", len(deref(b)))
    }
    PutDataBufferChunk(b)
    // request bigger than all classes should return nil
    if c := GetDataBufferChunk(int64(1<<20)); c != nil {
        t.Fatalf("expected nil for huge size")
    }
    // nil put should be no-op
    PutDataBufferChunk(nil)
}

func deref(p *[]byte) []byte { if p==nil { return nil }; return *p }

func TestBufferPoolGetPut(t *testing.T) {
    buf := Get()
    if buf == nil { t.Fatalf("nil buffer") }
    buf.WriteString("x")
    Put(buf)
    // Big buffer should be dropped
    big := bytes.NewBuffer(make([]byte, 0, 128<<10))
    Put(big)
}
