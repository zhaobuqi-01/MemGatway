package globals

const (
	LoadTypeHTTP = iota
	LoadTypeTCP
	LoadTypeGRPC

	HTTPRuleTypePrefixURL = 0
	HTTPRuleTypeDomain    = 1

	ValidatorKey               = "ValidatorKey"
	TranslatorKey              = "TranslatorKey"
	AdminSessionInfoKey string = "AdminSessionInfoKey"

	DataChange = "data_change"
)

var (
	LoadTypeMap = map[int]string{
		LoadTypeHTTP: "HTTP",
		LoadTypeTCP:  "TCP",
		LoadTypeGRPC: "GRPC",
	}
)

type DataChangeMessage struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}
