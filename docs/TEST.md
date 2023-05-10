# TEST

## LSOF的使用

`lsof`（list open files）是一个用于显示当前系统中打开的文件和网络套接字的实用工具。在UNIX和类UNIX系统（如Linux）上，一切皆文件，因此`lsof`可以帮助查找文件、目录、设备文件和网络套接字等资源的使用情况。

`lsof`命令有许多参数可用，以下是一些常用的参数及其用途：

1. `-a`：将多个条件组合在一起，并使用逻辑AND进行过滤。

```sh
lsof -u <username> -a -i <protocol>
```

1. `-c <name>`：列出名为`<name>`的进程打开的文件。这里的`<name>`可以是一个完整的进程名，也可以是一个简写。

```sh
lsof -c ssh
```

1. `-i <address>`：列出与指定网络地址相关的文件。`<address>`可以是端口号、主机名、IPv4或IPv6地址。

```sh
lsof -i :80
```

1. `-p <PID>`：列出指定进程ID（PID）打开的文件。

```sh
lsof -p 12345
```

1. `-s <protocol>:<state>`：显示处于特定状态的网络文件。`<protocol>`可以是TCP或UDP，`<state>`可以是网络连接的状态（如：LISTEN, ESTABLISHED, CLOSE_WAIT等）。

```sh
lsof -i tcp -sTCP:LISTEN
```

1. `-t`：仅显示文件描述符（通常是进程ID）。

```sh
lsof -t -i :80
```

1. `-u <username>`：列出特定用户打开的文件。`<username>`可以是一个完整的用户名，也可以是一个用户ID。

```sh
lsof -u john
```

这些参数可以帮助您根据不同条件筛选出`lsof`的输出结果。有时，您可能需要组合多个参数来获取您需要的信息。例如，要查找用户名为`john`且正在监听TCP端口80的进程，可以使用以下命令：

```sh
lsof -u john -a -i :80 -sTCP:LISTEN
```

## HTTP_PROXY

#### 开启真实服务器

两台服务器演示**负载均衡**

```go
go run /test/real_server.go
```

该命令启动了监听在127.0.0.1:2003和127.0.0.1:2004两台服务器

#### 结束其中的一台服务器：

结束一台服务器会来演示**服务发现**

1. 使用`lsof`命令找到监听特定端口的进程

```go
lsof -i :2003
```

2. 使用`kill`命令发送一个SIGTERM信号（代码15）来优雅地关闭进程	

```shell
kill -15 <PID>
```

如果服务器没有立即关闭，您也可以尝试使用SIGKILL信号（代码9）强制关闭：

```shell
kill -9 <PID>
```

#### 限流演示

```shell
hey -n 1000 -c 100  -t 10 -z 10s http://127.0.0.1:8080/test_http_string/abbb
```

```
telnet 127.0.0.1 8011
 grpcurl -plaintext -proto pingpong.proto -d '{"message": "ping"}' localhost:50051 pingpong.PingPong/Ping
```

