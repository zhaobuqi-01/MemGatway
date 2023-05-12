package dao

import "gateway/enity"

// Model is an interface representing various types of database models.
// It includes Admin, ServiceInfo, AccessControl, GrpcRule,
// HttpRule, TcpRule, LoadBalance, and App.
type Model interface {
	enity.Admin | enity.ServiceInfo | enity.AccessControl | enity.GrpcRule |
		enity.HttpRule | enity.TcpRule | enity.LoadBalance | enity.App
}
