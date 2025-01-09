# OpsDev CLI Tool

OpsDev is a command-line tool that helps you set up your DevOps development environment by installing and managing common DevOps tools.

## Features

- Automatic detection of installed tools
- Version management for supported tools
- Cross-platform support (Linux and macOS)
- Easy installation and updates
- Support for multiple tool versions

## Supported Tools

- AWS CLI
- Azure CLI
- Terraform
- Packer
- Vault
- Consul
- Go
- Python
- (more tools can be easily added)

## Installation

### From Source

1. Clone the repository:
```bash
git clone https://github.com/kelleyblackmore/opsdev.git
cd opsdev
```

2. Build and install:
```bash
make install
```

This will install the `opsdev` binary to `/usr/local/bin/`.

### Using Pre-built Binaries

Download the latest release for your platform from the releases page and extract it:

```bash
# For Linux:
tar xzf opsdev-<version>-linux-amd64.tar.gz
sudo mv opsdev /usr/local/bin/

# For macOS:
tar xzf opsdev-<version>-darwin-amd64.tar.gz
sudo mv opsdev /usr/local/bin/
```

## Usage

Simply run:
```bash
opsdev
```

The tool will:
1. Check your operating system
2. Detect already installed tools
3. Show available versions for each tool
4. Let you choose which tools to install/update
5. Handle the installation process

## Development

### Prerequisites

- Go 1.16 or later
- Make

### Building

```bash
# Build for current platform
make build

# Build for all supported platforms
make build-all

# Run tests
make test

# Run linter
make lint
```

### Project Structure

```
opsdev/
├── cmd/
│   └── opsdev/
│       └── main.go
├── internal/
│   ├── installer/
│   │   ├── installer.go
│   │   └── tools.go
│   └── utils/
│       ├── colors.go
│       └── system.go
├── go.mod
├── Makefile
└── README.md
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.