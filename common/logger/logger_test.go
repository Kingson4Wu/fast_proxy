package logger

import (
    "bytes"
    stdlog "log"
    "testing"
    "time"
)

func TestVerboseLoggerOutputs(t *testing.T) {
    var buf bytes.Buffer
    l := VerbosePrintfLogger(stdlog.New(&buf, "", 0), &buf)
    l.Debug("dbg", "k", 1)
    l.Debugf("dbg %d", 2)
    l.Info("info", "k", 3)
    l.Infof("info %d", 4)
    l.Warn("warn", "k", 5)
    l.Warnf("warn %d", 6)
    l.Error("err", "k", 7, "ts", time.Unix(0,0))
    l.Errorf("err %d", 8)
    l.Flush()
    if buf.Len() == 0 {
        t.Fatalf("expected output from verbose logger")
    }
}

func TestPrintfLoggerOnlyErrors(t *testing.T) {
    var buf bytes.Buffer
    l := PrintfLogger(stdlog.New(&buf, "", 0), &buf)
    l.Info("should not appear")
    l.Debug("should not appear")
    l.Error("err")
    if got := buf.String(); got == "" {
        t.Fatalf("expected error output")
    }
}
