# ABD IM 后端

高性能即时通讯后端，支持万级并发连接。

## 项目结构

```
backend/
├── cmd/server/          # 应用入口
├── config/              # 配置文件
├── internal/
│   ├── api/             # HTTP API (Gin)
│   ├── im/              # WebSocket 即时通讯
│   │   ├── client.go    # 客户端连接管理
│   │   ├── ws_server.go # WebSocket 服务器
│   │   ├── distributor/ # 消息分发器
│   │   ├── pusher/      # 消息推送器
│   │   └── imrepo/      # IM 数据仓库
│   ├── model/           # 数据模型
│   ├── service/         # 业务逻辑层
│   └── pkg/             # 内部公共包
│       ├── database/    # 数据库封装
│       ├── cache/redis/ # Redis 缓存
│       ├── kafka/       # Kafka 消息队列
│       └── prommetrics/ # Prometheus 指标
├── pkg/                 # 外部公共包
└── deployments/         # 部署配置
    └── build/           # Docker/K8s 配置
```

## 技术亮点

- **三层架构** — Gateway / Logic / Pusher 分离，可独立扩展
- **高并发** — sync.Pool 复用连接对象，Channel 串行化避免锁竞争
- **缓存防护** — Singleflight 防止缓存击穿
- **消息可靠** — Kafka 异步持久化 + Seq 机制保证消息有序

## 技术栈

Gin · WebSocket · Kafka · Redis · MySQL · GORM · JWT · Snowflake · Prometheus

## 快速开始

### 本地运行

```bash
# 1. 配置 config/config.yaml 中的中间件地址

# 2. 运行
make run
```

### Docker 部署

```bash
cd deployments/build

# 启动所有服务（包含中间件）
docker-compose up -d

# 查看日志
docker-compose logs -f backend
```

## 配置说明

配置文件: `config/config.yaml`

| 配置项 | 说明 |
|--------|------|
| `server.http_addr` | HTTP API 端口 (默认 :8080) |
| `websocket.addr` | WebSocket 端口 (默认 :8082) |
| `server.metrics_addr` | Prometheus 指标端口 (默认 :9090) |
| `database.*` | MySQL 数据库配置 |
| `redis.*` | Redis 缓存配置 |
| `kafka.*` | Kafka 消息队列配置 |

## API 端点

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/user/register` | 用户注册 |
| POST | `/user/login` | 用户登录 |
| GET | `/user/info` | 获取用户信息 |
| WS | `/ws` | WebSocket 连接 (端口 8082) |
| GET | `/metrics` | Prometheus 指标 (端口 9090) |

## 测试

```bash
# 运行所有测试
go test ./... -v
```

> 注意: 部分测试依赖 Redis/Kafka 服务 (192.168.6.130)
