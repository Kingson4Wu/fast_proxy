//go:build !test

package main

import (
	"embed"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Kingson4Wu/fast_proxy/common/logger/zap"
	"github.com/Kingson4Wu/fast_proxy/common/security"
	"github.com/Kingson4Wu/fast_proxy/common/server"
	"github.com/Kingson4Wu/fast_proxy/examples/center"
	"github.com/Kingson4Wu/fast_proxy/inproxy/inconfig"
)

//go:embed *
var SecureConfigFs embed.FS

func main() {
	configBytes, err := SecureConfigFs.ReadFile("secure_config.yaml")
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	c := inconfig.LoadYamlConfig(configBytes)

	sc := center.GetSC(func() string { return c.ServiceRpcHeaderName() })

	// Create security middleware
	securityMiddleware := security.NewSecurityMiddleware(
		[]byte("your-jwt-secret-key"), 
		"fastproxy-inproxy",
	)

	// Generate a token for testing
	token, err := securityMiddleware.GenerateServiceToken(
		"in_proxy", 
		[]string{"proxy", "read", "write"}, 
		24*time.Hour,
	)
	if err != nil {
		fmt.Printf("Failed to generate token: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated service token: %s\n", token)

	// Create a custom handler that uses security middleware
	originalHandler := func(res http.ResponseWriter, req *http.Request) {
		// Your original proxy logic here
		// For demonstration, we'll just return a simple response
		res.WriteHeader(http.StatusOK)
		res.Write([]byte("Secure proxy response"))
	}

	// Wrap the original handler with security middleware
	secureHandler := securityMiddleware.CombinedMiddleware(originalHandler)

	// Create the server with our secure handler
	p := server.NewServer(c, zap.DefaultLogger(), secureHandler)
	
	// Start the server with options
	p.Start(
		server.WithServiceCenter(sc),
		server.WithLogger(zap.DefaultLogger()),
	)
}