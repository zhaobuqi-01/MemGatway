package dao

import (
	"fmt"
	"gateway/backend/dto"
	"gateway/enity"
	"gateway/globals"
	"gateway/pkg/log"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type APP interface {
	Getter[enity.App]
	Saver[enity.App]
	Deleter[enity.App]
	PagedLister[enity.App]
	AllGetter[enity.App]
}

func NewApp() APP {
	return New[enity.App]()
}

type Admin interface {
	Getter[enity.Admin]
	Saver[enity.Admin]
}

func NewAdmin() Admin {
	return New[enity.Admin]()
}

type TcpService interface {
	Getter[enity.TcpRule]
	Saver[enity.TcpRule]
	Deleter[enity.TcpRule]
}

func NewTcpService() TcpService {
	return New[enity.TcpRule]()
}

type GrpcService interface {
	Getter[enity.GrpcRule]
	Saver[enity.GrpcRule]
	Deleter[enity.GrpcRule]
}

func NewGrpcService() GrpcService {
	return New[enity.GrpcRule]()
}

type HttpService interface {
	Getter[enity.HttpRule]
	Saver[enity.HttpRule]
	Deleter[enity.HttpRule]
}

func NewHttpService() HttpService {
	return New[enity.HttpRule]()
}

type LoadBalanceService interface {
	Getter[enity.LoadBalance]
	Saver[enity.LoadBalance]
	Deleter[enity.LoadBalance]
}

func NewLoadBalanceService() LoadBalanceService {
	return New[enity.LoadBalance]()
}

type AccessControlService interface {
	Getter[enity.AccessControl]
	Saver[enity.AccessControl]
	Deleter[enity.AccessControl]
}

func NewAccessControlService() AccessControlService {
	return New[enity.AccessControl]()
}

type ServiceInfoService interface {
	Getter[enity.ServiceInfo]
	Saver[enity.ServiceInfo]
	Deleter[enity.ServiceInfo]
	PagedLister[enity.ServiceInfo]
	AllGetter[enity.ServiceInfo]
	LoadTypeGrouper[enity.ServiceInfo]
	ServiceDetailGetter[enity.ServiceInfo]
}

func NewServiceInfoService() ServiceInfoService {
	return New[enity.ServiceInfo]()
}

type Getter[T Model] interface {
	Get(c *gin.Context, db *gorm.DB, search *T) (*T, error)
}

type Saver[T Model] interface {
	Save(c *gin.Context, db *gorm.DB, data *T) error
}

type Deleter[T Model] interface {
	Delete(c *gin.Context, db *gorm.DB, data *T) error
}

type PagedLister[T Model] interface {
	PageList(c *gin.Context, db *gorm.DB, queryConditions []func(db *gorm.DB) *gorm.DB, PageNo, PageSize int) ([]T, int64, error)
}

type AllGetter[T Model] interface {
	GetAll(c *gin.Context, db *gorm.DB, queryConditions []func(db *gorm.DB) *gorm.DB) ([]T, error)
}

type LoadTypeGrouper[T Model] interface {
	GetLoadTypeByGroup(c *gin.Context, tx *gorm.DB) ([]dto.DashServiceStatItemOutput, error)
}

type ServiceDetailGetter[T Model] interface {
	GetServiceDetail(c *gin.Context, db *gorm.DB, search *enity.ServiceInfo) (*enity.ServiceDetail, error)
}
type gormDao[T Model] struct{}

func New[T Model]() *gormDao[T] {
	return &gormDao[T]{}
}

func (dao *gormDao[T]) Get(c *gin.Context, db *gorm.DB, search *T) (*T, error) {
	// log记录查询信息
	log.Info("start getting", zap.Any("search", search), zap.String("trace_id", c.GetString("TraceID")))

	var out T

	result := db.Set("gorm:query_option", "FOR UPDATE").Where(search).First(&out)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			log.Error("record not found", zap.Any("search", search), zap.String("trace_id", c.GetString("TraceID")))
			return nil, result.Error
		}

		log.Error(fmt.Sprintf("error retrieving :%v ", search), zap.Error(result.Error), zap.String("trace_id", c.GetString("TraceID")))
		return nil, result.Error
	}

	log.Info("got successfully", zap.Any("out", out), zap.String("trace_id", c.GetString("TraceID")))
	return &out, nil
}

func (dao *gormDao[T]) Save(c *gin.Context, db *gorm.DB, data *T) error {
	// log记录保存信息
	log.Info("start saving", zap.Any("data", data), zap.String("trace_id", c.GetString("TraceID")))
	if err := db.Save(data).Error; err != nil {
		log.Error(fmt.Sprintf("error saving : %v ", data), zap.Error(err), zap.String("trace_id", c.GetString("TraceID")))
		return err
	}

	log.Info("saved successfully", zap.Any("data", data), zap.String("trace_id", c.GetString("TraceID")))
	return nil
}

func (dao *gormDao[T]) Delete(c *gin.Context, db *gorm.DB, data *T) error {
	// log记录删除信息
	log.Info("start deleting", zap.Any("data", data), zap.String("trace_id", c.GetString("TraceID")))

	if err := db.Delete(data).Error; err != nil {
		log.Error("error deleting", zap.Error(err), zap.String("trace_id", c.GetString("TraceID")))
		return err
	}

	log.Info("deleted successfully", zap.Any("data", data), zap.String("trace_id", c.GetString("TraceID")))
	return nil
}

