# zookeeper

## zookeeper使用

#### 1.安装`go-zookeeper/zk`

```go
go get github.com/go-zookeeper/zk
```

#### 2.连接到zookeeper

```go
package main

import (
	"fmt"
	"time"
	"github.com/go-zookeeper/zk"
)

func main() {
	servers := []string{"localhost:2181"}
	conn, _, err := zk.Connect(servers, time.Second)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// 与Zookeeper交互的代码
}
```

#### 3.创建节点（ZNode）：使用`conn.Create()`方法创建一个ZNode

```go
path := "/myapp"
data := []byte("some data")
flags := int32(0)
acl := zk.WorldACL(zk.PermAll)

createdPath, err := conn.Create(path, data, flags, acl)
if err != nil {
	panic(err)
}
fmt.Printf("Created ZNode: %s\n", createdPath)
```

#### 4.读取ZNode：使用`conn.Get()`方法读取ZNode的数据

```go
data, stat, err := conn.Get(path)
if err != nil {
	panic(err)
}
fmt.Printf("Data: %s, Version: %d\n", string(data), stat.Version)
```

#### 5.更新ZNode：使用`conn.Set()`方法更新ZNode的数据

```go
newData := []byte("updated data")
stat, err = conn.Set(path, newData, stat.Version)
if err != nil {
	panic(err)
}
fmt.Printf("Updated ZNode, New Version: %d\n", stat.Version)
```

#### 6.删除ZNode：使用`conn.Delete()`方法删除ZNode

```go
err = conn.Delete(path, stat.Version)
if err != nil {
	panic(err)
}
fmt.Println("Deleted ZNode")
```

#### 7.监听ZNode变化：可以使用`conn.GetW()`方法在读取ZNode数据的同时设置一个监听器，当ZNode发生变化时会收到通知

```go
data, stat, watch, err := conn.GetW(path)
if err != nil {
	panic(err)
}

go func() {
	event := <-watch
	fmt.Printf("Event Type: %s, Path: %s\n", event.Type, event.Path)
}()

// 模拟其他客户端更新ZNode
conn.Set(path, []byte("new data"), stat.Version)

// 等待事件通知
time.Sleep(2 * time.Second)
```

