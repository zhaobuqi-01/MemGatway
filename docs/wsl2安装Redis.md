# wsl2安装Redis

### 安装

```sh
#更新ubuntubao
sudo apt update
#安装redis
sudo apt install redis-server
#查看版本号
redis-server --version
```

### 启动redis

```sh
sudo service redis-server start
```

### 检查 redis 是否正常工作

```sh
$ redis-cli
127.0.0.1:6379> ping
PONG
127.0.0.1:6379>
```

### redis 操作

```sh
sudo service redis-server stop
sudo service redis-server start
sudo service redis-server status
sudo service redis-server restart
```

### redis开机启动

```sh
sudo systemctl enable redis
```

- 启动 Redis 服务：`sudo systemctl start redis`
- 停止 Redis 服务：`sudo systemctl stop redis`
- 重启 Redis 服务：`sudo systemctl restart redis`
- 查看 Redis 服务状态：`sudo systemctl status redis`

### 修改redis配置

```sh
#切换到root
su
#修改配置文件
gedit /etc/redis/redis.conf

#允许来自任何 IP 地址的连接
bind 0.0.0.0
#重启生效
sudo service redis-server restart
```

### 查看redis监听的端口

```sh
LISTEN 0      511          0.0.0.0:6379       0.0.0.0:*    users:(("redis-server",pid=13438,fd=6))
```

