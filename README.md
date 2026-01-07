# pionex-administrative-sys

基于 Go 的管理后台系统，提供 REST API 服务。

## 技术栈

- **Web 框架**: Gin
- **数据库**: SQLite + GORM
- **日志**: Zap
- **认证**: JWT

## 快速开始

### 运行

```bash
# 默认启动（端口 8080）
go run main.go

# 自定义端口
go run main.go -p 3000

# 守护进程模式
go run main.go -d

# 启用文件日志
go run main.go -fl
```

### 构建

```bash
go build -o pionex-administrative-sys
```

### 测试

```bash
go test ./...
go test -v ./db/...
```

## 项目结构

```
├── main.go                 # 入口，优雅关闭处理
├── server/
│   ├── server.go           # Gin 引擎配置
│   ├── handler/            # 路由处理器
│   │   ├── user/           # 用户认证、管理
│   │   ├── coupon/         # 优惠券管理
│   │   └── my_coupon/      # 我的优惠券
│   └── middleware/         # 中间件（日志、恢复、JWT）
├── db/                     # 数据模型和数据访问
│   ├── db.go               # 数据库连接
│   ├── user.go             # 用户模型
│   ├── role.go             # 角色定义
│   ├── coupon.go           # 优惠券模型
│   └── coupon_type.go      # 优惠券类型
├── utils/                  # 工具函数
│   ├── app/                # 命令行解析、守护进程
│   ├── logger/             # 日志配置
│   ├── crypto.go           # 加密工具
│   ├── jwt.go              # JWT 工具
│   └── resp.go             # 响应封装
└── static/                 # 静态资源
```

## API 接口

所有接口在 `/api/v1` 路径下：

| 路径 | 说明 |
|------|------|
| `GET /health` | 健康检查 |
| `/api/v1/user/*` | 用户相关接口 |
| `/api/v1/coupon/*` | 优惠券管理 |
| `/api/v1/my_coupon/*` | 我的优惠券 |

## 配置

| 配置项 | 默认值 | 说明 |
|--------|--------|------|
| 端口 | 8080 | 使用 `-p` 参数覆盖 |
| 数据目录 | `~/.pas/data/` | 可通过 `PAS_HOME` 环境变量覆盖 |
| 日志目录 | `~/.pas/logs/` | - |

## 默认账户

- 账号: `admin`
- 密码: `123456`
