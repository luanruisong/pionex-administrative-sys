# CLAUDE.md

本文件为 Claude Code (claude.ai/code) 提供代码库操作指南。

## 项目概述

基于 Go 的管理后台系统，提供 REST API。使用 Gin 框架、SQLite + GORM，实现用户管理和权限控制。

## 常用命令

```bash
# 运行应用
go run main.go
go run main.go -p 3000        # 自定义端口
go run main.go -d             # 守护进程模式
go run main.go -fl            # 启用文件日志

# 构建
go build -o pionex-administrative-sys

# 测试
go test ./...
go test -run TestName ./path/to/package
go test -v ./db/...           # 详细测试某个包

# 代码质量
go fmt ./...
go vet ./...
go mod tidy
```

## 架构

### 分层结构

```
main.go                    入口，优雅关闭处理
├── server/
│   ├── server.go          Gin 引擎配置，HTTP 服务生命周期
│   ├── handler/           路由注册，按业务域划分
│   │   └── user/          认证、用户管理
│   └── middleware/        日志、恢复、JWT 认证
├── db/                    GORM 模型和数据访问函数
│   ├── db.go              SQLite 连接，自动迁移
│   └── user.go            用户模型，位掩码权限
└── utils/
    ├── app/               命令行参数解析，守护进程模式
    ├── logger/            基于 Zap 的日志，集成 GORM
    └── *.go               加密、JWT、响应工具
```

### 关键模式

**API 路由**：所有 API 路由在 `/api/v1` 下，按业务域组织处理器。健康检查接口：`/health`。

**用户角色**：位掩码权限定义在 `db/user.go`：
- `RoleAdmin=1`（管理员）
- `RoleLogin=2`（登录）

**响应模式**：使用 `utils.Resp(code, msg, data).Success(c)` 保持 API 响应格式统一。

**数据库**：SQLite，启动时自动迁移。模型通过 `TableName()` 方法定义表名。

### 环境与配置

- 默认端口：8080（使用 `-p` 参数覆盖）
- 数据目录：`~/.pas/data/`（可通过 `PAS_HOME` 环境变量覆盖）
- 日志目录：`~/.pas/logs/`
- 默认管理员：账号 `admin`，密码 `123456`