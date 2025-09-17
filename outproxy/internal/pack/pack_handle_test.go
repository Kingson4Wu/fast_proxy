package pack

import (
    "reflect"
    "testing"

    "github.com/Kingson4Wu/fast_proxy/common/config"
)

func TestPackHandler_Fmap_SuccessAndError(t *testing.T) {
    sc := &config.ServiceConfig{EncryptEnable: true}
    h := NewPackHandler([]byte("a"), sc)

    // success path: mapper uppercases
    m1 := h.Fmap(func(b []byte, _ *config.ServiceConfig) ([]byte, error) {
        return []byte("A"), nil
    }).(packHandlerImpl)
    if string(m1.b) != "A" {
        t.Fatalf("unexpected mapped bytes: %q", string(m1.b))
    }
    if !reflect.DeepEqual(m1.sc, sc) {
        t.Fatalf("service config should be preserved")
    }

    // error path: original state must be preserved
    m2 := m1.Fmap(func(b []byte, _ *config.ServiceConfig) ([]byte, error) {
        return nil, assertErr{}
    }).(packHandlerImpl)
    if string(m2.b) != "A" {
        t.Fatalf("on error, bytes must be unchanged, got %q", string(m2.b))
    }
    if !reflect.DeepEqual(m2.sc, sc) {
        t.Fatalf("service config should still be preserved")
    }
}

type assertErr struct{}
func (assertErr) Error() string { return "err" }

