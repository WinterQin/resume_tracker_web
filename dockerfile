# 构建阶段
FROM golang:1.23.7-alpine3.21 AS builder

# 设置工作目录
WORKDIR /app

# 复制 go mod 文件
COPY go.* ./
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# 最终阶段
FROM alpine:latest

# 安装基本工具和时区数据
RUN apk --no-cache add tzdata

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/main .

# 暴露端口
EXPOSE 8080

# 设置默认环境变量
ENV DB_HOST=localhost \
    DB_PORT=3306 \
    DB_USER=resume_winter \
    DB_PASSWORD=resume_qwt123456 \
    DB_NAME=internship_manager \
    DB_CHARSET=utf8mb4 \
    DB_MAX_IDLE_CONNS=10 \
    DB_MAX_OPEN_CONNS=100 \
    JWT_KEY=winter-key \
    SERVER_PORT=8080

# 运行应用
CMD ["./main"]
