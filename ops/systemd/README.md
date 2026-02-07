# gQMD Systemd 服务配置

本目录包含 gQMD 的 systemd 服务配置文件，用于 Linux 系统的服务管理。

## 文件说明

| 文件 | 说明 |
|------|------|
| `ollama.service` | Ollama 本地 LLM 服务（仅供参考） |
| `gqmd-scan.service` | gqmd 文档索引扫描服务 |
| `gqmd-scan.timer` | gqmd 定时扫描触发器 |

## 快速部署

详细说明请参考项目根目录 [README.md](../../README.md) 中的 Linux Systemd 部署章节。

## 配置说明

### gqmd 扫描参数

| 参数 | 值 | 说明 |
|------|-----|------|
| OLLAMA_HOST | `http://127.0.0.1:11434` | Ollama API 地址 |
| GQMD_EMBEDDING_MODEL | `nomic-embed-text` | 嵌入模型名称 |
| 定时执行 | 每天 03:00 | 自动扫描更新索引 |
| 启动执行 | 开机 5 分钟后 | 首次扫描 |

### 自定义配置

部署前请根据实际环境修改 `gqmd-scan.service` 中的：
- `User` / `Group` - 运行服务的用户和组
- `WorkingDirectory` - 工作目录
- `ExecStart` - gqmd 二进制文件路径

## 服务管理命令

| 服务 | 启动 | 停止 | 状态 | 日志 |
|------|------|------|------|------|
| Ollama | `sudo systemctl start ollama` | `sudo systemctl stop ollama` | `sudo systemctl status ollama` | `journalctl -u ollama -f` |
| gqmd 扫描 | `sudo systemctl start gqmd-scan` | - | `sudo systemctl status gqmd-scan` | `journalctl -u gqmd-scan` |
| gqmd 定时器 | `sudo systemctl start gqmd-scan.timer` | `sudo systemctl stop gqmd-scan.timer` | `systemctl list-timers` | - |
