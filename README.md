# gQMD

```
        _____ ____  __  __ _____
   __ _|  _  |    \/  ||  _  \
  / _` | | | | |\/| | | | | |
 | (_| | |_| | |  | | | |_| |
  \__, |\___/|_|  |_|_|____/
   __/ |    (\__/)
  |___/     (o.o)  <- Gopher
             (> <)
```

> A Go rewrite of [qmd](https://github.com/tobi/qmd) - Query your Markdown documents with AI

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://go.dev)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

## Overview

gQMD is a single-binary MCP (Model Context Protocol) server that indexes your Markdown documents and provides semantic search capabilities for AI assistants like Claude Code.

**Key Benefits:**
- Single binary, no dependencies (pure Go, no CGO)
- Cross-platform (macOS, Linux, Windows)
- FTS5 full-text search with BM25 ranking
- Vector semantic search via Ollama embeddings
- Dramatically reduces token consumption when analyzing local files

## Features

| Feature | Description |
|---------|-------------|
| **MCP Server** | stdio transport for Claude Code integration |
| **FTS5 Search** | SQLite full-text search with BM25 ranking |
| **Vector Search** | Semantic search using Ollama embeddings |
| **Multi-Collection** | Organize documents into collections |
| **CLI Tools** | Manage collections and search from terminal |

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/NOTAschool/gqmd.git
cd gqmd

# Build
make build

# Or directly with Go
go build -o gqmd ./cmd/gqmd
```

### Pre-built Binaries

Download from [Releases](https://github.com/NOTAschool/gqmd/releases).

## Quick Start

### 1. Add a Document Collection

```bash
# Add a collection of markdown files
./gqmd add docs ~/Documents/notes

# List collections
./gqmd list
```

### 2. Scan and Index

```bash
# Scan all collections
./gqmd scan

# Search documents
./gqmd search "golang tutorial"
```

### 3. Use with Claude Code

Add to your Claude Code MCP configuration:

```json
{
  "mcpServers": {
    "gqmd": {
      "command": "/path/to/gqmd",
      "args": ["mcp"]
    }
  }
}
```

## MCP Tools

| Tool | Description |
|------|-------------|
| `status` | Show index status and health |
| `search` | FTS5 full-text search |
| `get` | Get document by collection/path |
| `multi_get` | Get multiple documents |
| `vector_search` | Semantic vector search (requires Ollama) |

## CLI Commands

```bash
gqmd add <name> <path>    # Add a collection
gqmd list                 # List collections
gqmd remove <name>        # Remove a collection
gqmd scan                 # Scan and index documents
gqmd search <query>       # Search documents
gqmd mcp                  # Start MCP server
```

## Vector Search Setup

Vector search requires [Ollama](https://ollama.ai) running locally:

```bash
# Install Ollama and pull embedding model
ollama pull nomic-embed-text

# Start Ollama service
ollama serve

# Use vector search
./gqmd mcp
# Then use vector_search tool via MCP
```

## Building

### Requirements

- Go 1.23+
- Make (optional)

### Build Commands

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Run tests
make test

# Clean build artifacts
make clean
```

### Cross-compilation

```bash
# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o gqmd-darwin-arm64 ./cmd/gqmd

# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o gqmd-darwin-amd64 ./cmd/gqmd

# Linux
GOOS=linux GOARCH=amd64 go build -o gqmd-linux-amd64 ./cmd/gqmd

# Windows
GOOS=windows GOARCH=amd64 go build -o gqmd-windows-amd64.exe ./cmd/gqmd
```

## Architecture

```
gqmd/
├── cmd/gqmd/          # Main entry point
├── internal/
│   ├── cli/           # CLI commands (Cobra)
│   ├── mcp/           # MCP server (stdio)
│   ├── store/         # SQLite storage & search
│   └── embed/         # Ollama embedding client
└── docs/              # Documentation
```

## Technical Details

- **Database**: SQLite with WASM bindings (ncruces/go-sqlite3)
- **Search**: FTS5 with BM25 ranking algorithm
- **Vector Storage**: BLOB format (float32 little-endian)
- **Similarity**: Cosine similarity in pure Go
- **MCP Protocol**: mark3labs/mcp-go library

## License

MIT License - see [LICENSE](LICENSE) for details.

## Credits

- Original [qmd](https://github.com/tobi/qmd) by Tobias Lutke
- [ncruces/go-sqlite3](https://github.com/ncruces/go-sqlite3) - Pure Go SQLite
- [mark3labs/mcp-go](https://github.com/mark3labs/mcp-go) - MCP Protocol
- [spf13/cobra](https://github.com/spf13/cobra) - CLI Framework
