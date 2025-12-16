---
sidebar_position: 2
---

<head>
  <meta name="description" content="Get started with FastProxy - Learn how to install and configure the production-grade Go service proxy for your distributed systems." />
  <meta name="keywords" content="fastproxy installation, golang proxy setup, service proxy tutorial, microservices security, go proxy configuration" />
  <meta property="og:title" content="Getting Started with FastProxy" />
  <meta property="og:description" content="Learn how to install, configure and deploy FastProxy for securing and accelerating service-to-service communication in your distributed systems." />
  <meta property="og:type" content="article" />
  <meta name="twitter:card" content="summary_large_image" />
</head>

# Getting Started with FastProxy

FastProxy is designed to be easy to integrate into your existing Go projects. This guide will walk you through installing, configuring, and running FastProxy.

## Prerequisites

Before you begin, ensure you have:

- Go 1.20 or newer installed
- GNU Make (for building from source)
- Git for version control
- Access to the `github.com/Kingson4Wu/fast_proxy` module

## Installation Methods

### Install as a Go Module

To use FastProxy in your Go project, simply add it as a dependency:

```bash
go get github.com/Kingson4Wu/fast_proxy
```

This downloads FastProxy and adds it to your `go.mod` file. You can then import the appropriate packages in your code such as `github.com/Kingson4Wu/fast_proxy/inproxy` and `github.com/Kingson4Wu/fast_proxy/outproxy`.

### Build from Source

For development or custom builds, clone and build the project:

```bash
# Clone the repository
git clone https://github.com/Kingson4Wu/fast_proxy.git
cd fast_proxy

# Build all components
make
```

The build process creates these executables:
- `center` - Service discovery and coordination component
- `in-proxy` - Inbound traffic proxy component
- `out-proxy` - Outbound traffic proxy component
- `server` - Example server component
- `client` - Example client component

## Quick Start

### 1. Basic Configuration

Create a configuration file `config.yaml`:

```yaml
proxy:
  forwardAddress: http://127.0.0.1:9833/inProxy

application:
  name: in_proxy
  port: 8033
  contextPath: /inProxy

rpc:
  serviceHeaderName: C_ServiceName

serviceConfig:
  song_service:
    encryptKeyName: encrypt.key.room.v2
    signKeyName: sign.key.room.v1
    encryptEnable: true
    signEnable: true
    compressEnable: true

signKeyConfig:
  sign.key.room.v1: abcd
  sign.key.room.v2: abcd

encryptKeyConfig:
  encrypt.key.room.v1: ABCDABCDABCDABCDW
  encrypt.key.room.v2: ABCDABCDABCDABCD

serviceCallTypeConfig:
  song_service:
    /token_service/api/service:
      callType: 1
      qps: 10

httpClient:
  MaxIdleConns: 5000
  MaxIdleConnsPerHost: 3000

fastHttp:
  enable: true
```

### 2. Run the Proxy

Start the proxy with your configuration:

```bash
# Run in-proxy
./in-proxy

# In another terminal, run out-proxy
./out-proxy
```

### 3. Embedded Usage

You can embed FastProxy directly in your Go application:

```go
package main

import (
    "github.com/Kingson4Wu/fast_proxy/inproxy"
    "github.com/Kingson4Wu/fast_proxy/inproxy/inconfig"
)

func main() {
    // Load configuration from your YAML file
    configBytes := /* your config YAML as bytes */
    config := inconfig.LoadYamlConfig(configBytes)

    // Start the proxy server
    inproxy.NewServer(config)
}
```

## Understanding the Components

### InProxy vs OutProxy

- **InProxy**: Handles incoming requests, validates signatures, decrypts requests, forwards to services
- **OutProxy**: Handles outgoing requests, encrypts requests, adds signatures, forwards to target services

### Service Configuration

Each service can have its own security settings:

```yaml
serviceConfig:
  service_name:
    encryptKeyName: "key_name_for_encryption"  # Name of encryption key to use
    signKeyName: "key_name_for_signing"       # Name of signing key to use
    encryptEnable: true                       # Enable encryption for this service
    signEnable: true                          # Enable signing for this service
    compressEnable: true                      # Enable compression for this service
```

### Apollo Configuration Center

For dynamic configuration updates, FastProxy supports Apollo:

```go
import "github.com/Kingson4Wu/fast_proxy/common/apollo"

// Initialize Apollo configuration
apolloConfig := &apollo.Config{
    ApolloAddr:    "http://your-apollo-server:8080",
    AppID:         "fastproxy",
    Cluster:       "default",
    NamespaceName: "application",
}

apollo.InitApolloConfig(apolloConfig)
```

## First Steps

1. **Start Simple**: Begin with a basic YAML configuration to understand the flow
2. **Enable FastHTTP**: Set `fastHttp.enable: true` for better performance
3. **Configure Services**: Add your services to the configuration with appropriate security settings
4. **Use Apollo**: For dynamic updates, configure Apollo Configuration Center
5. **Monitor**: Check logs to ensure requests are flowing properly with encryption/signing

For more detailed configuration options, see the [Configuration Guide](./configuration.mdx).