# Distributed Cloud Storage

A distributed cloud storage system leveraging blockchain technology

## Overview

This project builds a distributed cloud storage system that combines blockchain and P2P technologies, creating a system that doesn't rely on centralized servers.

### Key Features

- **Distributed File Storage**: Encrypt files and store them distributed across multiple nodes
- **Blockchain Management**: Manage file metadata and access permissions on blockchain
- **P2P Network**: High availability through peer-to-peer communication
- **Incentive System**: Reward system for storage providers
- **Redundancy**: Replication features to prevent data loss

## Technology Stack

- **Go**: Backend development language
- **libp2p**: P2P network communication
- **Ethereum/Polygon**: Smart contracts
- **IPFS**: Distributed file system
- **Gin**: REST API framework

## Quick Start

### Prerequisites

- Go 1.21+
- Docker & Docker Compose
- Make

### Installation

```bash
git clone https://github.com/nshmdayo/Distributed-cloud-storage-sample.git
cd Distributed-cloud-storage-sample
make setup
```

### Start Storage Node

```bash
make run-node
```

### Start API Server

```bash
make run-api
```

## Project Structure

```
/
├── cmd/              # Executables
├── internal/         # Internal packages
├── pkg/             # Public packages
├── contracts/       # Smart contracts
├── scripts/         # Build and deploy scripts
├── docs/           # Documentation
└── tests/          # Test files
```

## Contributing

We welcome contributions to this project. Please see [CONTRIBUTING.md](CONTRIBUTING.md) for details.

## License

This project is published under the MIT License. See the [LICENSE](LICENSE) file for details.