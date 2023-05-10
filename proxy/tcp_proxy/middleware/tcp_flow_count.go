package middleware

import (
	"gateway/enity"
	"gateway/globals"
)

func TCPFlowCountMiddleware() func(c *TcpSliceRouterContext) {
	return func(c *TcpSliceRouterContext) {
		serverInterface := c.Get("service")
		if serverInterface == nil {
			c.conn.Write([]byte("get service empty"))
			c.Abort()
			return
		}

		serviceDetail := serverInterface.(*enity.ServiceDetail)

		totalCounter, err := globals.FlowCounter.GetCounter("flow_total")
		if err != nil {
			c.conn.Write([]byte(err.Error()))
			c.Abort()
			return
		}
		totalCounter.Increase()
		serviceCounter, err := globals.FlowCounter.GetCounter(serviceDetail.Info.ServiceName)
		if err != nil {
			c.conn.Write([]byte(err.Error()))
			c.Abort()
			return
		}
		serviceCounter.Increase()
		c.Next()
	}
}
