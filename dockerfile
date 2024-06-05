# 使用官方的 Golang 镜像作为构建阶段
FROM golang:1.22.0 as builder

# 设置工作目录
WORKDIR /app

# 复制 go.mod 和 go.sum 文件，并下载依赖项
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 编译 Go 应用程序，确保它是静态链接的
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o main .


# 使用一个更小的基础镜像作为运行阶段
FROM alpine:latest

# 安装 ca-certificates 以确保 HTTPS 请求正常工作
RUN apk --no-cache add ca-certificates

# 设置工作目录
WORKDIR /root/

# 从构建阶段复制编译好的 Go 应用程序
COPY --from=builder /app/main .

COPY --from=builder /app/config.dev.yml .
COPY --from=builder /app/config.yml .

# 暴露应用程序运行的端口
EXPOSE 8080

# 运行应用程序并传递命令行参数
CMD ["./main"]
