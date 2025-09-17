package config

import (
    "testing"
    "github.com/apolloconfig/agollo/v4/storage"
    "github.com/Kingson4Wu/fast_proxy/common/logger"
)

func TestCustomChangeListener_OnChange(t *testing.T) {
    // ensure logger is non-nil to avoid panics
    log = logger.DiscardLogger
    // seed listener list with a no-op to exercise broadcast
    RegisterApolloListener(func(e *storage.ChangeEvent) {})

    c := &customChangeListener{}
    evt := &storage.ChangeEvent{Changes: map[string]*storage.ConfigChange{
        serviceConfigName:        {NewValue: `{"svc":{"encrypt.key.name":"enc","sign.key.name":"sig","encrypt.enable":true,"sign.enable":true,"compress.enable":true,"compress.algorithm":0}}`},
        signKeyConfigName:        {NewValue: `{"sig":"S"}`},
        encryptKeyConfigName:     {NewValue: `{"enc":"E"}`},
        serviceTimeoutConfigName: {NewValue: `{"svc":{"/x":{"timeout":123}}}`},
    }}
    c.OnChange(evt)
    if v := serviceConfigMap["svc"]; !v.EncryptEnable || !v.SignEnable || !v.CompressEnable {
        t.Fatalf("serviceConfigMap not updated: %+v", v)
    }
    if encryptKeyMap["enc"] != "E" || signKeyMap["sig"] != "S" {
        t.Fatalf("key maps not updated")
    }
    if serviceTimeoutConfigMap["svc"]["/x"].Timeout != 123 {
        t.Fatalf("timeout map not updated")
    }
}
