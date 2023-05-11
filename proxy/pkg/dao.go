package pkg

import (
	"fmt"
	"gateway/enity"
	"gateway/globals"
	"gateway/pkg/log"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Model interface {
	enity.Admin | enity.ServiceInfo | enity.AccessControl | enity.GrpcRule |
		enity.HttpRule | enity.TcpRule | enity.LoadBalance | enity.App
}

// PageList 分页查询
func getAll[T Model](db *gorm.DB, queryConditions []func(db *gorm.DB) *gorm.DB) ([]T, error) {
	// log记录查询信息
	log.Info("start pageList ")

	list := []T{}

	query := db.Where("is_delete=0")
	for _, condition := range queryConditions {
		query = condition(query)
	}
	if err := query.Order("id desc").Find(&list).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Error(fmt.Sprintf("error retrieving :%v ", query), zap.Error(err))
		return nil, err
	}

	// log记录成功信息
	log.Info("pageList successfully", zap.Any("list", list))
	return list, nil
}

// GetServiceDetail 获取服务详情
// func getServiceDetail(db *gorm.DB, search *enity.ServiceInfo) (*enity.ServiceDetail, error) {
// 	// log记录查询信息
// 	log.Info("start getting service detail")

// 	if search.ServiceName == "" {
// 		info, err := get(db, search)
// 		if err != nil {
// 			return nil, err
// 		}
// 		search = info
// 	}

// 	httpRule, err := get(db, &enity.HttpRule{ServiceID: search.ID})
// 	if err != nil && err != gorm.ErrRecordNotFound {
// 		log.Error("error retrieving http rule", zap.Error(err))
// 		return nil, err
// 	}
// 	log.Info("get http rule successful", zap.Any("httpRule", httpRule))

// 	tcpRule, err := get(db, &enity.TcpRule{ServiceID: search.ID})
// 	if err != nil && err != gorm.ErrRecordNotFound {
// 		log.Error("error retrieving tcp rule", zap.Error(err))
// 		return nil, err
// 	}
// 	log.Info("get tcp rule successful", zap.Any("tcpRule", tcpRule))

// 	grpcRule, err := get(db, &enity.GrpcRule{ServiceID: search.ID})
// 	if err != nil && err != gorm.ErrRecordNotFound {
// 		log.Error("error retrieving grpc rule", zap.Error(err))
// 		return nil, err
// 	}
// 	log.Info("get grpc rule successful", zap.Any("grpcRule", grpcRule))

// 	accessControl, err := get(db, &enity.AccessControl{ServiceID: search.ID})
// 	if err != nil && err != gorm.ErrRecordNotFound {
// 		log.Error("error retrieving access control", zap.Error(err))
// 		return nil, err
// 	}
// 	log.Info("get access control successful", zap.Any("accessControl", accessControl))

// 	loadBalance, err := get(db, &enity.LoadBalance{ServiceID: search.ID})
// 	if err != nil && err != gorm.ErrRecordNotFound {
// 		log.Error("error retrieving load balance", zap.Error(err))
// 		return nil, err
// 	}
// 	log.Info("get load balance successfully", zap.Any("loadBalance", loadBalance))

// 	detail := &enity.ServiceDetail{
// 		Info:          search,
// 		HTTPRule:      httpRule,
// 		TCPRule:       tcpRule,
// 		GRPCRule:      grpcRule,
// 		LoadBalance:   loadBalance,
// 		AccessControl: accessControl,
// 	}

// 	// log记录成功取到信息
// 	log.Info("get service detail successful", zap.Any("detail", detail))
// 	return detail, nil
// }

// Get查询单条数据
func get[T Model](db *gorm.DB, search *T) (*T, error) {
	// log记录查询信息
	log.Info("start getting", zap.Any("search", search))

	var out T
	result := db.Where(search).First(&out)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			log.Error("record not found", zap.Any("search", search))
			return nil, result.Error
		}

		log.Error(fmt.Sprintf("error retrieving :%v ", search), zap.Error(result.Error))
		return nil, result.Error
	}

	log.Info("got successfully", zap.Any("search", search))
	return &out, nil
}

func getServiceDetail(db *gorm.DB, search *enity.ServiceInfo) (*enity.ServiceDetail, error) {
	// log记录查询信息
	log.Info("start getting service detail")

	if search.ServiceName == "" {
		info, err := get(db, search)
		if err != nil {
			return nil, err
		}
		search = info
	}

	var httpRule, tcpRule, grpcRule interface{}
	var err error

	// 优化后的查询代码
	switch search.LoadType {
	case globals.LoadTypeHTTP:
		httpRule, err = get(db, &enity.HttpRule{ServiceID: search.ID})
		if err != nil && err != gorm.ErrRecordNotFound {
			log.Error("error retrieving http rule", zap.Error(err))
			return nil, err
		}
		log.Info("get http rule successful", zap.Any("httpRule", httpRule))
	case globals.LoadTypeTCP:
		tcpRule, err = get(db, &enity.TcpRule{ServiceID: search.ID})
		if err != nil && err != gorm.ErrRecordNotFound {
			log.Error("error retrieving tcp rule", zap.Error(err))
			return nil, err
		}
		log.Info("get tcp rule successful", zap.Any("tcpRule", tcpRule))
	case globals.LoadTypeGRPC:
		grpcRule, err = get(db, &enity.GrpcRule{ServiceID: search.ID})
		if err != nil && err != gorm.ErrRecordNotFound {
			log.Error("error retrieving grpc rule", zap.Error(err))
			return nil, err
		}
		log.Info("get grpc rule successful", zap.Any("grpcRule", grpcRule))
	}

	accessControl, err := get(db, &enity.AccessControl{ServiceID: search.ID})
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Error("error retrieving access control", zap.Error(err))
		return nil, err
	}
	log.Info("get access control successful", zap.Any("accessControl", accessControl))

	loadBalance, err := get(db, &enity.LoadBalance{ServiceID: search.ID})
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Error("error retrieving load balance", zap.Error(err))
		return nil, err
	}
	log.Info("get load balance successfully", zap.Any("loadBalance", loadBalance))

	detail := &enity.ServiceDetail{
		Info:          search,
		LoadBalance:   loadBalance,
		AccessControl: accessControl,
	}

	if httpRule != nil {
		if rule, ok := httpRule.(*enity.HttpRule); ok {
			detail.HTTPRule = rule
		} else {
			log.Error("error retrieving http rule: unexpected type", zap.Any("tcpRule", tcpRule))
			return nil, fmt.Errorf("unexpected type for TCP rule: %T", tcpRule)
		}
	}
	if grpcRule != nil {
		if rule, ok := grpcRule.(*enity.GrpcRule); ok {
			detail.GRPCRule = rule
		} else {
			log.Error("error retrieving grpc rule: unexpected type", zap.Any("tcpRule", tcpRule))
			return nil, fmt.Errorf("unexpected type for TCP rule: %T", tcpRule)
		}
	}
	if tcpRule != nil {
		if rule, ok := tcpRule.(*enity.TcpRule); ok {
			detail.TCPRule = rule
		} else {
			log.Error("error retrieving tcp rule: unexpected type", zap.Any("tcpRule", tcpRule))
			return nil, fmt.Errorf("unexpected type for TCP rule: %T", tcpRule)
		}
	}

	// log记录成功取到信息
	log.Info("get service detail successful", zap.Any("detail", detail))
	return detail, nil
}
