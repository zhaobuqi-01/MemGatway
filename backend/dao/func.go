package dao

import (
	"fmt"
	"gateway/backend/dto"
	"gateway/enity"
	"gateway/globals"
	"gateway/pkg/log"
	"sync"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Get查询单条数据
func Get[T Model](c *gin.Context, db *gorm.DB, search *T) (*T, error) {
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

// Save保存对象
func Save[T Model](c *gin.Context, db *gorm.DB, data *T) error {
	// log记录保存信息
	log.Info("start saving", zap.Any("data", data), zap.String("trace_id", c.GetString("TraceID")))
	if err := db.Save(data).Error; err != nil {
		log.Error(fmt.Sprintf("error saving : %v ", data), zap.Error(err), zap.String("trace_id", c.GetString("TraceID")))
		return err
	}

	log.Info("saved successfully", zap.Any("data", data), zap.String("trace_id", c.GetString("TraceID")))
	return nil
}

// delete删除对象
func Delete[T Model](c *gin.Context, db *gorm.DB, data *T) error {
	// log记录删除信息
	log.Info("start deleting", zap.Any("data", data), zap.String("trace_id", c.GetString("TraceID")))

	if err := db.Delete(data).Error; err != nil {
		log.Error("error deleting", zap.Error(err), zap.String("trace_id", c.GetString("TraceID")))
		return err
	}

	log.Info("deleted successfully", zap.Any("data", data), zap.String("trace_id", c.GetString("TraceID")))
	return nil
}

// ListByServiceID 根据服务ID查询列表
func ListByServiceID[T Model](c *gin.Context, db *gorm.DB, serviceID int64) ([]T, int64, error) {
	// log记录查询信息
	log.Info("start listByServiceID ", zap.String("trace_id", c.GetString("TraceID")))

	var list []T
	var count int64
	query := db.Select("*")
	query = query.Where("service_id=?", serviceID)
	err := query.Order("id desc").Find(&list).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Error(fmt.Sprintf("error retrieving :%v ", serviceID), zap.Error(err), zap.String("trace_id", c.GetString("TraceID")))
		return nil, 0, err
	}
	errCount := query.Count(&count).Error
	if errCount != nil {
		log.Error(fmt.Sprintf("error retrieving :%v ", serviceID), zap.Error(errCount), zap.String("trace_id", c.GetString("TraceID")))
		return nil, 0, err
	}

	log.Info("ListByServiceID successfully", zap.Any("list", list), zap.String("trace_id", c.GetString("TraceID")))
	return list, count, nil
}

// PageList 分页查询
func PageList[T Model](c *gin.Context, db *gorm.DB, queryConditions []func(db *gorm.DB) *gorm.DB, PageNo, PageSize int) ([]T, int64, error) {
	// log记录查询信息
	log.Info("start pageList ", zap.String("trace_id", c.GetString("TraceID")))

	total := int64(0)
	list := []T{}
	offset := (PageNo - 1) * PageSize

	query := db.Where("is_delete=0")
	for _, condition := range queryConditions {
		query = condition(query)
	}
	if err := query.Limit(PageSize).Offset(offset).Order("id desc").Find(&list).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Error(fmt.Sprintf("error retrieving :%v ", query), zap.Error(err), zap.String("trace_id", c.GetString("TraceID")))
		return nil, 0, err
	}
	query.Limit(PageSize).Offset(offset).Count(&total)

	// log记录成功信息
	log.Info("pageList successfully", zap.Any("list", list), zap.String("trace_id", c.GetString("TraceID")))
	return list, total, nil
}

// GetServiceDetail 获取服务详情
// func GetServiceDetail(c *gin.Context, db *gorm.DB, search *enity.ServiceInfo) (*enity.ServiceDetail, error) {
// 	// log记录查询信息
// 	log.Info("start getting service detail", zap.String("trace_id", c.GetString("TraceID")))

// 	if search.ServiceName == "" {
// 		info, err := Get(c, db, search)
// 		if err != nil {
// 			return nil, err
// 		}
// 		search = info
// 	}

// 	httpRule, err := Get(c, db, &enity.HttpRule{ServiceID: search.ID})
// 	if err != nil && err != gorm.ErrRecordNotFound {
// 		log.Error("error retrieving http rule", zap.Error(err), zap.String("trace_id", c.GetString("TraceID")))
// 		return nil, err
// 	}
// 	log.Info("get http rule successful", zap.Any("httpRule", httpRule), zap.String("trace_id", c.GetString("TraceID")))

// 	tcpRule, err := Get(c, db, &enity.TcpRule{ServiceID: search.ID})
// 	if err != nil && err != gorm.ErrRecordNotFound {
// 		log.Error("error retrieving tcp rule", zap.Error(err), zap.String("trace_id", c.GetString("TraceID")))
// 		return nil, err
// 	}
// 	log.Info("get tcp rule successful", zap.Any("tcpRule", tcpRule), zap.String("trace_id", c.GetString("TraceID")))

