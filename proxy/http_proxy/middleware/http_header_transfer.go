package middleware

import (
	"fmt"
	"gateway/enity"
	"gateway/pkg/response"
	"strings"

	"github.com/gin-gonic/gin"
)

// HTTPHeaderTransferMiddleware is a Gin middleware function that modifies
// the HTTP headers of a request according to the specified transformation rules.
// The transformation rules are retrieved from the service detail's HTTPRule.
// Supported operations are "add", "edit", and "del".
//
// The function returns a Gin HandlerFunc that can be used as a middleware
// in a Gin HTTP server.
func HTTPHeaderTransferMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the service detail from the context
		serverInterface, ok := c.Get("service")
		if !ok {
			response.ResponseError(c, response.ServiceNotFoundErrCode, fmt.Errorf("service not found"))
			c.Abort()
			return
		}
		serviceDetail := serverInterface.(*enity.ServiceDetail)

		// Split the header transformation rules and apply them to the request
		for _, item := range strings.Split(serviceDetail.HTTPRule.HeaderTransfor, ",") {
			items := strings.Split(item, " ")
			if len(items) != 3 {
				continue
			}

			operation, headerKey, headerValue := items[0], items[1], items[2]

			switch operation {
			case "add", "edit":
				c.Request.Header.Set(headerKey, headerValue)
			case "del":
				c.Request.Header.Del(headerKey)
			}
		}

		c.Next()
	}
}
