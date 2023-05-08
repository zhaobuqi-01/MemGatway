package middleware

import (
	"fmt"
	"gateway/enity"
	"gateway/pkg/response"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

// HTTPUrlRewriteMiddleware is a Gin middleware function that rewrites the
// request's URL path based on the service detail's HTTPRule. The rewrite
// rules are defined as a comma-separated list of regex patterns and
// replacement strings in the UrlRewrite field.
//
// The function returns a Gin HandlerFunc that can be used as a middleware
// in a Gin HTTP server.
//
// Example:
//
//	If the UrlRewrite is "/old-path/(.*) /new-path/$1", a request with the
//	URL path "/old-path/test" will be rewritten as "/new-path/test".
//
// Input:
//   - A Gin context containing the request
//
// Output:
//   - The request's URL path is updated according to the rewrite rules
func HTTPUrlRewriteMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the service detail from the context
		serverInterface, ok := c.Get("service")
		if !ok {
			response.ResponseError(c, response.ServiceNotFoundErrCode, fmt.Errorf("service not found"))
			c.Abort()
			return
		}
		serviceDetail := serverInterface.(*enity.ServiceDetail)

		// Apply the URL rewrite rules
		for _, item := range strings.Split(serviceDetail.HTTPRule.UrlRewrite, ",") {
			items := strings.Split(item, " ")
			if len(items) != 2 {
				continue
			}

			pattern, replacement := items[0], items[1]

			regexp, err := regexp.Compile(pattern)
			if err != nil {
				continue
			}

			rewrittenPath := regexp.ReplaceAll([]byte(c.Request.URL.Path), []byte(replacement))
			c.Request.URL.Path = string(rewrittenPath)
		}

		c.Next()
	}
}
