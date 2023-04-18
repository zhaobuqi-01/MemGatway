HTTP 状态码是服务器在响应客户端请求时返回的 3 位数字，用于表示请求的处理结果。状态码分为五类，根据第一个数字进行分类：

1. 1xx (信息性状态码)：表示请求已被接收，需要继续处理。
2. 2xx (成功状态码)：表示请求已成功处理。
3. 3xx (重定向状态码)：表示需要进一步操作以完成请求。
4. 4xx (客户端错误状态码)：表示请求包含无效数据或无法完成。
5. 5xx (服务器错误状态码)：表示服务器在处理请求时发生错误。

以下是一些常用的 HTTP 状态码：

1xx 信息性状态码：

- 100 Continue：表示客户端应继续发送请求。

2xx 成功状态码：

- 200 OK：请求成功处理。
- 201 Created：请求成功处理，并创建了新的资源。
- 202 Accepted：请求已接受，但尚未处理。
- 204 No Content：请求成功处理，但没有返回任何内容。

3xx 重定向状态码：

- 300 Multiple Choices：表示多个资源可供选择。
- 301 Moved Permanently：请求的资源已被永久移动到新位置。
- 302 Found：请求的资源临时存在不同的 URI。
- 304 Not Modified：资源未修改，使用缓存的版本。

4xx 客户端错误状态码：

- 400 Bad Request：请求格式错误或包含无效数据。
- 401 Unauthorized：请求需要身份验证。
- 403 Forbidden：客户端没有访问资源的权限。
- 404 Not Found：请求的资源在服务器上不存在。
- 405 Method Not Allowed：不允许使用请求行中指定的方法。
- 429 Too Many Requests：客户端发送的请求过多，触发了限速。

5xx 服务器错误状态码：

- 500 Internal Server Error：服务器在处理请求时发生错误。
- 501 Not Implemented：服务器不支持实现请求所需的功能。
- 502 Bad Gateway：上游服务器返回了无效响应。
- 503 Service Unavailable：服务器暂时无法处理请求（可能由于过载或维护）。
- 504 Gateway Timeout：上游服务器未及时响应。

