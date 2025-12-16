---
sidebar_position: 1
---

<head>
  <meta name="description" content="FastProxy is a production-grade, high-performance service proxy for modern distributed systems. Secure east-west traffic with wire-speed encryption, signature verification, and traffic governance." />
  <meta name="keywords" content="fastproxy, golang proxy, service proxy, microservices, api gateway, http proxy, security, encryption, authentication" />
  <meta property="og:title" content="Introduction to FastProxy - Enterprise-Grade Service Proxy" />
  <meta property="og:description" content="Learn about FastProxy, a production-grade, high-performance service proxy designed to secure and accelerate east-west traffic in modern distributed systems." />
  <meta property="og:type" content="article" />
  <meta name="twitter:card" content="summary_large_image" />
</head>

# Introduction to FastProxy

<div align="center">
  <img src="/img/fast_proxy.png" alt="FastProxy logo" width="300"/>
</div>

<br/>

[FastProxy](https://github.com/Kingson4Wu/fast_proxy) is a **production-grade, high-performance service proxy** designed to secure and accelerate east-west traffic in modern distributed systems. It combines wire-speed cryptography, signature verification, traffic governance, and observability into a single Go-native runtime that is easy to embed, automate, and operate.

## Why FastProxy?

FastProxy addresses critical challenges in modern distributed systems:

- **Service-to-service security**: Purpose-built for mesh-style microservices, serverless functions, and data pipelines where service-to-service trust boundaries are critical.
- **Performance at scale**: Battle-tested modules for encryption, authentication, compression, and flow shaping with minimal latency overhead.
- **Go-native implementation**: Lean core written in Go with first-class support for protobuf payloads and FastHTTP, enabling extreme throughput on commodity hardware.
- **Flexible deployment**: Run it embedded within your application, drop it in as a sidecar, or operate it as a central ingress/egress gateway.

## Key Features

### High-Performance Architecture
FastProxy is built with performance in mind:
- Optimized with FastHTTP for maximum throughput
- Wire-speed encryption and decryption
- Minimal latency overhead for typical operations
- Designed to handle thousands of concurrent connections

### Security-First Design
- Symmetric and asymmetric encryption for inbound and outbound traffic
- Signature verification and tamper detection across every hop
- Per-endpoint throttling and circuit breaking
- Built-in compression to reduce network footprint

### Ecosystem Compatibility
- First-class support for protobuf payloads
- Integrates with common authentication providers
- Compatible with standard HTTP/HTTPS traffic
- Hooks for custom codecs and policy engines

## Core Capabilities

- **Secure Transport**: Symmetric and asymmetric encryption for inbound and outbound traffic with transparent key orchestration.
- **Traffic Integrity**: Signature verification, request tracing, and tamper detection across every hop.
- **Adaptive Flow Control**: Per-endpoint throttling, circuit breaking, and concurrency guards to protect downstream services.
- **Payload Optimization**: Built-in compression/decompression to reduce network footprint without sacrificing fidelity.
- **Ecosystem Friendly**: Protobuf-based transport primitives plus hooks for custom codecs, auth backends, and policy engines.
- **Operational Excellence**: Native metrics, structured logging (zap), and deep observability integrations for production monitoring.

## Architecture Overview

FastProxy is composed of modular components that can be combined to match your topology:

- **Center** orchestrates service metadata, configuration, and dynamic rulesets.
- **InProxy / OutProxy** handle ingress and egress enforcement, terminating secure channels and applying policy.
- **Server** hosts upstream business logic or routes to existing services.
- **Client SDK** offers first-class APIs for integrating FastProxy directly into Go applications.

The control/data plane split lets you evolve policies in real time without redeploying workloads.

<div align="center">
  <img src="/img/design-overview-fast-proxy.png" alt="Architecture Overview" width="800"/>
</div>

## Project Status

[![CI/CD Pipeline](https://github.com/Kingson4Wu/fast_proxy/actions/workflows/go.yml/badge.svg)](https://github.com/Kingson4Wu/fast_proxy/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/kingson4wu/fast_proxy)](https://goreportcard.com/report/github.com/kingson4wu/fast_proxy)
[![GitHub top language](https://img.shields.io/github/languages/top/kingson4wu/fast_proxy)](https://github.com/kingson4Wu/fast_proxy/search?l=go)
[![GitHub stars](https://img.shields.io/github/stars/kingson4wu/fast_proxy)](https://github.com/kingson4wu/fast_proxy/stargazers)
[![codecov](https://codecov.io/gh/kingson4wu/fast_proxy/branch/main/graph/badge.svg)](https://codecov.io/gh/kingson4wu/fast_proxy)
[![Go Reference](https://pkg.go.dev/badge/github.com/kingson4wu/fast_proxy.svg)](https://pkg.go.dev/github.com/kingson4wu/fast_proxy)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go#database)
[![LICENSE](https://img.shields.io/github/license/kingson4wu/fast_proxy.svg?style=flat-square)](https://github.com/kingson4wu/fast_proxy/blob/main/LICENSE)