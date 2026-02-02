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

> [qmd](https://github.com/tobi/qmd) 的 Go 语言重写版 - 用 AI 查询你的 Markdown 文档

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://go.dev)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

## 概述

gQMD 是一个单文件 MCP (Model Context Protocol) 服务器，可以索引你的 Markdown 文档并为 Claude Code 等 AI 助手提供语义搜索能力。

**核心优势:**
- 单一可执行文件，无外部依赖 (纯 Go，无 CGO)
- 跨平台支持 (macOS, Linux, Windows)
- FTS5 全文搜索，BM25 排序算法
- 基于 Ollama 嵌入的向量语义搜索
- 大幅减少分析本地文件时的 token 消耗

## 功能特性

| 功能 | 描述 |
|------|------|
| **MCP 服务器** | stdio 传输，集成 Claude Code |
| **FTS5 搜索** | SQLite 全文搜索，BM25 排序 |
| **向量搜索** | 基于 Ollama 嵌入的语义搜索 |
| **多集合管理** | 将文档组织到不同集合 |
| **CLI 工具** | 从终端管理集合和搜索 |

## 安装

### 从源码编译

```bash
# 克隆仓库
git clone https://github.com/NOTAschool/gqmd.git
cd gqmd

# 编译
make build

# 或直接使用 Go
go build -o gqmd ./cmd/gqmd
```

### 预编译二进制

从 [Releases](https://github.com/NOTAschool/gqmd/releases) 下载。

## 快速开始

### 1. 添加文档集合

```bash
# 添加 markdown 文件集合
./gqmd add docs ~/Documents/notes

# 列出集合
./gqmd list
```

### 2. 扫描和索引

```bash
# 扫描所有集合
./gqmd scan

# 搜索文档
./gqmd search "golang 教程"
```

### 3. 配合 Claude Code 使用

在 Claude Code MCP 配置中添加:

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

## MCP 工具

| 工具 | 描述 |
|------|------|
| `status` | 显示索引状态和健康信息 |
| `search` | FTS5 全文搜索 |
| `get` | 按 collection/path 获取文档 |
| `multi_get` | 批量获取多个文档 |
| `vector_search` | 语义向量搜索 (需要 Ollama) |

## CLI 命令

```bash
gqmd add <name> <path>    # 添加集合
gqmd list                 # 列出集合
gqmd remove <name>        # 删除集合
gqmd scan                 # 扫描并索引文档
gqmd search <query>       # 搜索文档
gqmd mcp                  # 启动 MCP 服务器
```

## 向量搜索配置

向量搜索需要本地运行 [Ollama](https://ollama.ai):

```bash
# 安装 Ollama 并拉取嵌入模型
ollama pull nomic-embed-text

# 启动 Ollama 服务
ollama serve

# 使用向量搜索
./gqmd mcp
# 然后通过 MCP 使用 vector_search 工具
```

## 编译

### 环境要求

- Go 1.23+
- Make (可选)

### 编译命令

```bash
# 编译当前平台
make build

# 编译所有平台
make build-all

# 运行测试
make test

# 清理构建产物
make clean
```

### 交叉编译

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

## 项目结构

```
gqmd/
├── cmd/gqmd/          # 主入口
├── internal/
│   ├── cli/           # CLI 命令 (Cobra)
│   ├── mcp/           # MCP 服务器 (stdio)
│   ├── store/         # SQLite 存储和搜索
│   └── embed/         # Ollama 嵌入客户端
└── docs/              # 文档
```

## 技术细节

- **数据库**: SQLite WASM 绑定 (ncruces/go-sqlite3)
- **搜索**: FTS5 全文搜索，BM25 排序算法
- **向量存储**: BLOB 格式 (float32 小端序)
- **相似度**: 纯 Go 实现余弦相似度
- **MCP 协议**: mark3labs/mcp-go 库

## 许可证

MIT License - 详见 [LICENSE](LICENSE)。

## 致谢

- 原版 [qmd](https://github.com/tobi/qmd) by Tobias Lutke
- [ncruces/go-sqlite3](https://github.com/ncruces/go-sqlite3) - 纯 Go SQLite
- [mark3labs/mcp-go](https://github.com/mark3labs/mcp-go) - MCP 协议
- [spf13/cobra](https://github.com/spf13/cobra) - CLI 框架