[http状态码](https://developer.mozilla.org/en-US/docs/Web/HTTP/Status#client_error_responses)

const (
	StatusContinue           = 100 // RFC 7231, 6.2.1
	StatusSwitchingProtocols = 101 // RFC 7231, 6.2.2
	StatusProcessing         = 102 // RFC 2518, 10.1

```go
StatusOK                   = 200 // RFC 7231, 6.3.1
StatusCreated              = 201 // RFC 7231, 6.3.2
StatusAccepted             = 202 // RFC 7231, 6.3.3
StatusNonAuthoritativeInfo = 203 // RFC 7231, 6.3.4
StatusNoContent            = 204 // RFC 7231, 6.3.5
StatusResetContent         = 205 // RFC 7231, 6.3.6
StatusPartialContent       = 206 // RFC 7233, 4.1
StatusMultiStatus          = 207 // RFC 4918, 11.1
StatusAlreadyReported      = 208 // RFC 5842, 7.1
StatusIMUsed               = 226 // RFC 3229, 10.4.1

StatusMultipleChoices   = 300 // RFC 7231, 6.4.1
StatusMovedPermanently  = 301 // RFC 7231, 6.4.2
StatusFound             = 302 // RFC 7231, 6.4.3
StatusSeeOther          = 303 // RFC 7231, 6.4.4
StatusNotModified       = 304 // RFC 7232, 4.1
StatusUseProxy          = 305 // RFC 7231, 6.4.5
StatusTemporaryRedirect = 307 // RFC 7231, 6.4.7
StatusPermanentRedirect = 308 // RFC 7538, 3

StatusBadRequest                   = 400 // RFC 7231, 6.5.1
StatusUnauthorized                 = 401 // RFC 7235, 3.1
StatusPaymentRequired              = 402 // RFC 7231, 6.5.2
StatusForbidden                    = 403 // RFC 7231, 6.5.3
StatusNotFound                     = 404 // RFC 7231, 6.5.4
StatusMethodNotAllowed             = 405 // RFC 7231, 6.5.5
StatusNotAcceptable                = 406 // RFC 7231, 6.5.6
StatusProxyAuthRequired            = 407 // RFC 7235, 3.2
StatusRequestTimeout               = 408 // RFC 7231, 6.5.7
StatusConflict                     = 409 // RFC 7231, 6.5.8
StatusGone                         = 410 // RFC 7231, 6.5.9
StatusLengthRequired               = 411 // RFC 7231, 6.5.10
StatusPreconditionFailed           = 412 // RFC 7232, 4.2
StatusRequestEntityTooLarge        = 413 // RFC 7231, 6.5.11
StatusRequestURITooLong            = 414 // RFC 7231, 6.5.12
StatusUnsupportedMediaType         = 415 // RFC const (
	StatusContinue           = 100 // RFC 7231, 6.2.1
	StatusSwitchingProtocols = 101 // RFC 7231, 6.2.2
	StatusProcessing         = 102 // RFC 2518, 10.1

	StatusOK                   = 200 // RFC 7231, 6.3.1
	StatusCreated              = 201 // RFC 7231, 6.3.2
	StatusAccepted             = 202 // RFC 7231, 6.3.3
	StatusNonAuthoritativeInfo = 203 // RFC 7231, 6.3.4
	StatusNoContent            = 204 // RFC 7231, 6.3.5
	StatusResetContent         = 205 // RFC 7231, 6.3.6
	StatusPartialContent       = 206 // RFC 7233, 4.1
	StatusMultiStatus          = 207 // RFC 4918, 11.1
	StatusAlreadyReported      = 208 // RFC 5842, 7.1
	StatusIMUsed               = 226 // RFC 3229, 10.4.1

	StatusMultipleChoices   = 300 // RFC 7231, 6.4.1
	StatusMovedPermanently  = 301 // RFC 7231, 6.4.2
	StatusFound             = 302 // RFC 7231, 6.4.3
	StatusSeeOther          = 303 // RFC 7231, 6.4.4
	StatusNotModified       = 304 // RFC 7232, 4.1
	StatusUseProxy          = 305 // RFC 7231, 6.4.5
	StatusTemporaryRedirect = 307 // RFC 7231, 6.4.7
	StatusPermanentRedirect = 308 // RFC 7538, 3

	StatusBadRequest                   = 400 // RFC 7231, 6.5.1
	StatusUnauthorized                 = 401 // RFC 7235, 3.1
	StatusPaymentRequired              = 402 // RFC 7231, 6.5.2
	StatusForbidden                    = 403 // RFC 7231, 6.5.3
	StatusNotFound                     = 404 // RFC 7231, 6.5.4
	StatusMethodNotAllowed             = 405 // RFC 7231, 6.5.5
	StatusNotAcceptable                = 406 // RFC 7231, 6.5.6
	StatusProxyAuthRequired            = 407 // RFC 7235, 3.2
	StatusRequestTimeout               = 408 // RFC 7231, 6.5.7
	StatusConflict                     = 409 // RFC 7231, 6.5.8
	StatusGone                         = 410 // RFC 7231, 6.5.9
	StatusLengthRequired               = 411 // RFC 7231, 6.5.10
	StatusPreconditionFailed           = 412 // RFC 7232, 4.2
	StatusRequestEntityTooLarge        = 413 // RFC 7231, 6.5.11
	StatusRequestURITooLong            = 414 // RFC 7231, 6.5.12
	StatusUnsupportedMediaType         = 415 // RFC 7231, 6.5.13
   StatusRequestedRangeNotSatisfiable = 416 // RFC 7233, 4.4
	StatusExpectationFailed = 417 // RFC 7231, 6.5.14
	StatusTeapot = 418 // RFC 7168, 2.3.3
	StatusMisdirectedRequest = 421 // RFC 7540, 9.1.2
    StatusUnprocessableEntity = 422 // RFC 4918, 11.2
    StatusLocked = 423 // RFC 4918, 11.3
    StatusFailedDependency = 424 // RFC 4918, 11.4
    StatusUpgradeRequired = 426 // RFC 7231, 6.5.15
    StatusPreconditionRequired = 428 // RFC 6585, 3
    StatusTooManyRequests = 429 // RFC 6585, 4
    StatusRequestHeaderFieldsTooLarge = 431 // RFC 6585, 5
    StatusUnavailableForLegalReasons = 451 // RFC 7725, 3
	
    StatusInternalServerError           = 500 // RFC 7231, 6.6.1
    StatusNotImplemented                = 501 // RFC 7231, 6.6.2
    StatusBadGateway                    = 502 // RFC 7231, 6.6.3
    StatusServiceUnavailable            = 503 // RFC 7231, 6.6.4
    StatusGatewayTimeout                = 504 // RFC 7231, 6.6.5
    StatusHTTPVersionNotSupported       = 505 // RFC 7231, 6.6.6
    StatusVariantAlsoNegotiates         = 506 // RFC 2295, 8.1
    StatusInsufficientStorage           = 507 // RFC 4918, 11.5
    StatusLoopDetected                  = 508 // RFC 5842, 7.2
    StatusNotExtended                   = 510 // RFC 2774, 7
    StatusNetworkAuthenticationRequired = 511 // RFC 6585, 6

```

