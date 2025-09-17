![fast_proxy.png](https://raw.githubusercontent.com/Kingson4Wu/fast_proxy/main/resource/img/fast_proxy.png)
---
![fast_proxy_ico.png](https://raw.githubusercontent.com/Kingson4Wu/fast_proxy/main/resource/img/logo.jpg)

[![CI/CD Pipeline](https://github.com/Kingson4Wu/fast_proxy/actions/workflows/go.yml/badge.svg)](https://github.com/Kingson4Wu/fast_proxy/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/kingson4wu/fast_proxy)](https://goreportcard.com/report/github.com/kingson4wu/fast_proxy)
[![GitHub top language](https://img.shields.io/github/languages/top/kingson4wu/fast_proxy)](https://github.com/kingson4wu/fast_proxy/search?l=go)
[![GitHub stars](https://img.shields.io/github/stars/kingson4wu/fast_proxy)](https://github.com/kingson4wu/fast_proxy/stargazers)
[![codecov](https://codecov.io/gh/kingson4wu/fast_proxy/branch/main/graph/badge.svg)](https://codecov.io/gh/kingson4wu/fast_proxy)
[![Go Reference](https://pkg.go.dev/badge/github.com/kingson4wu/fast_proxy.svg)](https://pkg.go.dev/github.com/kingson4wu/fast_proxy)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go#database)
[![LICENSE](https://img.shields.io/github/license/kingson4wu/fast_proxy.svg?style=flat-square)](https://github.com/kingson4wu/fast_proxy/blob/main/LICENSE)

[English](https://github.com/kingson4wu/fast_proxy#fastproxy) | [简体中文](https://github.com/kingson4wu/fast_proxy/blob/main/README-CN.md) | [Deepwiki](https://deepwiki.com/Kingson4Wu/fast_proxy)

---

# FastProxy

FastProxy is a production-grade, high-performance service proxy designed to secure and accelerate east-west traffic in modern distributed systems. It combines wire-speed cryptography, signature verification, traffic governance, and observability into a single Go-native runtime that is easy to embed, automate, and operate.

## Why FastProxy

- Purpose-built for mesh-style microservices, serverless functions, and data pipelines where service-to-service trust boundaries are critical.
- Battle-tested modules for encryption, authentication, compression, and flow shaping with minimal latency overhead.
- Lean core written in Go with first-class support for protobuf payloads and FastHTTP, enabling extreme throughput on commodity hardware.
- Flexible deployment surface: run it embedded within your application, drop it in as a sidecar, or operate it as a central ingress/egress gateway.

## Core Capabilities

- **Secure Transport**: Symmetric and asymmetric encryption for inbound and outbound traffic with transparent key orchestration.
- **Traffic Integrity**: Signature verification, request tracing, and tamper detection across every hop.
- **Adaptive Flow Control**: Per-endpoint throttling, circuit breaking, and concurrency guards to protect downstream services.
- **Payload Optimization**: Built-in compression/decompression to reduce network footprint without sacrificing fidelity.
- **Ecosystem Friendly**: Protobuf-based transport primitives plus hooks for custom codecs, auth backends, and policy engines.
- **Operational Excellence**: Native metrics, structured logging (zap), and deep observability integrations for production monitoring.

## Architecture Overview

![Design Overview](https://github.com/kingson4wu/fast_proxy/blob/main/resource/img/design-overview-fast-proxy.png)

FastProxy is composed of modular components that can be combined to match your topology:

- **Center** orchestrates service metadata, configuration, and dynamic rulesets.
- **InProxy / OutProxy** handle ingress and egress enforcement, terminating secure channels and applying policy.
- **Server** hosts upstream business logic or routes to existing services.
- **Client SDK** offers first-class APIs for integrating FastProxy directly into Go applications.

The control/data plane split lets you evolve policies in real time without redeploying workloads.

## Getting Started

### Prerequisites

- Go 1.20 or newer
- GNU Make (for building binaries)
- Access to the `github.com/Kingson4Wu/fast_proxy` module

### Install via Go Modules

```bash
go get github.com/Kingson4Wu/fast_proxy
```

### Build Binaries

```bash
make
```

This produces the core executables in the project root (`center`, `server`, `in-proxy`, `out-proxy`, `client`).

### Quick Start Scenarios

1. **Embedded SDK**: Explore the [examples](https://github.com/kingson4wu/fast_proxy/tree/main/examples) or the extended [fast_proxy_examples](https://github.com/Kingson4Wu/fast_proxy_examples) repository to see how to embed FastProxy in your Go services.
2. **Sidecar / Gateway**: Run the compiled binaries to stand up a full proxy pipeline.

```bash
./center &
./server &
./in-proxy &
./out-proxy &
./client
```

Review the logs to confirm secure channels and policy enforcement are active.

## Configuration Highlights

FastProxy ships with composable configuration primitives (`common`, `inproxy`, `outproxy`, `resource`) for describing security posture, routing tables, and QoS constraints. Key themes include:

- Cryptographic suites and key rotation policies
- Authentication providers and signature requirements
- Per-route throttling, timeouts, and circuit breakers
- Observability sinks (Prometheus, Zap, tracing exporters)

Refer to the [wiki](https://github.com/kingson4wu/fast_proxy/wiki) for full schema documentation and operational guides.

## Performance & Hardening

- Optimized around FastHTTP with zero-copy buffers to minimize GC pressure.
- Proven support for high concurrency workloads, validated via the [`benchmark/`](benchmark) suite.
- Comprehensive unit and integration coverage (`coverage-*.out`) to guard against regressions.
- Security-first defaults aligned with the project's [Code of Conduct](CODE_OF_CONDUCT.md) and [Contributing](CONTRIBUTING.md) guidelines.

## Tooling & Ecosystem

- `Makefile`: streamline build, lint, and coverage workflows.
- `golangci-lint.sh`: curated static analysis for Go best practices.
- `escape_analysis.sh`: quick insights into heap allocations for performance tuning.
- Examples for integrating with protobuf pipelines, HTTP gateways, and enterprise auth providers.

## Roadmap & Community

FastProxy is actively maintained and evolves alongside modern distributed systems requirements. Upcoming focus areas include:

- Advanced policy authoring UX and multi-cluster federation
- Extended telemetry coverage (OpenTelemetry exporters, tracing adapters)
- Managed key lifecycle integrations with cloud KMS providers

We welcome feedback, feature suggestions, and contributions from the community.

## Contributing

Interested in shaping the roadmap? Review the [CONTRIBUTING.md](CONTRIBUTING.md) guidelines, open an issue, or submit a pull request. All participation is governed by our [Code of Conduct](CODE_OF_CONDUCT.md).

## License

FastProxy is released under the [Apache 2.0 License](LICENSE).

---

Looking for more detail? Dive into the [project wiki](https://github.com/kingson4wu/fast_proxy/wiki) or read the [Chinese documentation](README-CN.md).
