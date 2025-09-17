package sign

import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
    "net/http"
    "testing"
    "time"

    "github.com/Kingson4Wu/fast_proxy/common/config"
    "github.com/Kingson4Wu/fast_proxy/common/logger/zap"
    "github.com/Kingson4Wu/fast_proxy/common/server"
)

type cfgForSign struct{ key string }

func (c *cfgForSign) ServerName() string                    { return "sign" }
func (c *cfgForSign) ServerPort() int                        { return 0 }
func (c *cfgForSign) ServiceRpcHeaderName() string           { return "C_ServiceName" }
func (c *cfgForSign) GetServiceConfig(string) *config.ServiceConfig { return &config.ServiceConfig{SignEnable: true, SignKeyName: "k"} }
func (c *cfgForSign) GetSignKey(sc *config.ServiceConfig) string    { return c.key }
func (c *cfgForSign) GetSignKeyByName(string) string         { return c.key }
func (c *cfgForSign) GetEncryptKey(*config.ServiceConfig) string    { return "" }
func (c *cfgForSign) GetEncryptKeyByName(string) string      { return "" }
func (c *cfgForSign) GetTimeoutConfigByName(string, string) int { return 0 }
func (c *cfgForSign) HttpClientMaxIdleConns() int            { return 10 }
func (c *cfgForSign) HttpClientMaxIdleConnsPerHost() int     { return 10 }
func (c *cfgForSign) FastHttpEnable() bool                   { return true }

func startSignServerOnce(key string) {
    cfg := &cfgForSign{key: key}
    go func() {
        p := server.NewServer(cfg, zap.DefaultLogger(), func(http.ResponseWriter, *http.Request) {})
        p.Start()
    }()
    // brief wait for global to be set
    time.Sleep(20 * time.Millisecond)
}

func TestGenerateSign(t *testing.T) {
    got, err := GenerateSign([]byte("hello"), "secret")
    if err != nil { t.Fatalf("GenerateSign err: %v", err) }

    m := hmac.New(sha256.New, []byte("secret"))
    m.Write([]byte("hello"))
    want := hex.EncodeToString(m.Sum(nil))
    if got != want { t.Fatalf("mismatch: got %s want %s", got, want) }
}

func TestGenerateBodySignWithName(t *testing.T) {
    startSignServerOnce("sek")
    got, err := GenerateBodySignWithName([]byte("body"), "k")
    if err != nil { t.Fatalf("err: %v", err) }
    m := hmac.New(sha256.New, []byte("sek"))
    m.Write([]byte("body"))
    want := hex.EncodeToString(m.Sum(nil))
    if got != want { t.Fatalf("mismatch: got %s want %s", got, want) }
}

func TestGenerateBodySign_ErrorWhenNoKey(t *testing.T) {
    startSignServerOnce("")
    if _, err := GenerateBodySign([]byte("b"), &config.ServiceConfig{}); err == nil {
        t.Fatalf("expected error when no sign key")
    }
}
