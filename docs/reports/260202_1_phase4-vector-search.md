# Phase 4 工程报告: 向量搜索

> 日期: 2026-02-02
> Issue: #1
> 模型: Claude Opus 4.5 (claude-opus-4-5-20251101)

---

## 一、完成内容

### 1.1 向量存储

- `embeddings` 表 - 存储向量嵌入 (BLOB 格式)
- 支持多 chunk 索引
- 记录模型和维度信息

### 1.2 向量搜索

- 纯 Go 实现余弦相似度计算
- 无外部依赖，跨平台兼容
- 支持 top-k 结果排序

### 1.3 嵌入接口

- Ollama API 客户端
- 默认模型: `nomic-embed-text`
- 可配置 baseURL 和 model

### 1.4 MCP 工具

| 工具 | 功能 |
|------|------|
| `vector_search` | 语义向量搜索 |

---

## 二、新增文件

```
internal/
├── embed/
│   └── client.go     # Ollama 嵌入客户端
└── store/
    └── vector.go     # 向量存储和搜索
```

---

## 三、技术决策

### sqlite-vec WASM 兼容性问题

sqlite-vec 的 ncruces WASM 绑定存在原子操作兼容性问题：
```
i32.atomic.store invalid as feature "" is disabled
```

**替代方案**: 使用纯 Go 实现
- 向量存储为 BLOB (float32 little-endian)
- 余弦相似度在 Go 中计算
- 无 CGO 依赖，保持跨平台兼容

---

## 四、使用方法

```bash
# 需要运行 Ollama 服务
ollama serve

# 向量搜索 (通过 MCP)
echo '{"method":"tools/call","params":{"name":"vector_search","arguments":{"query":"golang"}}}' | ./gqmd mcp
```

---

## 五、后续优化建议

1. **CLI embed 命令**: 批量生成文档嵌入
2. **混合搜索**: FTS5 + 向量搜索结果融合
3. **本地嵌入**: 集成 gollama.cpp 实现离线嵌入

