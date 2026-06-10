# memos-cli

> **Go ≥ 1.21 要求** | 跨平台：Windows / Linux / macOS

**⚠️ 本项目无 UI**

memos-cli 是被 skill / agent 调用的 shell 工具，可观察入口在 stdout / stderr / 退出码。无 Web UI / TUI / 图形界面。

## 概述

memos-cli 提供了对 Memos HTTP API 的完整命令行封装，支持备忘录、评论、附件等全部功能。

## 项目结构

```
memos-cli/
├── cmd/                  # 命令定义
│   ├── root.go           # 根命令
│   ├── memo.go           # 备忘录命令
│   ├── comment.go        # 评论命令
│   ├── attachment.go     # 附件命令
│   └── whoami.go         # 用户信息命令
├── internal/             # 私有实现
│   ├── client/           # API 客户端
│   └── config/           # 配置管理
├── pkg/                  # 可复用的公共包
│   ├── config/           # 配置加载
│   ├── httpclient/       # HTTP 客户端
│   ├── output/           # 输出格式化
│   └── version/          # 版本信息
├── integration/          # 集成测试
├── main.go               # 入口文件
├── go.mod
├── Makefile
└── README.md
```

## 安装

### 从源码安装

```bash
git clone https://github.com/ANIAN0/memos-cli.git
cd memos-cli
go build -o memos-cli .
```

### 使用 go install

```bash
go install github.com/ANIAN0/memos-cli@latest
```

### 使用 Makefile

```bash
make install
```

## 配置

### 配置文件位置（按优先级）

1. `--config <path>` 命令行参数
2. `MEMOS_CLI_CONFIG` 环境变量
3. 二进制同级目录的 `config.yaml`（项目安装模式）
4. 用户目录：`~/.config/memos-cli/config.yaml`（Unix）或 `%APPDATA%\memos-cli\config.yaml`（Windows）

### 配置文件示例

```yaml
version: 1
instance_url: "https://memos.example.com"
access_token: "${MEMOS_TOKEN}"  # 支持环境变量插值
default_page_size: 10
default_visibility: "PRIVATE"
```

### 环境变量插值

配置文件中支持 `${ENV_VAR}` 格式的环境变量插值：

```yaml
access_token: "${MEMOS_TOKEN}"
```

运行前设置环境变量：

```bash
export MEMOS_TOKEN="your-token"
```

## 子命令

### 备忘录（Memo）

| 命令 | 说明 |
|------|------|
| `memo create` | 创建新备忘录 |
| `memo get <id>` | 获取备忘录详情 |
| `memo list` | 列出备忘录 |
| `memo update <id>` | 更新备忘录 |
| `memo delete <id>` | 删除备忘录 |
| `memo search <query>` | 搜索备忘录 |

### 评论（Comment）

| 命令 | 说明 |
|------|------|
| `comment list <memo-id>` | 列出评论 |
| `comment create <memo-id>` | 创建评论 |

### 附件（Attachment）

| 命令 | 说明 |
|------|------|
| `attachment upload <file>` | 上传附件 |
| `attachment list` | 列出附件 |
| `attachment get <id>` | 下载附件 |
| `attachment delete <id>` | 删除附件 |

## 全局选项

| 选项 | 说明 |
|------|------|
| `--config <path>` | 指定配置文件路径 |
| `--json` | 输出 JSON 格式（可被 `jq` 解析） |
| `--verbose, -v` | 详细日志到 stderr |
| `--timeout <seconds>` | HTTP 请求超时（默认 60） |
| `--no-color` | 禁用颜色输出 |
| `--version` | 输出版本信息 |
| `--help, -h` | 显示帮助 |

## 退出码约定

| 退出码 | 含义 | 触发条件 |
|--------|------|----------|
| `0` | 成功 | 请求成功 |
| `1` | 客户端错误 | Memos 错误码 3/5/7/16 |
| `2` | 服务端错误 | 其他 Memos 错误码 |
| `3` | 网络错误 | DNS 失败、连接超时、连接拒绝 |
| `4` | 配置错误 | 配置文件不存在、字段缺失、环境变量未设置 |

错误详情输出到 stderr，成功数据输出到 stdout。

## 使用示例

```bash
# 创建备忘录
memos-cli memo create --content "Hello World" --visibility PRIVATE

# 列出备忘录
memos-cli memo list --page-size 10

# 搜索备忘录
memos-cli memo search "Hello"

# 上传附件
memos-cli attachment upload ./image.png

# 使用 JSON 输出
memos-cli memo list --json | jq '.items'
```

## 开发

```bash
# 运行测试
make test

# 清理构建产物
make clean
```

## 许可证

MIT License
