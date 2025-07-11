FROM golang:1.23.4 AS builder

# 设置工作目录
WORKDIR /app
ENV GOPROXY=https://mirrors.aliyun.com/goproxy/
# 复制源代码
COPY . .
RUN go mod tidy && go build -o basic-app

# 创建最终运行时镜像
FROM ubuntu:latest
WORKDIR /app
# 拷贝构建产物
COPY --from=builder /app/basic-app /app/basic-app
COPY ./config /app/config
RUN chmod +x /app/basic-app

EXPOSE 8080
# 设置入口点
ENTRYPOINT ["./basic-app"]
