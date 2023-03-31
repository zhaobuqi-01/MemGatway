# redis,mysql脚本

```sh
#!/bin/bash

redis_start() {
  if ! pgrep -x "redis-server" > /dev/null; then
    redis-server /etc/redis/redis.conf &
    echo "Redis started."
  else
    echo "Redis is already running."
  fi
}

redis_stop() {
  if pgrep -x "redis-server" > /dev/null; then
    pkill -x "redis-server"
    echo "Redis stopped."
  else
    echo "Redis is not running."
  fi
}

redis_restart() {
  redis_stop
  sleep 1
  redis_start
}

redis_status() {
  if pgrep -x "redis-server" > /dev/null; then
    echo "Redis is running."
    echo "Redis process information:"
    ps aux | grep "redis-server" | grep -v "grep"
  else
    echo "Redis is not running."
  fi
}

mysql_start() {
  if ! pgrep -x "mysqld" > /dev/null; then
    sudo service mysql start
    echo "MySQL started."
  else
    echo "MySQL is already running."
  fi
}

mysql_stop() {
  if pgrep -x "mysqld" > /dev/null; then
    sudo service mysql stop
    echo "MySQL stopped."
  else
    echo "MySQL is not running."
  fi
}

mysql_restart() {
  mysql_stop
  sleep 1
  mysql_start
}

mysql_status() {
  if pgrep -x "mysqld" > /dev/null; then
    echo "MySQL is running."
    echo "MySQL process information:"
    ps aux | grep "mysqld" | grep -v "grep"
  else
    echo "MySQL is not running."
  fi
}

services_start() {
  redis_start
  mysql_start
}

services_stop() {
  redis_stop
  mysql_stop
}

services_restart() {
  redis_restart
  mysql_restart
}

services_status() {
  redis_status
  mysql_status
}

check_ports() {
  echo "Redis listening ports:"
  sudo ss -tulpn | grep "redis-server"
  echo "MySQL listening ports:"
  sudo ss -tulpn | grep "mysqld"
}

case "$1" in
  start)
    case "$2" in
      redis)
        redis_start
        ;;
      mysql)
        mysql_start
        ;;
      *)
        services_start
        ;;
    esac
    ;;
  stop)
    case "$2" in
      redis)
        redis_stop
        ;;
      mysql)
        mysql_stop
        ;;
      *)
        services_stop
        ;;
    esac
    ;;
  restart)
    case "$2" in
      redis)
        redis_restart
        ;;
      mysql)
        mysql_restart
        ;;
      *)
        services_restart
        ;;
    esac
    ;;
  status)
    case "$2" in
      redis)
        redis_status
        ;;
      mysql)
        mysql_status
        ;;
      *)
        services_status
        ;;
    esac
    ;;
  ports)
    check_ports
    ;;
  *)
    echo "Usage: $0

```

```sh
# Alias to manage Redis and MySQL services
alias service.start="~/services.sh start"
alias service.stop="~/services.sh stop"
alias service.restart="~/services.sh restart"
alias service.status="~/services.sh status"
alias service.ports="~/services.sh ports"

alias service.redis.start="~/services.sh start redis"
alias service.redis.stop="~/services.sh stop redis"
alias service.redis.restart="~/services.sh restart redis"
alias service.redis.status="~/services.sh status redis"

alias service.mysql.start="~/services.sh start mysql"
alias service.mysql.stop="~/services.sh stop mysql"
alias service.mysql.restart="~/services.sh restart mysql"
alias service.mysql.status="~/services.sh status mysql"
```

