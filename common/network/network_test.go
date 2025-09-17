package network

import (
    "net"
    "testing"
)

func TestTelnet(t *testing.T) {
    // start a local listener on a free port
    ln, err := net.Listen("tcp", "127.0.0.1:0")
    if err != nil { t.Skip("cannot listen on local port") }
    defer ln.Close()
    addr := ln.Addr().(*net.TCPAddr)
    if !Telnet("127.0.0.1", addr.Port) {
        t.Fatalf("expected telnet success to local listener")
    }
}

func TestTelnet_Failure(t *testing.T) {
    // pick a free port and close the listener; then telnet should fail
    ln, err := net.Listen("tcp", "127.0.0.1:0")
    if err != nil { t.Skip("cannot listen on local port") }
    addr := ln.Addr().(*net.TCPAddr)
    _ = ln.Close()
    if Telnet("127.0.0.1", addr.Port) {
        t.Fatalf("expected telnet failure to closed port")
    }
}
