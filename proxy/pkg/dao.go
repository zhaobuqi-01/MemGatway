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
func pageList[T Model](db *gorm.DB, queryConditions []func(db *gorm.DB) *gorm.DB, PageNo, PageSize int) ([]T, int64, error) {
	// log记录查询信息
	log.Info("start pageList ")

	total := int64(0)
	list := []T{}
	offset := (PageNo - 1) * PageSize

	query := db.Where("is_delete=0")
	for _, condition := range queryConditions {
		query = condition(query)
	}
	if err := query.Limit(PageSize).Offset(offset).Order("id desc").Find(&list).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Error(fmt.Sprintf("error retrieving :%v ", query), zap.Error(err))
		return nil, 0, err
	}
	query.Limit(PageSize).Offset(offset).Count(&total)

	// log记录成功信息
	log.Info("pageList successfully", zap.Any("list", list))
	return list, total, nil
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
		detail.HTTPRule = httpRule.(*enity.HttpRule)
	}
	if tcpRule != nil {
		detail.TCPRule = tcpRule.(*enity.TcpRule)
	}
	if grpcRule != nil {
		detail.GRPCRule = grpcRule.(*enity.GrpcRule)
	}

	// log记录成功取到信息
	log.Info("get service detail successful", zap.Any("detail", detail))
	return detail, nil
}
