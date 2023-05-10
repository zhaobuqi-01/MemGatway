package middleware

import (
	"fmt"
	"gateway/enity"
	"gateway/globals"
	"gateway/pkg/response"

	"github.com/gin-gonic/gin"
)

func HTTPTrafficStats() gin.HandlerFunc {
	return func(c *gin.Context) {

		serverInterface, ok := c.Get("service")
		if !ok {
			response.ResponseError(c, response.ServiceNotFoundErrCode, fmt.Errorf("service not found"))
			c.Abort()
			return
		}

		serviceDetail := serverInterface.(*enity.ServiceDetail)

		totalCounter, err := globals.FlowCounter.GetCounter(globals.FlowTotal)
		if err != nil {
			response.ResponseError(c, response.CommErrCode, err)
			c.Abort()
			return
		}
		totalCounter.Increase()
		serviceCounter, err := globals.FlowCounter.GetCounter(serviceDetail.Info.ServiceName)
		if err != nil {
			response.ResponseError(c, response.CommErrCode, err)
			c.Abort()
			return
		}
		serviceCounter.Increase()
		c.Next()
	}
}
