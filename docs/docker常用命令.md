# docker常用命令

Docker提供了一系列命令来管理容器、镜像、网络和数据卷。以下是一些常用的Docker命令：

1. `docker run` - 创建并启动一个新容器 例如：`docker run -d -p 80:80 --name my_container nginx`
2. `docker ps` - 查看当前正在运行的容器 例如：`docker ps` 查看所有容器（包括已停止的）：`docker ps -a`
3. `docker stop` - 停止一个正在运行的容器 例如：`docker stop my_container`
4. `docker start` - 启动一个已停止的容器 例如：`docker start my_container`
5. `docker restart` - 重启一个容器 例如：`docker restart my_container`
6. `docker rm` - 删除一个容器 例如：`docker rm my_container`
7. `docker images` - 查看本地镜像列表 例如：`docker images`
8. `docker rmi` - 删除本地镜像 例如：`docker rmi nginx`
9. `docker pull` - 从远程仓库拉取镜像 例如：`docker pull nginx`
10. `docker build` - 使用Dockerfile构建镜像 例如：`docker build -t my_image .`
11. `docker push` - 将镜像推送到远程仓库 例如：`docker push my_image`
12. `docker exec` - 在运行的容器中执行命令 例如：`docker exec -it my_container bash`
13. `docker logs` - 查看容器日志 例如：`docker logs my_container`
14. `docker inspect` - 查看容器详细信息 例如：`docker inspect my_container`
15. `docker cp` - 在容器与主机之间复制文件或目录 例如：`docker cp my_container:/path/to/file /path/on/host`
16. `docker network` - 管理Docker网络 例如：`docker network ls` 创建网络：`docker network create my_network` 删除网络：`docker network rm my_network`
17. `docker volume` - 管理Docker数据卷 例如：`docker volume ls` 创建数据卷：`docker volume create my_volume` 删除数据卷：`docker volume rm my_volume`

这些命令涵盖了Docker的基本操作，实际使用中可能会遇到更多其他命令和选项。你可以通过阅读Docker文档或使用`docker <command> --help`来查找详细的命令用法。