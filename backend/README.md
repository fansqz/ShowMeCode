# FanCode

FanCode是一个可视化OJ系统，愿景是实现一个集判题系统，知识库等教学一体的编程系统。取代传统的算法学习流程。fancode_backend是FanCode的后端仓库。

目前实现的功能如下：

- 简单的沙盒
- 判题系统
- 调试功能





## 启动程序

1. Linux环境

   由于FanCode实现了简单的沙盒，会依赖于Linux的一些功能，比如它的cgroup。所以FanCode必须要在Linux环境下才能启动。如果你的系统不是Linux系统，你可能需要去搭建你的Linux环境才能启动FanCode。下面提供了一些搭建Linux环境的参考：

    - Window下可以使用Linux子系统[什么是适用于 Linux 的 Windows 子系统 | Microsoft Learn](https://learn.microsoft.com/zh-cn/windows/wsl/about)

2. 依赖程序

   FanCode依赖了工具，因为FanCode在本地需要进行编译运行等，需要安装各种语言的编译器。可以运行以下程序在linux上安装需要的工具。

   ```
   // Ubuntu
   apt-get update
   apt-get install gcc
   apt-get install  gdb
   // 待补充
   ```

3. 数据库以及其他依赖

   数据库sql在项目目录下的./fan_code.sql文件下，可以直接创建自己的数据库。当然目前系统的数据库已经在服务器中有部署，所以可以不创建。配置文件中的其他依赖的配置也可以不需要动，直接与运行程序就课可以了。

4. 运行后端程序

   
## docker
1. 创建调试器docker镜像
   ```shell
   docker build -t go-debugger -f Dockerfile-debugger .
   ```
2. 创建并启动showmecode容器
   1. 创建镜像
   ```shell
   docker build -t showmecode -f Dockerfile .
   ```
   2. 启动容器，启动docker容器的时候需要使用`--privileged`给容器提高的权限，并且通过/var/run/docker.sock:/var/run/docker.sock将宿主机的Docker套接字挂载到容器内部，实现docker in docker

   ```
   docker run -d -v /usr/bin/docker:/usr/bin/docker -v /var/run/docker.sock:/var/run/docker.sock -v /var/fanCode:/var/fanCode --network=host --name showmecode showmecode
   ```