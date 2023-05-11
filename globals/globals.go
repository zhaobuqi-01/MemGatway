package globals

import (
	"gateway/flow_counter"
	"gateway/mq"
	"gateway/pkg/database/redis"
	"sync"
)

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

	FlowTotal = "flow_total"

	JwtSignKey = "my_sign_key"
	JwtExpires = 60 * 60
)

var (
	LoadTypeMap = map[int]string{
		LoadTypeHTTP: "HTTP",
		LoadTypeTCP:  "TCP",
		LoadTypeGRPC: "GRPC",
	}
)

type DataChangeMessage struct {
	Type        string `json:"type"`
	Payload     string `json:"payload"`
	ServiceType int    `json:"service_type"`
	Operation   string `json:"operation"`
}

const (
	DataDelete = "delete"
	DataUpdate = "update"
	DataInsert = "insert"
)

var (
	MessageQueue mq.MQ
	FlowCounter  flow_counter.FlowCounter
	once         sync.Once
)

func Init() {
	once.Do(func() {
		MessageQueue = mq.Default(redis.GetRedisConnection())
		FlowCounter = flow_counter.NewFlowCounter()
	})
}
