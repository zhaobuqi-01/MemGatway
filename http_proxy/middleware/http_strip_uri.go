package middleware

import (
	"fmt"
	"gateway/dao"
	"gateway/utils"

	"strings"

	"github.com/gin-gonic/gin"
)

// HTTPStripUriMiddleware is a Gin middleware function that conditionally
// strips the prefix from the request's URL path based on the service detail's
// HTTPRule. The prefix is removed if the rule type is HTTPRuleTypePrefixURL
// and NeedStripUri is set to 1.
//
// The function returns a Gin HandlerFunc that can be used as a middleware
// in a Gin HTTP server.
func HTTPStripUriMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the service detail from the context
		serverInterface, ok := c.Get("service")
		if !ok {
			utils.ResponseError(c, utils.ServiceNotFoundErrCode, fmt.Errorf("service not found"))
			c.Abort()
			return
		}
		serviceDetail := serverInterface.(*dao.ServiceDetail)

		// Check if the URL path should be stripped
		if serviceDetail.HTTPRule.RuleType == utils.HTTPRuleTypePrefixURL && serviceDetail.HTTPRule.NeedStripUri == 1 {
			// Remove the prefix from the request's URL path
			c.Request.URL.Path = strings.Replace(c.Request.URL.Path, serviceDetail.HTTPRule.Rule, "", 1)
		}

		c.Next()
	}
}
