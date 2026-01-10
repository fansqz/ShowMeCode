# 第一阶段： 编译阶段
FROM ubuntu:18.04 AS builder

# 设定工作目录
WORKDIR /var/go_debugger

# 设置 Docker 时间为上海时区
RUN apt-get update
RUN apt-get install -y tzdata
RUN ln -fs /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

# 安装go
RUN apt-get install -y software-properties-common || true && \
    add-apt-repository -y ppa:longsleep/golang-backports && \
    apt-get install -y golang-go
RUN go version
# 启动go mod
ENV GO111MODULE=on
COPY go.mod go.sum ./
#RUN go env -w GOPROXY=https://goproxy.io,direct # 加了这行github的actions跑不动= =
RUN go mod download
# 复制项目到镜像中
COPY . .
# 静态编译 Go 程序
RUN GOOS=linux go build -a -installsuffix cgo -o godebugger /var/go_debugger


# 第二阶段：运行阶段
FROM ubuntu:18.04

# 设定工作目录
WORKDIR /var/go_debugger

USER root

# 设置 Docker 时间为上海时区
RUN apt-get update
RUN apt-get install -y tzdata
RUN ln -fs /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

# 安装编译器和其他软件
RUN apt-get install -y software-properties-common || true && \
    add-apt-repository -y ppa:longsleep/golang-backports && \
    apt-get install -y golang-go
    # 验证安装 \
RUN go version
RUN apt-get install -y gcc
RUN apt-get install -y gdb
RUN apt-get install -y python3 default-jre


# 应用程序监听端口
EXPOSE 8080

# 运行应用程序
CMD ["/var/go_debugger/godebugger"] 