// 	grpcRule, err := Get(c, db, &enity.GrpcRule{ServiceID: search.ID})
// 	if err != nil && err != gorm.ErrRecordNotFound {
// 		log.Error("error retrieving grpc rule", zap.Error(err), zap.String("trace_id", c.GetString("TraceID")))
// 		return nil, err
// 	}
// 	log.Info("get grpc rule successful", zap.Any("grpcRule", grpcRule), zap.String("trace_id", c.GetString("TraceID")))

// 	accessControl, err := Get(c, db, &enity.AccessControl{ServiceID: search.ID})
// 	if err != nil && err != gorm.ErrRecordNotFound {
// 		log.Error("error retrieving access control", zap.Error(err), zap.String("trace_id", c.GetString("TraceID")))
// 		return nil, err
// 	}
// 	log.Info("get access control successful", zap.Any("accessControl", accessControl), zap.String("trace_id", c.GetString("TraceID")))

// 	loadBalance, err := Get(c, db, &enity.LoadBalance{ServiceID: search.ID})
// 	if err != nil && err != gorm.ErrRecordNotFound {
// 		log.Error("error retrieving load balance", zap.Error(err), zap.String("trace_id", c.GetString("TraceID")))
// 		return nil, err
// 	}
// 	log.Info("get load balance successfully", zap.Any("loadBalance", loadBalance), zap.String("trace_id", c.GetString("TraceID")))

// 	detail := &enity.ServiceDetail{
// 		Info:          search,
// 		HTTPRule:      httpRule,
// 		TCPRule:       tcpRule,
// 		GRPCRule:      grpcRule,
// 		LoadBalance:   loadBalance,
// 		AccessControl: accessControl,
// 	}

// 	// log记录成功取到信息
// 	log.Info("get service detail successful", zap.Any("detail", detail), zap.String("trace_id", c.GetString("TraceID")))
// 	return detail, nil
// }

func GetServiceDetail(c *gin.Context, db *gorm.DB, search *enity.ServiceInfo) (*enity.ServiceDetail, error) {
	// log记录查询信息
	log.Info("start getting service detail")

	if search.ServiceName == "" {
		info, err := Get(c, db, search)
		if err != nil {
			return nil, err
		}
		search = info
	}

	var (
		httpRule      *enity.HttpRule
		tcpRule       *enity.TcpRule
		grpcRule      *enity.GrpcRule
		accessControl *enity.AccessControl
		loadBalance   *enity.LoadBalance
	)

	var err, acErr, lbErr error
	var wg sync.WaitGroup

	wg.Add(3) // we will launch three goroutines

	// goroutine for getting rules
	go func() {
		defer wg.Done()

		switch search.LoadType {
		case globals.LoadTypeHTTP:
			httpRule, err = Get(c, db, &enity.HttpRule{ServiceID: search.ID})
			if err != nil && err != gorm.ErrRecordNotFound {
				log.Error("error retrieving http rule", zap.Error(err))
			} else {
				log.Info("get http rule successful", zap.Any("httpRule", httpRule))
			}
		case globals.LoadTypeTCP:
			tcpRule, err = Get(c, db, &enity.TcpRule{ServiceID: search.ID})
			if err != nil && err != gorm.ErrRecordNotFound {
				log.Error("error retrieving tcp rule", zap.Error(err))
			} else {
				log.Info("get tcp rule successful", zap.Any("tcpRule", tcpRule))
			}
		case globals.LoadTypeGRPC:
			grpcRule, err = Get(c, db, &enity.GrpcRule{ServiceID: search.ID})
			if err != nil && err != gorm.ErrRecordNotFound {
				log.Error("error retrieving grpc rule", zap.Error(err))
			} else {
				log.Info("get grpc rule successful", zap.Any("grpcRule", grpcRule))
			}
		}
	}()

	// goroutine for getting access control
	go func() {
		defer wg.Done()

		accessControl, acErr = Get(c, db, &enity.AccessControl{ServiceID: search.ID})
		if acErr != nil && acErr != gorm.ErrRecordNotFound {
			log.Error("error retrieving access control", zap.Error(acErr))
		} else {
			log.Info("get access control successful", zap.Any("accessControl", accessControl))
		}
	}()

	// goroutine for getting load balance
	go func() {
		defer wg.Done()

		loadBalance, lbErr = Get(c, db, &enity.LoadBalance{ServiceID: search.ID})
		if lbErr != nil && lbErr != gorm.ErrRecordNotFound {
			log.Error("error retrieving load balance", zap.Error(lbErr))
		} else {
			log.Info("get load balance successfully", zap.Any("loadBalance", loadBalance))
		}
	}()

	wg.Wait() // wait for all goroutines to finish

	// check errors after all goroutines have finished
	if err != nil || acErr != nil || lbErr != nil {
		return nil, fmt.Errorf("error getting service detail: %v, %v, %v", err, acErr, lbErr)
	}

	detail := &enity.ServiceDetail{
		Info:          search,
		HTTPRule:      httpRule,
		TCPRule:       tcpRule,
		GRPCRule:      grpcRule,
		LoadBalance:   loadBalance,
		AccessControl: accessControl,
	}

	// log记录成功取到信息
	log.Info("get service detail successful", zap.Any("detail", detail))
	return detail, nil
}
func GetLoadTypeByGroup(c *gin.Context, tx *gorm.DB) ([]dto.DashServiceStatItemOutput, error) {
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
