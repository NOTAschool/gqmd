# Phase 2 工程报告: 文档索引与搜索

> 日期: 2026-02-02
> Issue: #1
> 模型: Claude Opus 4.5 (claude-opus-4-5-20251101)

---

## 一、完成内容

### 1.1 数据库 Schema 扩展

新增表结构：
- `content` - 内容寻址存储
- `documents` - 文档元数据
- `documents_fts` - FTS5 全文搜索虚拟表

### 1.2 Collection 管理

| 方法 | 功能 |
|------|------|
| `AddCollection` | 创建集合 |
| `ListCollections` | 列出所有集合 |
| `GetCollection` | 获取集合详情 |
| `RemoveCollection` | 删除集合及其文档 |

### 1.3 文档索引

- `IndexDocument` - 索引单个文档
- `ScanCollection` - 扫描目录并索引文档
- 支持 `**/*.md` 风格的 glob 模式
- 自动提取 Markdown 标题

### 1.4 搜索功能

- FTS5 全文搜索 (BM25 排序)
- `Search` - 搜索文档
- `Get` - 获取单个文档
- `MultiGet` - 批量获取文档 (带字节限制)

### 1.5 MCP 工具

| 工具 | 功能 |
|------|------|
| `status` | 显示索引状态 |
| `search` | FTS5 全文搜索 |
| `get` | 获取单个文档 |
| `multi_get` | 批量获取文档 |

---

## 二、新增文件

```
internal/store/
├── db.go          # 扩展: Collection/Document/Search 方法
└── scanner.go     # 新增: 目录扫描和文档索引
```

---

## 三、测试结果

```
=== RUN   TestOpenPath
--- PASS: TestOpenPath (0.33s)
=== RUN   TestGetStatus
--- PASS: TestGetStatus (0.01s)
=== RUN   TestCollectionCRUD
--- PASS: TestCollectionCRUD (0.01s)
=== RUN   TestIndexAndSearch
--- PASS: TestIndexAndSearch (0.01s)
PASS
```

---

## 四、后续优化建议

1. **Phase 3 准备**: 实现 CLI 命令 (add/scan/search)
2. **增量索引**: 基于文件 hash 跳过未修改文件
3. **删除检测**: 标记已删除文件为 inactive

---

## 五、使用示例

```go
// 创建集合
store.AddCollection("docs", "/path/to/docs", "**/*.md")

// 扫描并索引
store.ScanCollection("docs")

// 搜索
results, _ := store.Search("golang", 10)

// 获取文档
doc, content, _ := store.Get("docs", "readme.md")
```

