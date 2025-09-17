package config

import "testing"

func TestOnNewestChange_NoOp(t *testing.T) {
    var c customChangeListener
    c.OnNewestChange(nil)
}

