# wsl2安装Mysql

```shell
sudo apt update
#安装mysql-server
sudo apt install mysql-server、
#查看mysql版本
mysql --version
#启动mysql
sudo /etc/init.d/mysql start
#以root用户进入mysql
sudo mysql --user=root mysql
```
### 获取所有的权限

```mysql
mysql> UPDATE mysql.user SET authentication_string=null WHERE User='root';
mysql> flush privileges; 
```

### 修改密码

```mysql
mysql> ALTER USER 'root'@'localhost' IDENTIFIED WITH mysql_native_password BY '123456';
mysql> flush privileges;
```

### 验证mysql密码修改是否成功

```sh
# 关掉 mysql 所有的进程
$ sudo killall -u mysql

# 重启 mysql
$ sudo /etc/init.d/mysql start
Starting mysql (via systemctl): mysql.service.

# 进入mysql
$ mysql -uroot -p123456
>mysql
```

### 登录

```shell
sudo mysql -u root -p
```

### 设置开机启动

```sh
$ sudo systemctl enable mysql
Synchronizing state of mysql.service with SysV service script with /lib/systemd/systemd-sysv-install.
Executing: /lib/systemd/systemd-sysv-install enable mysql
```

### mysql操作

```sh
# 查看 mysql 服务运行状态
$ sudo service mysql status

# 开启 mysql 服务
$ sudo service mysql start

# 停止 mysql 服务
$ sudo service mysql stop
```



# 参考

##### 第一步 ：更新软件包

```bash
$ sudo apt update 
```

##### 第二步 ：安装 mysql8.0

```bash
$ sudo apt install mysql-server-8.0 -y
```

##### 第三步：查询安装的 mysql 版本

```bash
$ mysql --version
mysql  Ver 8.0.32-0ubuntu0.20.04.2 for Linux on x86_64 ((Ubuntu))
或者
$ mysql -V
mysql  Ver 8.0.32-0ubuntu0.20.04.2 for Linux on x86_64 ((Ubuntu))
```

##### 第四步：停止 mysql 服务

```bash
$ sudo /etc/init.d/mysql stop
```

##### 第五步：创建特定mysql运行目录

```bash
$ sudo mkdir /var/run/mysqld
$ sudo chown mysql /var/run/mysqld
```

##### 第六步：取消授权登录的限制，允许你可以匿名登录

```bash
$ sudo mysqld_safe --skip-grant-tables&
```

##### 第七步：进入 mysql

```bash
$ sudo mysql --user=root mysql
```

##### 第八步：获取所有的权限

```sql
mysql> UPDATE mysql.user SET authentication_string=null WHERE User='root';
mysql> flush privileges; 
```

##### 第九步：修改 mysql 密码

```sql
mysql> ALTER USER 'root'@'localhost' IDENTIFIED WITH mysql_native_password BY '123456';
mysql> flush privileges;
```

##### 第十步：退出 mysql

快捷键：Ctrl+z

```bash
mysql > exit
Bye
```

##### 第十一步：验证 mysql 密码修改是否成功

```bash
# 关掉 mysql 所有的进程
$ sudo killall -u mysql

# 重启 mysql
$ sudo /etc/init.d/mysql start
Starting mysql (via systemctl): mysql.service.

# 进入mysql
$ mysql -uroot -p123456
>mysql
```

##### 第十二步：设置开机启动 mysql 服务

```bash
# 设置开机启动 mysql 服务
$ sudo update-rc.d -f mysql defaults

# 查看是否开机启动 mysql 服务
$ sudo service mysql status
mysql.service - MySQL Community Server
     Loaded: loaded (/lib/systemd/system/mysql.service; enabled; vendor preset: enabled)
     Active: active (running) since Fri 2023-02-03 02:53:16 UTC; 3min 18s ago
   Main PID: 6212 (mysqld)
     Status: "Server is operational"
      Tasks: 39 (limit: 9406)
     Memory: 367.0M
     CGroup: /system.slice/mysql.service
             └─6212 /usr/sbin/mysqld

Feb 03 02:53:15 zgxt systemd[1]: Starting MySQL Community Server...
Feb 03 02:53:16 zgxt systemd[1]: Started MySQL Community Server.

# 取消开机启动 mysql 服务
$ sudo update-rc.d -f mysql remove
```

##### 第十三步：mysql 服务操作

```bash
# 查看 mysql 服务运行状态
$ sudo service mysql status

# 开启 mysql 服务
$ sudo service mysql start

# 停止 mysql 服务
$ sudo service mysql stop
```

##### 第十四步：mysql 开启 root 用户远程连接

> 注意：mysql 出于安全方面考虑默认只允许本机(localhost, 127.0.0.1)来连接访问

```bash
# 进入 mysql
$ mysql -uroot -p123456

# 切换到 mysql 数据库
mysql> use mysql

# 查看 root 用户权限
mysql> select user,host,plugin from user;
+------------------+-----------+-----------------------+
| user             | host      | plugin                |
+------------------+-----------+-----------------------+
| debian-sys-maint | localhost | caching_sha2_password |
| mysql.infoschema | localhost | caching_sha2_password |
| mysql.session    | localhost | caching_sha2_password |
| mysql.sys        | localhost | caching_sha2_password |
| root             | localhost | mysql_native_password |
+------------------+-----------+-----------------------+
```

可以看到 root 用户只有 localhost 本机权限，就是只有本机能访问

```bash
# 给 root 用户授权
mysql> update user set host = '%' where user ='root';
mysql> flush privileges;

# 查看 root 用户权限
mysql> select user,host,plugin from user;
+------------------+-----------+-----------------------+
| user             | host      | plugin                |
+------------------+-----------+-----------------------+
| root             | %         | mysql_native_password |
| debian-sys-maint | localhost | caching_sha2_password |
| mysql.infoschema | localhost | caching_sha2_password |
| mysql.session    | localhost | caching_sha2_password |
| mysql.sys        | localhost | caching_sha2_password |
+------------------+-----------+-----------------------+
```

注意：其中 % 表示任意远程 IP 可以访问

修改 mysql 配置文件，将默认的 bind-address=127.0.0.1 修改如下：

```bash
# 修改 mysql IP地址绑定
$ cd /etc/mysql/mysql.conf.d
bind-address  = 0.0.0.0

## 关掉 mysql 所有的进程
$ sudo killall -u mysql

# 重启 mysql
$ sudo /etc/init.d/mysql start

# 查询 mysql 服务监听端口
$ sudo ss -tulnp | grep LISTEN | grep mysql
tcp    LISTEN  0       151                    0.0.0.0:3306         0.0.0.0:*     users:(("mysqld",pid=9340,fd=23)) 
```

注意：0.0.0.0:3306 表示支持监听远程连接

使用 Navicat 远程连接 mysql 数据库：

![image-20230203122853322](https://gitee.com/binbingg/pic-bed/raw/master/img/image-20230203122853322.png)

