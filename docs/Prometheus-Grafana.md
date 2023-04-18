# win

在Windows 10上安装Prometheus和Grafana的步骤如下：

1. 安装Prometheus

   - 访问Prometheus官方下载页面：https://prometheus.io/download/

   - 从"prometheus"部分下载适用于Windows的最新版本（`.zip` 文件）。

   - 将下载的`.zip`文件解压到一个目录，例如 `C:\prometheus`。

   - 在解压缩后的目录中，找到`prometheus.yml`文件，根据需要对其进行配置。例如，可以添加监控目标、更改抓取间隔等。

   - 打开命令提示符，转到Prometheus目录，运行以下命令启动Prometheus：

     ```
     arduinoCopy codecd C:\prometheus
     .\prometheus.exe --config.file=prometheus.yml
     ```

   - Prometheus现在已经运行在 `http://localhost:9090`。

2. 安装Grafana

   - 访问Grafana官方下载页面：https://grafana.com/grafana/download

   - 从"Standalone Windows Binaries"部分下载适用于Windows的最新版本（`.zip` 文件）。

   - 将下载的`.zip`文件解压到一个目录，例如 `C:\grafana`。

   - 在解压缩后的目录中，找到`conf`目录，根据需要对`defaults.ini`（或`custom.ini`）进行配置。

   - 打开命令提示符，转到Grafana目录，运行以下命令启动Grafana：

     ```
     bashCopy codecd C:\grafana\bin
     .\grafana-server.exe
     ```

   - Grafana现在已经运行在 `http://localhost:3000`。

现在，Prometheus和Grafana都已经在Windows 10上安装并运行。您可以按照前面的回答中提到的步骤将Prometheus与Grafana集成，并将其嵌入到vue-element-admin中。

# linux

在WSL2中安装Prometheus和Grafana的步骤如下：

1. 安装Prometheus

   - 打开WSL2终端。

   - 下载最新版本的Prometheus：

     ```
     wget https://github.com/prometheus/prometheus/releases/download/v2.31.1/prometheus-2.31.1.linux-amd64.tar.gz
     ```

   - 解压下载的文件：

     ```
     tar -xvf prometheus-2.31.1.linux-amd64.tar.gz
     ```

   - 进入解压后的目录，并根据需要编辑`prometheus.yml`配置文件。

     ```
     prometheus-2.31.1.linux-amd64
     vi prometheus.yml
     ```

   - 启动Prometheus：

     ```
     ./prometheus --config.file=prometheus.yml
     ```

   - Prometheus现在已经运行在 `http://localhost:9090`。

2. 安装Grafana

   - 打开另一个WSL2终端。

   - 首先，添加Grafana的APT存储库：

     ```
     curl -s https://packages.grafana.com/gpg.key | sudo apt-key add -
     echo "deb https://packages.grafana.com/oss/deb stable main" | sudo tee -a /etc/apt/sources.list.d/grafana.list
     ```

   - 更新APT存储库并安装Grafana：

     ```
     apt-get update
     sudo apt-get install grafana
     ```

   - 启动Grafana服务：

     ```
     systemctl daemon-reload
     sudo systemctl start grafana-server
     ```

   - Grafana现在已经运行在 `http://localhost:3000`。

要在Windows 10中打开这些页面，您需要获取WSL2的IP地址。可以通过运行以下命令获得：

```
ip addr show eth0 | grep "inet\b" | awk '{print $2}' | cut -d/ -f1
```

假设您的WSL2 IP地址为`192.168.1.100`，您可以通过以下网址在Windows 10中访问Prometheus和Grafana：

- Prometheus：`http://192.168.1.100:9090`
- Grafana：`http://192.168.1.100:3000`

现在，Prometheus和Grafana都已经在WSL2上安装并运行。您可以按照前面的回答中提到的步骤将Prometheus与Grafana集成，并将其嵌入到vue-element-admin中。



对于WSL2，它不支持Systemd。在这种情况下，您需要手动启动Grafana服务。请按照以下步骤操作：

1. 启动Grafana服务：

   ```
   sudo grafana-server --config=/etc/grafana/grafana.ini --homepath=/usr/share/grafana
   ```

2. Grafana现在应该已经运行在 `http://localhost:3000`。但是，您需要获取WSL2的IP地址才能在Windows 10上访问它。您可以通过运行以下命令获得：

   ```
   ip addr show eth0 | grep "inet\b" | awk '{print $2}' | cut -d/ -f1
   ```

假设您的WSL2 IP地址为`192.168.1.100`，您可以通过以下网址在Windows 10中访问Grafana：

- Grafana：`http://192.168.1.100:3000`

现在Grafana已经在WSL2上安装并运行，您可以按照前面的回答中提到的步骤将Prometheus与Grafana集成，并将其嵌入到vue-element-admin中。

# docker

在Docker中安装Grafana和Prometheus很简单。首先确保已经安装了Docker，然后按照以下步骤操作：

1. 运行Prometheus容器：

   创建一个配置文件`prometheus.yml`，可以使用以下示例配置：

   ```
   yamlCopy codeglobal:
     scrape_interval: 15s
     external_labels:
       monitor: 'codelab-monitor'
   
   scrape_configs:
     - job_name: 'prometheus'
       static_configs:
         - targets: ['localhost:9090']
   ```

   在包含`prometheus.yml`文件的目录中，运行以下命令来启动Prometheus容器：

   ```
   docker run -d --name prometheus -p 9090:9090 -v $(pwd)/prometheus.yml:/etc/prometheus/prometheus.yml prom/prometheus
   ```

   这将在后台运行Prometheus容器，并将主机上的9090端口映射到容器的9090端口。Prometheus服务器将使用提供的配置文件运行。

2. 运行Grafana容器：

   使用以下命令启动Grafandocker run -d --name grafana -p 3000:3000 grafana/grafana

   这将在后台运行Grafana容器，并将主机上的3000端口映射到容器的3000端口。

现在，Prometheus和Grafana已经在Docker容器中运行。您可以通过访问 [http://localhost:3000](http://localhost:3000/)（默认用户名和密码均为`admin`）来访问Grafana UI，并通过访问 [http://localhost:9090](http://localhost:9090/) 来访问Prometheus UI。

接下来，在Grafana中配置Prometheus数据源，然后就可以创建和查看仪表板了。