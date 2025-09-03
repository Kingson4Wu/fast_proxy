# fastproxy

![fast_proxy_ico.png](https://raw.githubusercontent.com/Kingson4Wu/fast_proxy/main/resource/img/logo.jpg)

[![CI/CD Pipeline](https://github.com/Kingson4Wu/fast_proxy/actions/workflows/go.yml/badge.svg)](https://github.com/Kingson4Wu/fast_proxy/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/kingson4wu/fast_proxy)](https://goreportcard.com/report/github.com/kingson4wu/fast_proxy)
[![GitHub top language](https://img.shields.io/github/languages/top/kingson4wu/fast_proxy)](https://github.com/kingson4wu/fast_proxy/search?l=go)
[![GitHub stars](https://img.shields.io/github/stars/kingson4wu/fast_proxy)](https://github.com/kingson4wu/fast_proxy/stargazers)
[![codecov](https://codecov.io/gh/kingson4wu/fast_proxy/branch/main/graph/badge.svg)](https://codecov.io/gh/kingson4wu/fast_proxy)
[![Go Reference](https://pkg.go.dev/badge/github.com/kingson4wu/fast_proxy.svg)](https://pkg.go.dev/github.com/kingson4wu/fast_proxy)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go#database)
[![LICENSE](https://img.shields.io/github/license/kingson4wu/fast_proxy.svg?style=flat-square)](https://github.com/kingson4wu/fast_proxy/blob/main/LICENSE)

[English](https://github.com/kingson4wu/fast_proxy#fast_proxy) | [简体中文](https://github.com/kingson4wu/fast_proxy/blob/main/README-CN.md)|[deepwiki](https://deepwiki.com/Kingson4Wu/fast_proxy)


fastproxy is a high-performance proxy that can encrypt and decrypt inbound and outbound traffic, authenticate, and limit flow control.

## Key Features

* **Microservice Proxy**: Can be used as a microservice proxy for a service discovery ecosystem
* **Encryption/Decryption**: Supports encryption and decryption of data
* **Compression/Decompression**: Supports data compression and decompression
* **Signature Verification**: Supports signature verification for request traffic
* **Flow Control**: Supports flow limit control for request traffic
* **Protocol**: Uses protobuf protocol as the transmission format of intermediate data

## Design Overview

![Design Overview](https://github.com/kingson4wu/fast_proxy/blob/main/resource/img/design-overview-fast-proxy.png)

## Installation

```bash
go get github.com/Kingson4Wu/fast_proxy
```

## Quick Start

### 1. Embedded Usage

See [examples](https://github.com/kingson4wu/fast_proxy/tree/main/examples) or [more comprehensive examples](https://github.com/Kingson4Wu/fast_proxy_examples)

### 2. Command Line Usage

Start fast proxy server:

```bash
cd fast_proxy
make
./center
./server
./in-proxy 
./out-proxy 
```

Start client:

```bash
./client 
```

## Documentation

See [wiki](https://github.com/kingson4wu/fast_proxy/wiki)

## Performance

Fastproxy is optimized for high performance with:

* Minimal latency overhead
* Efficient memory usage
* Support for concurrent connections
* FastHTTP integration for improved performance

## Contributing

If you are interested in contributing to fastproxy, see [CONTRIBUTING](https://github.com/kingson4wu/fast_proxy/blob/main/CONTRIBUTING.md)

## License

fastproxy is licensed under the term of the [Apache 2.0 License](https://github.com/kingson4wu/fast_proxy/blob/main/LICENSE)