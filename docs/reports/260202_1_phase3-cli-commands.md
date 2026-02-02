# Phase 3 工程报告: CLI 命令实现

> 日期: 2026-02-02
> Issue: #1
> 模型: Claude Opus 4.5 (claude-opus-4-5-20251101)

---

## 一、完成内容

### 1.1 CLI 命令

| 命令 | 功能 |
|------|------|
| `gqmd add <name> <path>` | 添加集合 |
| `gqmd list` | 列出所有集合 |
| `gqmd remove <name>` | 删除集合 |
| `gqmd scan [name]` | 扫描并索引文档 |
| `gqmd search <query>` | 搜索文档 |
| `gqmd status` | 显示索引状态 |

### 1.2 命令参数

- `add --pattern, -p` : 文件匹配模式 (默认 `**/*.md`)
- `search --limit, -n` : 最大结果数 (默认 10)

---

## 二、新增文件

```
internal/cli/
├── add.go      # 添加集合
├── list.go     # 列出集合
├── remove.go   # 删除集合
├── scan.go     # 扫描索引
└── search.go   # 搜索文档
```

---

## 三、测试结果

```bash
$ ./gqmd add docs ./docs
Added collection "docs" -> /opt/src/41490/NOTA/gqmd/docs

$ ./gqmd list
docs -> /opt/src/41490/NOTA/gqmd/docs (**/*.md)

$ ./gqmd scan docs
Scanning docs...
Added: 3, Errors: 0

$ ./gqmd search "Phase"
1. docs/reports/260202_1_phase2-document-indexing.md
   Phase 2 工程报告: 文档索引与搜索
2. docs/reports/260202_1_phase1-core-framework.md
   Phase 1 工程报告: 核心框架实现

$ ./gqmd status
gqmd Index Status:
  Database: ~/.cache/gqmd/index.sqlite
  Total documents: 3
  Collections: 1
```

---

## 四、后续优化建议

1. **Phase 4 准备**: 向量搜索 (sqlite-vec + GGUF 嵌入)
2. **增量扫描**: 基于 hash 跳过未修改文件
3. **交互模式**: 支持 REPL 交互搜索

