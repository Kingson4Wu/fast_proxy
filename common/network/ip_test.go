package network

import "testing"

func TestGetIntranetIp_DoesNotPanic(t *testing.T) {
    _ = GetIntranetIp() // may be empty on CI; just ensure it runs
}

