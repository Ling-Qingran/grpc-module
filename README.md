# README for the `metricinterceptor` Package

## Overview
The `metricinterceptor` package provides an efficient and seamless way to intercept gRPC requests and responses, collecting various metrics such as request/response size, duration, error rates, and more. These metrics are then recorded in InfluxDB, making it an ideal tool for monitoring and observability in microservices architecture.

## Features
- **gRPC Interception**: Intercepts gRPC requests and responses.
- **Metric Collection**: Measures request/response sizes, execution duration, and error rates.
- **InfluxDB Integration**: Seamlessly sends metrics to an InfluxDB server.
- **Peer Information Extraction**: Retrieves IP address from gRPC requests.
- **User-Agent Extraction**: Extracts User-Agent from the request metadata.

## Prerequisites
Before using the `metricinterceptor` package, ensure you have the following:
- Go (version 1.13 or higher)
- Access to an InfluxDB server
- gRPC and Protocol Buffers installed in your Go environment

## Installation
Install the package using `go get`:

```bash
go get github.com/yourusername/metricinterceptor
