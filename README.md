![fast_proxy_ico.png](https://raw.githubusercontent.com/Kingson4Wu/fast_proxy/main/resource/img/logo.png)

[![Go Report Card](https://goreportcard.com/badge/github.com/kingson4wu/fast_proxy)&nbsp;](https://goreportcard.com/report/github.com/kingson4wu/fast_proxy)![GitHub top language](https://img.shields.io/github/languages/top/kingson4wu/fast_proxy)&nbsp;[![GitHub stars](https://img.shields.io/github/stars/kingson4wu/fast_proxy)&nbsp;](https://github.com/kingson4wu/fast_proxy/stargazers)[![codecov](https://codecov.io/gh/kingson4wu/fast_proxy/branch/main/graph/badge.svg)](https://codecov.io/gh/kingson4wu/fast_proxy) [![Go Reference](https://pkg.go.dev/badge/github.com/kingson4wu/fast_proxy.svg)](https://pkg.go.dev/github.com/kingson4wu/fast_proxy) [![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go#database) [![LICENSE](https://img.shields.io/github/license/kingson4wu/fast_proxy.svg?style=flat-square)](https://github.com/kingson4wu/fast_proxy/blob/main/LICENSE)

English| [简体中文](https://github.com/kingson4wu/fast_proxy/blob/main/README-CN.md)

fastproxy is a high-performance agent that can encrypt and decrypt inbound and outbound traffic, authenticate, and limit flow control, etc.

Key features:

* **Can be used as a microservice proxy for a service discovery ecosystem**
* **Supports encryption and decryption of data**
* **Support data compression and decompression**
* **Support signature verification for request traffic**
* **Support flow limit control for request traffic**
* **Use protobuf protocol as the transmission format of intermediate data**

## Design Overview

![](https://github.com/kingson4wu/fast_proxy/blob/main/resource/img/design-overview-fast-proxy.png)

## Quick Start

**1. embedded usage:** see [examples](https://github.com/kingson4wu/fast_proxy/tree/main/examples)

**2. command line usage:**

start fast proxy server

```shell
cd fast_proxy
make
./center
./server
./in-proxy 
./out-proxy 
```

```shell
./client 
```

## Documentation

See [wiki](https://github.com/kingson4wu/fast_proxy/wiki)

## Contributing

If you are interested in contributing to fastproxy, see [CONTRIBUTING](https://github.com/kingson4wu/fast_proxy/blob/main/CONTRIBUTING.md) 

## License

fastproxy is licensed under the term of the [Apache 2.0 License](https://github.com/kingson4wu/fast_proxy/blob/main/LICENSE)

