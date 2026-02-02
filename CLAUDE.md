# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目概述

gqmd 是 [qmd](https://github.com/tobi/qmd) 的 Golang 重实现，目标是生成单一可执行文件，为 Claude Code 提供 MCP (Model Context Protocol) 服务，用于本地文档搜索和分析，以节省 token 消耗。

## 核心功能模块

- **MCP Server**: stdio 传输协议，暴露文档搜索工具
- **BM25 全文搜索**: 基于 SQLite FTS5
- **向量语义搜索**: sqlite-vec + GGUF embedding 模型
- **LLM 重排序**: qwen3-reranker GGUF 模型
- **查询扩展**: 微调的 query-expansion 模型
- **文档索引管理**: Collection/Context 管理

## 技术栈选型

| 模块 | 推荐库 | 备选 |
|------|--------|------|
| MCP Server | github.com/modelcontextprotocol/go-sdk | - |
| SQLite | github.com/mattn/go-sqlite3 | ncruces/go-sqlite3 (无 CGO) |
| sqlite-vec | asg017/sqlite-vec-go-bindings | - |
| llama.cpp | dianlight/gollama.cpp (无 CGO) | go-skynet/go-llama.cpp |
| CLI | spf13/cobra | urfave/cli |

## 构建与运行

```bash
# 构建
go build -o gqmd .

# 运行测试
go test ./...

# 运行单个测试
go test -run TestName ./path/to/package

# 带覆盖率测试
go test -cover ./...
```

## 跨平台编译注意事项

- 纯 Go 代码默认静态链接
- CGO 依赖（如 go-sqlite3）需要额外工具链
- 可使用 Zig 作为 CGO 交叉编译工具链
- 优先选择无 CGO 方案（如 gollama.cpp、ncruces/go-sqlite3）以简化分发
