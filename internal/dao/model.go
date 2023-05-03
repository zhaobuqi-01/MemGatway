package dao

// Model is an interface representing various types of database models.
// It includes Admin, ServiceInfo, AccessControl, GrpcRule,
// HttpRule, TcpRule, LoadBalance, and App.
type Model interface {
	Admin | ServiceInfo | AccessControl | GrpcRule |
		HttpRule | TcpRule | LoadBalance | App
}
