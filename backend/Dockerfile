# 第一阶段： 编译阶段
FROM golang:1.23-alpine AS builder

# 设定工作目录
WORKDIR /var/fancode

# 启动go mod
ENV GO111MODULE=on
COPY go.mod go.sum ./
RUN go env -w GOPROXY=https://goproxy.cn
RUN go mod download
# 复制项目到镜像中
COPY . .
# 静态编译 Go 程序
RUN GOOS=linux go build -a -installsuffix cgo -o fancode /var/fancode


# 第二阶段：运行阶段
FROM ubuntu:18.04

# 设定工作目录
WORKDIR /var/fancode

USER root

# 设置 Docker 时间为上海时区
RUN apt-get update
RUN apt-get install -y tzdata
RUN ln -fs /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

RUN apt-get install -y ca-certificates

# 从构建阶段拷贝编译好的二进制文件
COPY --from=builder /var/fancode/fancode .
# 从构建阶段拷贝配置文件
COPY --from=builder /var/fancode/conf ./conf
# 拷贝docker-debugger的dockerfile文件
COPY --from=builder /var/fancode/Dockerfile-debugger ./Dockerfile-debugger

# 配置环境变量，docker版本DOCKER_API_VERSION
ENV DOCKER_API_VERSION=1.43

# 应用程序监听端口
EXPOSE 8080

# 运行应用程序
CMD ["/var/fancode/fancode"]


