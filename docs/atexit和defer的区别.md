# atexit和defer的区别

`atexit` 和 `defer` 都是在 Go 语言中用于在程序执行完某些操作后执行特定代码的方法。但它们在实现方式和使用场景上有一些不同：

1. `defer`：
   - `defer` 是 Go 语言的关键字，它用于将函数推迟执行，直到包含它的函数返回。
   - 当有多个 `defer` 语句时，它们将以后进先出（LIFO）的顺序执行。
   - `defer` 通常用于确保资源（例如文件或网络连接）在函数执行完毕后得到释放，或者用于在函数返回前执行一些清理工作。
   - `defer` 只在当前函数的范围内有效。
2. `atexit`：
   - `atexit` 是一个第三方库，它提供了在程序退出时执行特定函数的功能。它类似于 C 语言中的 `atexit()` 函数。
   - `atexit` 注册的函数将在程序正常退出时执行。当有多个函数被注册时，它们将以相反的顺序（即最后注册的函数将首先执行）执行。
   - `atexit` 可以在整个程序范围内使用，通常用于在程序退出时执行全局的清理任务，例如关闭数据库连接或释放全局资源。
   - 要使用 `atexit`，需要导入 "github.com/tebeka/atexit" 包。

总结：`defer` 是 Go 语言内建的延迟执行功能，通常用于在函数范围内确保资源得到释放或执行清理工作；而 `atexit` 是一个第三方库，用于在程序退出时执行全局的清理任务。根据具体场景和需求选择使用哪种方法。