func (dao *gormDao[T]) PageList(c *gin.Context, db *gorm.DB, queryConditions []func(db *gorm.DB) *gorm.DB, pageNo, pageSize int) ([]T, int64, error) {
	// log records query information
	log.Info("start pageList ", zap.String("trace_id", c.GetString("TraceID")))

	total := int64(0)
	list := []T{}
	offset := (pageNo - 1) * pageSize

	query := db.Where("is_delete=0")
	for _, condition := range queryConditions {
		query = condition(query)
	}
	if err := query.Limit(pageSize).Offset(offset).Order("id desc").Find(&list).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			log.Error(fmt.Sprintf("error retrieving :%v ", query), zap.Error(err), zap.String("trace_id", c.GetString("TraceID")))
			return nil, 0, err
		}
	}
	query.Limit(pageSize).Offset(offset).Count(&total)

	// log records success information
	log.Info("pageList successfully", zap.Any("list", list), zap.String("trace_id", c.GetString("TraceID")))
	return list, total, nil
}

// GetLoadTypeByGroup 根据服务ID和分组获取负载类型
func (dao *gormDao[T]) GetLoadTypeByGroup(c *gin.Context, tx *gorm.DB) ([]dto.DashServiceStatItemOutput, error) {
	// log记录开始查询
	log.Info("searching for group by load type", zap.String("trace_id", c.GetString("TraceID")))

	list := []dto.DashServiceStatItemOutput{}
	if err := tx.Table(enity.ServiceInfo{}.TableName()).Where("is_delete=0").Select("load_type, count(*) as value").Group("load_type").Scan(&list).Error; err != nil {
		log.Error("error retrieving", zap.Error(err), zap.String("trace_id", c.GetString("TraceID")))
		return nil, err
	}

	// log记录成功取到信息
	log.Info("group by load type was found", zap.String("trace_id", c.GetString("TraceID")))
	return list, nil
}

// GetServiceDetail
func (dao *gormDao[T]) GetServiceDetail(c *gin.Context, db *gorm.DB, search *enity.ServiceInfo) (*enity.ServiceDetail, error) {
	// log记录查询信息
	log.Info("start getting service detail")

	if search.ServiceName == "" {
		info, err := get(c, db, search)
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
		httpRule, err = get(c, db, &enity.HttpRule{ServiceID: search.ID})
		if err != nil && err != gorm.ErrRecordNotFound {
			log.Error("error retrieving http rule", zap.Error(err))
			return nil, err
		}
		log.Info("get http rule successful", zap.Any("httpRule", httpRule))
	case globals.LoadTypeTCP:
		tcpRule, err = get(c, db, &enity.TcpRule{ServiceID: search.ID})
		if err != nil && err != gorm.ErrRecordNotFound {
			log.Error("error retrieving tcp rule", zap.Error(err))
			return nil, err
		}
		log.Info("get tcp rule successful", zap.Any("tcpRule", tcpRule))
	case globals.LoadTypeGRPC:
		grpcRule, err = get(c, db, &enity.GrpcRule{ServiceID: search.ID})
		if err != nil && err != gorm.ErrRecordNotFound {
			log.Error("error retrieving grpc rule", zap.Error(err))
			return nil, err
		}
		log.Info("get grpc rule successful", zap.Any("grpcRule", grpcRule))
	}

	accessControl, err := get(c, db, &enity.AccessControl{ServiceID: search.ID})
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Error("error retrieving access control", zap.Error(err))
		return nil, err
	}
	log.Info("get access control successful", zap.Any("accessControl", accessControl))

	loadBalance, err := get(c, db, &enity.LoadBalance{ServiceID: search.ID})
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

// GetAll
func (dao *gormDao[T]) GetAll(c *gin.Context, db *gorm.DB, queryConditions []func(db *gorm.DB) *gorm.DB) ([]T, error) {
	// log记录查询信息
	log.Info("start pageList ", zap.String("trace_id", c.GetString("TraceID")))

	list := []T{}

	query := db.Where("is_delete=0")
	for _, condition := range queryConditions {
		query = condition(query)
	}
	if err := query.Order("id desc").Find(&list).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Error(fmt.Sprintf("error retrieving :%v ", query), zap.Error(err), zap.String("trace_id", c.GetString("TraceID")))
		return nil, err
	}

	// log记录成功信息
	log.Info("pageList successfully", zap.Any("list", list), zap.String("trace_id", c.GetString("TraceID")))
	return list, nil
}

func get[T Model](c *gin.Context, db *gorm.DB, search *T) (*T, error) {
	// log记录查询信息
	log.Info("start getting", zap.Any("search", search), zap.String("trace_id", c.GetString("TraceID")))

	var out T
	result := db.Where(search).First(&out)

	if result.Error != nil {
		// buf := make([]byte, 1<<16)
		// stackSize := runtime.Stack(buf, false)
		// log.Error("Error stack trace", zap.ByteString("stack", buf[:stackSize]))

		if result.Error == gorm.ErrRecordNotFound {
			log.Error("record not found", zap.Any("search", search), zap.String("trace_id", c.GetString("TraceID")))
			return nil, result.Error
		}

		log.Error(fmt.Sprintf("error retrieving :%v ", search), zap.Error(result.Error), zap.String("trace_id", c.GetString("TraceID")))
		return nil, result.Error
	}

	log.Info("got successfully", zap.Any("search", search), zap.String("trace_id", c.GetString("TraceID")))
	return &out, nil
}
