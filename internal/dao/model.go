package dao

type Model interface {
	Admin | ServiceInfo | AccessControl | GrpcRule |
		HttpRule | TcpRule | LoadBalance | App
}
