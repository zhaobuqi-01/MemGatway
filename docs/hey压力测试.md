# hey wrk压力测试

## hey

### 安装

```go
go install github.com/rakyll/hey@latest
```

### 使用

```shell
hey -n 10000 -c 100 http://localhost:8080/
```

`hey` 命令的常用参数如下：

- `-n`：要发出的请求数（默认为 200）。
- `-c`：并发请求的数量（默认为 50）。
- `-t`：超时时间，单位为秒（默认为 20s）。
- `-z`：测试的持续时间。例如：`-z 10s` 表示持续 10 秒。这个参数与 `-n` 互斥，不能同时使用。
- `-m`：HTTP 方法，例如 GET、POST、PUT 等（默认为 GET）。
- `-H`：添加自定义 HTTP 头。例如：`-H "Authorization: Bearer mytoken"`。
- `-D`：发送 POST 请求时的请求体文件路径。
- `-o`：输出格式，可以是 `csv` 或者 `json`。例如：`-o csv`。

### 结果

`hey` 的测试结果分为多个部分，包括摘要、响应时间直方图、延迟分布和详细信息。下面解释了各个部分的意义：

#### 摘要

摘要部分提供了关于测试的总体信息。

- Total：整个测试运行的总时间。
- Slowest：最慢的请求响应时间。
- Fastest：最快的请求响应时间。
- Average：所有请求的平均响应时间。
- Requests/sec：每秒钟的请求处理速率。
- Total data：接收到的总数据量。
- Size/request：每个请求接收到的数据量。

#### 响应时间直方图

响应时间直方图展示了响应时间在不同区间的分布。每个区间的请求数以直方图的形式表示。这有助于了解请求响应时间的分布情况，以便找出潜在的性能问题。

#### 延迟分布

延迟分布部分显示了在不同百分比的请求中，响应时间所处的范围。例如，10% 的请求在某个时间范围内完成。这有助于了解请求延迟的整体情况。

#### 详细信息

详细信息部分提供了有关请求处理过程中各个阶段的平均、最快和最慢时间。

- DNS+dialup：建立连接所需的时间。
- DNS-lookup：进行 DNS 查找所需的时间。
- req write：发送请求所需的时间。
- resp wait：等待响应所需的时间。
- resp read：读取响应所需的时间。

#### 状态码分布

状态码分布部分显示了每个 HTTP 状态码的响应数量。这有助于了解服务器在压力测试期间的状态。

## wrk

###  安装

```shell
sudo apt-get update
sudo apt-get install -y build-essential libssl-dev git
git clone https://github.com/wg/wrk.git
cd wrk
make
sudo cp wrk /usr/local/bin/
```

### 使用

```shell
wrk [options] <URL>
```

常用参数：

- `-t` (Threads)：指定用于发起请求的线程数量。例如：`-t 4` 使用 4 个线程。
- `-c` (Connections)：指定并发连接的数量。例如：`-c 100` 使用 100 个并发连接。
- `-d` (Duration)：指定测试持续时间。可以用秒（s）、分钟（m）或小时（h）表示。例如：`-d 30s` 为 30 秒，`-d 10m` 为 10 分钟，`-d 1h` 为 1 小时。
- `-s` (Script)：指定一个 Lua 脚本，用于自定义请求行为。例如：`-s myscript.lua`。
- `-H` (Header)：添加自定义 HTTP 头。可以多次使用此选项以添加多个头。例如：`-H "Authorization: Bearer mytoken" -H "Accept-Encoding: gzip"`。
- `-R` (Rate)：指定每秒钟的请求速率。例如：`-R 1000` 每秒钟发起 1000 个请求。

## 快速测试服务器最大请求量

### hey

要使用 `hey` 快速测试服务器的最大请求量，可以使用以下命令：

```
hey -n 100000 -c 1000 -z 60s http://example.com/
```

该命令将：

1. 发起 100,000 个请求 (`-n 100000`)
2. 使用 1000 个并发连接 (`-c 1000`)
3. 测试持续 60 秒 (`-z 60s`)
4. 对目标 URL（http://example.com/）进行压力测试

您可以根据需要调整参数以满足您的测试需求。

### wrk

要使用 `wrk` 快速测试服务器的最大请求量，可以使用以下命令：

```
wrk -t 4 -c 1000 -d 60s http://example.com/
```

该命令将：

1. 使用 4 个线程 (`-t 4`)
2. 使用 1000 个并发连接 (`-c 1000`)
3. 测试持续 60 秒 (`-d 60s`)
4. 对目标 URL（http://example.com/）进行压力测试