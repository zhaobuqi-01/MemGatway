package pkg

const (
	ValidatorKey               = "ValidatorKey"
	TranslatorKey              = "TranslatorKey"
	AdminSessionInfoKey string = "AdminSessionInfoKey"
	LoadTypeHTTP               = 0
	LoadTypeTCP                = 1
	LoadTypeGRPC               = 2

	HTTPRuleTypePrefixURL = 0
	HTTPRuleTypeDomain    = 1

	FlowTotal         = "flow_total"
	FlowServicePrefix = "flow_service_"
	FlowAppPrefix     = "flow_app_"
)

var (
	LoadTypeMap = map[int]string{
		LoadTypeHTTP: "HTTP",
		LoadTypeTCP:  "TCP",
		LoadTypeGRPC: "GRPC",
	}
)
