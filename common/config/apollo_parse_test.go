package config

import "testing"

func TestParseApolloConfig_Failure(t *testing.T) {
    ok := ParseApolloConfig[string, int]("{not-json}", func(map[string]int){}, "name")
    if ok {
        t.Fatalf("expected false for invalid json")
    }
}

