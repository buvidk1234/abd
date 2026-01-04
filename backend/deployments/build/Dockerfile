# ============================================
# 阶段1: 构建阶段 (使用完整的 Go 环境编译代码)
# ============================================
FROM golang:1.24-alpine AS builder
 # 设置 Go 代理，加速依赖下载
ENV GOPROXY=https://goproxy.cn,direct

 # 设置工作目录
WORKDIR /app
 # 先复制 go.mod 和 go.sum 文件
COPY go.mod go.sum ./
 # 下载依赖，加速后续编译
RUN go mod download
 # 再复制源码
COPY . ./
 # 编译，CGO_ENABLED=0 表示纯 Go 编译，不依赖 C 库
RUN CGO_ENABLED=0 go build -o server ./cmd/server

# ============================================

# 阶段2: 运行阶段 (只包含运行所需的最小内容)
# ============================================
# 使用极小的 Alpine 镜像 (~5MB)
FROM alpine:3.21

 # 设置工作目录
WORKDIR /app
 # 从构建阶段复制编译好的二进制文件
COPY --from=builder /app/server .
 # 复制配置文件
COPY config/config.yaml ./config/

 # 声明容器监听的端口
EXPOSE 8080
 # 容器启动时执行的命令
ENTRYPOINT ["./server"]