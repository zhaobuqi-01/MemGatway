# go doc

`go doc` 是 Go 语言自带的文档工具，可以方便地查看和生成代码文档。

可以通过以下命令使用 `go doc` 查看包、函数、类型等的文档：

```
go doc <package>
go doc <package>.<function>
go doc <package>.<type>
```

其中 `<package>` 表示包的导入路径，`<function>` 和 `<type>` 分别表示函数和类型的全名。

例如，要查看 `fmt` 包的文档，可以执行以下命令：

```
go doc fmt
```

如果要查看 `fmt` 包中的 `Println` 函数的文档，可以执行以下命令：

```
go doc fmt.Println
```

如果要查看 `fmt` 包中的 `Stringer` 接口的文档，可以执行以下命令：

```
go doc fmt.Stringer
```

此外，`go doc` 还可以用于生成文档。在要生成文档的包目录下执行以下命令即可生成文档：

```
go doc -all -html > doc.html
```

其中 `-all` 参数表示生成所有内容的文档，`-html` 表示将生成的文档转换为 HTML 格式。执行后会生成一个 `doc.html` 文件，可以用浏览器打开查看文档。

另外，还可以使用第三方工具 `godoc` 来生成文档，使用方法类似于 `go doc`，不过更加方便易用。执行以下命令即可启动 `godoc` 服务：

```
godoc -http=:6060
```

然后在浏览器中访问 `http://localhost:6060` 即可查看文档。

`godocdown`是一个可以将`godoc`生成的HTML格式的文档转换为Markdown格式的工具。

首先，你需要安装`godocdown`：

```
go get github.com/robertkrimen/godocdown/godocdown
```

然后，你可以在你的包的目录下运行`godocdown`：

```
godocdown > README.md
```

这将会生成一个名为`README.md`的Markdown文件，里面包含了你的包的文档。

## 注释

为了让 `go doc` 正确地显示文档，需要在函数、方法、结构体、接口和包等的注释前添加特殊的注释，格式如下：

- 对于函数、方法、结构体、接口等：在声明前，使用 

  ```
  //
  ```

   或 

  ```
  /* */
  ```

   注释块注释以下信息：

  - 描述信息
  - 函数、方法的参数和返回值
  - 函数、方法的行为和副作用
  - 函数、方法的例子和用法

- 对于包注释：在包目录下添加一个名为 

  ```
  doc.go
  ```

   的文件，并在文件中使用 

  ```
  /* */
  ```

   注释块注释以下信息：

  - 描述信息
  - 包级别的变量、常量、函数、类型和接口等的声明和文档