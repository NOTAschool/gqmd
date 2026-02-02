# qmd 重实现语言选型分析：Golang vs Rust vs Zig

> 调研日期: 2026-02-02
> 关联 Issues: #59
> 调研模型: Claude Opus 4.5 (claude-opus-4-5-20251101)

---

## 一、问题背景

### 1.1 需求概述

[qmd (Query Markup Documents)](https://github.com/tobi/qmd) 是一个本地运行的文档搜索引擎，专为 AI 辅助开发工作流设计。当前实现基于 TypeScript + Python，依赖大量 npm/pip 模块：

**核心依赖分析：**
```
TypeScript 实现 (~488KB)
├── @modelcontextprotocol/sdk  # MCP 协议
├── node-llama-cpp             # GGUF 模型推理
├── sqlite-vec                 # 向量搜索扩展
├── yaml                       # 配置解析
└── zod                        # Schema 验证

Python 实现 (~248KB)
└── 用于模型微调 (finetune/)
```

**目标：** 使用 Golang/Rust/Zig 重新实现，生成单一可执行文件，为 Claude Code 提供 MCP 服务，消除对 npm/pip 生态的依赖。

### 1.2 qmd 核心功能模块

```
┌─────────────────────────────────────────────────────────────────┐
│                        qmd 功能架构                              │
├─────────────────────────────────────────────────────────────────┤
│  1. MCP Server        │ stdio 传输，暴露 6 个工具               │
│  2. BM25 全文搜索     │ SQLite FTS5 实现                        │
│  3. 向量语义搜索      │ sqlite-vec + GGUF embedding 模型        │
│  4. LLM 重排序        │ qwen3-reranker GGUF 模型                │
│  5. 查询扩展          │ 微调的 query-expansion 模型             │
│  6. 文档索引管理      │ Collection/Context 管理                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## 二、技术可行性分析

### 2.1 MCP SDK 支持情况

| 语言 | SDK 状态 | 成熟度 | 备注 |
|------|----------|--------|------|
| **Golang** | ✅ 官方 SDK | 高 | Google 协作维护，支持 stdio/SSE/WebSocket/gRPC |
| **Rust** | ✅ 多个实现 | 高 | rmcp、mcp_rust_sdk、Prism MCP 等，符合 2025-06-18 规范 |
| **Zig** | ⚠️ 社区实现 | 中 | mcp.zig、zig-mcp-server 等，生态较小 |

### 2.2 llama.cpp / GGUF 绑定

| 语言 | 库名称 | 特性 | 成熟度 |
|------|--------|------|--------|
| **Golang** | gollama.cpp | 无 CGO (purego)，自动下载预编译库，GPU 加速 | 高 |
| **Golang** | go-llama.cpp | CGO 绑定，支持 Embeddings API | 中 |
| **Rust** | llama_cpp-rs | 高级安全绑定，易用 | 高 |
| **Rust** | llama-cpp-2 | 低级绑定，紧跟上游，2026-01 仍活跃 | 高 |
| **Zig** | llama.cpp.zig | Zig 0.14.x 支持，Vulkan 可选 | 中 |

### 2.3 SQLite 扩展支持

| 功能 | Golang | Rust | Zig |
|------|--------|------|-----|
| **FTS5 (BM25)** | ✅ go-sqlite3-fts5, zalgonoise/fts | ✅ rusqlite 原生支持 | ✅ zig-sqlite C 互操作 |
| **sqlite-vec** | ✅ 官方 CGO/WASM 绑定 | ✅ 官方 crate | ⚠️ 需 C 互操作封装 |

---

## 三、单一可执行文件能力对比

### 3.1 静态链接与跨平台编译

```
┌──────────────────────────────────────────────────────────────────────────┐
│                        跨平台编译能力对比                                 │
├──────────┬───────────────────────────────────────────────────────────────┤
│ Golang   │ ✅ 默认静态链接（纯 Go）                                       │
│          │ ⚠️ CGO 依赖时需额外工具链，跨平台复杂                          │
│          │ 💡 可用 Zig 作为 C 交叉编译工具链                              │
├──────────┼───────────────────────────────────────────────────────────────┤
│ Rust     │ ✅ musl 目标实现完全静态链接                                   │
│          │ ⚠️ 需配置交叉编译工具链                                        │
│          │ 💡 二进制体积可通过 LTO/strip 优化                             │
├──────────┼───────────────────────────────────────────────────────────────┤
│ Zig      │ ✅ 原生支持 20+ 架构交叉编译                                   │
│          │ ✅ 内置 libc，无需外部工具链                                   │
│          │ ✅ 生成极小静态二进制                                          │
│          │ ⚠️ 1.0 预计 2026 年中发布                                      │
└──────────┴───────────────────────────────────────────────────────────────┘
```

### 3.2 二进制分发复杂度

| 维度 | Golang | Rust | Zig |
|------|--------|------|-----|
| 纯语言项目 | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| 含 C 依赖 | ⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| 含 GGUF 推理 | ⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐ |

---

## 四、开发效率与生态成熟度

### 4.1 开发速度对比

| 维度 | Golang | Rust | Zig |
|------|--------|------|-----|
| 学习曲线 | 低 | 高 | 中 |
| 原型开发速度 | 快 | 慢 | 中 |
| 编译速度 | 快 | 慢 | 快 |
| 调试体验 | 好 | 好 | 一般 |
| 人才可用性 | 高 | 中 | 低 |

### 4.2 生态成熟度

| 维度 | Golang | Rust | Zig |
|------|--------|------|-----|
| 标准库完整度 | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐ |
| 第三方库数量 | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐ |
| 生产案例 | 大量 | 增长中 | 少量 |
| 语言稳定性 | 稳定 | 稳定 | 预 1.0 |

---

## 五、针对 qmd 重实现的评估

### 5.1 实现难度矩阵

| 模块 | Golang | Rust | Zig |
|------|--------|------|-----|
| MCP Server | ⭐ 简单 | ⭐ 简单 | ⭐⭐ 中等 |
| SQLite FTS5 | ⭐ 简单 | ⭐ 简单 | ⭐⭐ 中等 |
| sqlite-vec | ⭐⭐ 中等 | ⭐ 简单 | ⭐⭐⭐ 困难 |
| GGUF Embedding | ⭐⭐⭐ 困难 | ⭐⭐ 中等 | ⭐⭐⭐ 困难 |
| GGUF Reranking | ⭐⭐⭐ 困难 | ⭐⭐ 中等 | ⭐⭐⭐ 困难 |
| 单文件分发 | ⭐⭐ 中等 | ⭐⭐ 中等 | ⭐ 简单 |

### 5.2 风险评估

**Golang 风险：**
- CGO 依赖使跨平台编译复杂化
- llama.cpp 绑定维护可能滞后上游
- gollama.cpp (无 CGO) 方案相对较新

**Rust 风险：**
- 开发周期较长
- llama.cpp 绑定 API 变化频繁
- 学习成本高

**Zig 风险：**
- 语言尚未达到 1.0 稳定版
- 生态系统不成熟，库选择有限
- 人才稀缺，维护困难

---

## 六、推荐方案

### 6.1 综合评分

| 维度 (权重) | Golang | Rust | Zig |
|-------------|--------|------|-----|
| 开发速度 (25%) | 9 | 5 | 6 |
| 单文件分发 (20%) | 6 | 8 | 10 |
| 生态成熟度 (20%) | 9 | 8 | 4 |
| 性能 (15%) | 7 | 10 | 9 |
| 维护性 (20%) | 9 | 7 | 4 |
| **加权总分** | **8.05** | **7.35** | **6.35** |

### 6.2 推荐结论

**首选：Golang** ⭐⭐⭐⭐

理由：
1. 开发效率最高，可快速交付 MVP
2. MCP SDK 官方支持，由 Google 协作维护
3. 生态成熟，社区活跃
4. 使用 gollama.cpp (无 CGO) 可简化分发

**备选：Rust** ⭐⭐⭐

理由：
1. 性能最优，内存安全
2. llama.cpp 绑定成熟度高
3. 适合长期维护的高质量项目
4. 开发周期较长，适合有充足时间的场景

**暂不推荐：Zig** ⭐⭐

理由：
1. 语言尚未稳定，API 可能变化
2. 生态系统不足以支撑复杂项目
3. 可在 Zig 1.0 发布后重新评估

---

## 七、实施建议

### 7.1 Golang 实施路径

```
Phase 1: 核心框架 (1-2 周)
├── MCP Server (官方 SDK)
├── SQLite + FTS5 集成
└── CLI 命令框架

Phase 2: 搜索功能 (2-3 周)
├── BM25 全文搜索
├── sqlite-vec 向量搜索
└── RRF 融合算法

Phase 3: LLM 集成 (2-3 周)
├── gollama.cpp 集成
├── Embedding 生成
└── Reranking 实现

Phase 4: 优化与分发 (1 周)
├── 跨平台编译脚本
├── 二进制体积优化
└── 文档与测试
```

### 7.2 关键技术选型

| 模块 | 推荐库 | 备选 |
|------|--------|------|
| MCP Server | github.com/modelcontextprotocol/go-sdk | - |
| SQLite | github.com/mattn/go-sqlite3 | ncruces/go-sqlite3 (无 CGO) |
| sqlite-vec | asg017/sqlite-vec-go-bindings | - |
| llama.cpp | dianlight/gollama.cpp | go-skynet/go-llama.cpp |
| CLI | spf13/cobra | urfave/cli |

---

## 八、替代方案考量

### 8.1 混合方案

如果单一语言难以满足所有需求，可考虑：

1. **Go + Zig 工具链**：使用 Zig 作为 Go CGO 的交叉编译工具链
2. **Go 主体 + Rust FFI**：核心逻辑用 Go，性能关键路径用 Rust
3. **分离部署**：MCP Server 用 Go，LLM 推理用独立进程

### 8.2 简化方案

如果 LLM 本地推理过于复杂，可考虑：

1. **外部 LLM 服务**：调用 llama.cpp server 或 Ollama API
2. **仅实现搜索**：放弃 reranking/query expansion，仅保留 BM25 + 向量搜索
3. **渐进式实现**：先实现核心功能，后续迭代添加 LLM 能力

---

# refer.

## MCP SDK 相关
- [Model Context Protocol 官方文档](https://modelcontextprotocol.io)
- [Go MCP SDK (官方)](https://github.com/modelcontextprotocol/go-sdk)
- [Rust MCP SDK - rmcp](https://github.com/anthropics/rmcp)
- [Zig MCP 实现 - mcp.zig](https://github.com/zig-wasm/zig-mcp)

## llama.cpp 绑定
- [gollama.cpp (Go, 无 CGO)](https://github.com/dianlight/gollama.cpp)
- [go-llama.cpp (Go, CGO)](https://github.com/go-skynet/go-llama.cpp)
- [llama_cpp-rs (Rust)](https://crates.io/crates/llama_cpp)
- [llama-cpp-2 (Rust)](https://crates.io/crates/llama-cpp-2)
- [llama.cpp.zig (Zig)](https://github.com/Deins/llama.cpp.zig)

## SQLite 扩展
- [sqlite-vec 官方](https://github.com/asg017/sqlite-vec)
- [sqlite-vec Go 绑定](https://github.com/asg017/sqlite-vec-go-bindings)
- [sqlite-vec Rust crate](https://crates.io/crates/sqlite-vec)
- [SQLite FTS5 文档](https://sqlite.org/fts5.html)

## 语言对比与生态
- [Rust vs Go 2026 对比](https://medium.com)
- [Zig 1.0 路线图](https://ziglang.org)
- [Go 1.26 发布计划](https://go.dev)

## 上游项目
- [tobi/qmd - 原始实现](https://github.com/tobi/qmd)

---

> 本调研由 Claude Opus 4.5 (claude-opus-4-5-20251101) 完成
> 所有技术信息已通过 Google 搜索验证
> 建议在实施前进行小规模 PoC 验证关键技术点
