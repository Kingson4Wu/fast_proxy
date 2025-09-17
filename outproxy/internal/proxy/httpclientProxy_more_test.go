package proxy

import (
    "net"
    "testing"
)

type nonTimeoutErr struct{}
func (nonTimeoutErr) Error() string   { return "err" }
func (nonTimeoutErr) Timeout() bool   { return false }
func (nonTimeoutErr) Temporary() bool { return false }

func TestShouldRetryErr_NonTimeout(t *testing.T) {
    if shouldRetryErr(nonTimeoutErr{}) {
        t.Fatalf("should not retry on non-timeout non-temporary error")
    }
    // also verify that plain error without net.Error does not retry
    if shouldRetryErr(errPlain("x")) {
        t.Fatalf("should not retry on plain error")
    }
}

type errPlain string
func (e errPlain) Error() string { return string(e) }

func TestBackoffDuration_RangeAndClip(t *testing.T) {
    // check jitter range for attempt 0 based on default baseBackoff
    d0 := backoffDuration(0)
    min := baseBackoff - (baseBackoff/5) // 80%
    max := baseBackoff + (baseBackoff/5) // 120%
    if d0 < min || d0 > max {
        t.Fatalf("backoff attempt0 out of range: %v not in [%v,%v]", d0, min, max)
    }
    // very large attempt should clip near maxBackoff with jitter ≤20%
    dX := backoffDuration(10)
    if dX > maxBackoff + (maxBackoff/5) {
        t.Fatalf("backoff beyond jittered max: %v > %v", dX, maxBackoff+(maxBackoff/5))
    }
    _ = net.IPv4len // keep net import used
}
