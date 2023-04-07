# swaggger

### 命令

```shell
swag init --generalInfo cmd/server/main.go --output api/v1
```

### 网址

```
 http://localhost:8080/swagger/index.html 
```

### 注解

1. @swagger.Info： 用于描述 API 文档的基本信息。这个注解通常只需要使用一次，位于文档的开头。它可以包含以下字段：

- @title：API 文档的标题。
- @description：API 文档的描述。
- @version：API 文档的版本。
- @termsOfService：API 使用的服务条款。
- @contact.name：联系人的名称。
- @contact.email：联系人的电子邮件。
- @contact.url：联系人的网址。
- @license.name：API 使用的许可证名称。
- @license.url：API 使用的许可证的网址。

1. @swagger.Schemes： 定义 API 支持的传输协议，如 HTTP 和 HTTPS。
2. @swagger.BasePath： 定义 API 的基本路径。所有 API 接口的路径都会以此路径为前缀。
3. @swagger.Tags： 定义 API 的分类标签。这些标签可以用于将相关的 API 接口分组到一起。
4. @swagger.Path： 定义 API 的具体路径和操作，包括 HTTP 方法、请求参数、响应结果等。这是最重要的注解，用于描述 API 接口的具体信息。以下是一些与 @swagger.Path 相关的注解：

- @tags：API 接口所属的分类标签。
- @summary：API 接口的简短描述。
- @description：API 接口的详细描述。
- @accept：API 接口接受的请求数据类型，如 JSON、XML 等。
- @produce：API 接口生成的响应数据类型，如 JSON、XML 等。
- @param：API 接口的请求参数。可以包含参数的名称、类型、是否必需等信息。
- @success：API 接口成功响应的结果。可以包含 HTTP 状态码、返回数据类型、返回数据的结构等信息。
- @failure：API 接口失败响应的结果。可以包含 HTTP 状态码、错误信息等。
- @router：API 接口的路径和 HTTP 方法。这个注解用于生成 API 文档中的实际 URL。